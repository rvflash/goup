// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/app"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/pkg/goup"
)

const (
	version   = "v0.1.0"
	notFound  = "/not-found"
	fileBuggy = "ok+err"
	fileErr   = "err"
	fileOK    = "ok"
	noop      = "no operation to do"
	nameInLog = ": "
	oops      = "oops"
)

func TestOpen(t *testing.T) {
	const wv = app.LogPrefix + "version " + version + "\n"
	var (
		are = is.New(t)
		dt  = map[string]struct {
			ctx    context.Context
			in     []string
			out    bool
			config goup.Config
			stderr string
			stdout string
		}{
			"default": {out: true, stderr: app.LogPrefix + "context canceled\n"},
			"empty":   {ctx: context.Background()},
			"version": {
				ctx:    context.Background(),
				config: goup.Config{PrintVersion: true},
				stderr: wv,
			},
			"invalid+version": {
				ctx:    context.Background(),
				in:     []string{notFound},
				out:    true,
				config: goup.Config{PrintVersion: true},
				stderr: wv + app.LogPrefix + "invalid go.mod\n",
			},
			"invalid": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				out:    true,
				config: goup.Config{PrintVersion: true, OnlyReleases: fileErr},
				stderr: wv + app.LogPrefix + nameInLog + oops + "\n",
			},
			"combo": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				out:    true,
				config: goup.Config{PrintVersion: true, OnlyReleases: fileBuggy},
				stderr: wv + app.LogPrefix + nameInLog + oops + "\n",
			},
			"verbose combo": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				out:    true,
				config: goup.Config{PrintVersion: true, OnlyReleases: fileBuggy, Verbose: true},
				stderr: wv + app.LogPrefix + nameInLog + oops + "\n",
				stdout: app.LogPrefix + nameInLog + noop + "\n",
			},
			"verbose ok": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				config: goup.Config{PrintVersion: true, OnlyReleases: "ok", Verbose: true},
				stderr: wv,
				stdout: app.LogPrefix + nameInLog + noop + "\n",
			},
			"ok": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				config: goup.Config{PrintVersion: true, OnlyReleases: "ok"},
				stderr: wv,
			},
			"recursive": {
				ctx:    context.Background(),
				in:     []string{"./..."},
				config: goup.Config{PrintVersion: true, OnlyReleases: "ok"},
				stderr: wv,
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			var stderr, stdout = new(strings.Builder), new(strings.Builder)
			a := newApp(t, stdout, stderr)
			a.Config = tt.config
			are.Equal(a.Check(tt.ctx, tt.in), tt.out) // mismatch result
			are.Equal(stdout.String(), tt.stdout)     // mismatch stdout
			are.Equal(stderr.String(), tt.stderr)     // mismatch stderr
		})
	}
}

func TestApp_Check(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("expected no panic")
		}
	}()
	var a app.App
	is.New(t).True(a.Check(nil, nil))
}

func TestWithChecker(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		a, err := app.Open(version, app.WithChecker(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // mismatch error
		are.True(a == nil)                         // mismatch result
	})
	t.Run("ok", func(t *testing.T) {
		a, err := app.Open(version, app.WithChecker(goup.CheckFile))
		are.Equal(err, nil) // mismatch error
		are.True(a != nil)  // mismatch result
	})
}

func TestWithErrOutput(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		a, err := app.Open(version, app.WithErrOutput(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // mismatch error
		are.True(a == nil)                         // mismatch result
	})
	t.Run("ok", func(t *testing.T) {
		b := new(strings.Builder)
		a, err := app.Open(version, app.WithErrOutput(b))
		are.Equal(err, nil) // mismatch error
		are.True(a != nil)  // mismatch result
	})
}

func TestWithOutput(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		a, err := app.Open(version, app.WithOutput(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // mismatch error
		are.True(a == nil)                         // mismatch result
	})
	t.Run("ok", func(t *testing.T) {
		b := new(strings.Builder)
		a, err := app.Open(version, app.WithOutput(b))
		are.Equal(err, nil) // mismatch error
		are.True(a != nil)  // mismatch result
	})
}

func TestWithParser(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		a, err := app.Open(version, app.WithParser(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // mismatch error
		are.True(a == nil)                         // mismatch result
	})
	t.Run("ok", func(t *testing.T) {
		a, err := app.Open(version, app.WithParser(mod.Parse))
		are.Equal(err, nil) // mismatch error
		are.True(a != nil)  // mismatch result
	})
}

var errOops = errors.New(oops)

type tip struct {
	err error
	msg string
}

func (t *tip) Err() error {
	return t.err
}

func (t *tip) String() string {
	return t.msg
}

type checker struct{}

// CheckFile implements the goup.Checker func.
func (p *checker) CheckFile(_ context.Context, _ mod.Mod, conf goup.Config) []goup.Tip {
	switch conf.OnlyReleases {
	case fileOK:
		return []goup.Tip{
			&tip{msg: noop},
		}
	case fileBuggy:
		return []goup.Tip{
			&tip{msg: noop},
			&tip{err: errOops},
		}
	case fileErr:
		return []goup.Tip{
			&tip{err: errOops},
		}
	default:
		return nil
	}
}

type parser struct{}

// Parse implements the mod.Parser func.
func (p *parser) Parse(path string) (*mod.File, error) {
	if strings.HasPrefix(path, notFound) {
		return nil, errup.ErrMod
	}
	return &mod.File{}, nil
}

func newApp(t *testing.T, stdout, stderr io.Writer) *app.App {
	var (
		c = &checker{}
		p = &parser{}
	)
	a, err := app.Open(
		version,
		app.WithChecker(c.CheckFile),
		app.WithErrOutput(stderr),
		app.WithOutput(stdout),
		app.WithParser(p.Parse),
	)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

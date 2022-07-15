// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app_test

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/app"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/log"
	"github.com/rvflash/goup/pkg/goup"
	"github.com/rvflash/goup/pkg/mod"
)

const (
	version      = "v0.1.0"
	notFound     = string(filepath.Separator) + "not-found"
	fileBuggy    = "ok+err"
	fileOutdated = "now ok"
	fileErr      = "err"
	fileOK       = "ok"
	noop         = "no operation to do"
	oops         = "oops"
)

func TestOpen(t *testing.T) {
	t.Parallel()
	const wv = log.Prefix + "version " + version + "\n"
	var (
		are = is.New(t)
		dt  = map[string]struct {
			ctx    context.Context
			in     []string
			out    bool
			config goup.Config
			stderr string
		}{
			"default": {out: true, stderr: log.Prefix + "context canceled\n"},
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
				stderr: wv + log.Prefix + "invalid go.mod\n",
			},
			"invalid": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				out:    true,
				config: goup.Config{PrintVersion: true, OnlyReleases: fileErr},
				stderr: wv + log.Prefix + oops + "\n",
			},
			"combo": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				out:    true,
				config: goup.Config{PrintVersion: true, OnlyReleases: fileBuggy},
				stderr: wv + log.Prefix + oops + "\n",
			},
			"verbose combo": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				out:    true,
				config: goup.Config{PrintVersion: true, OnlyReleases: fileBuggy, Verbose: true},
				stderr: wv + log.Prefix + oops + "\n",
			},
			"verbose ok": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				config: goup.Config{PrintVersion: true, OnlyReleases: "ok", Verbose: true},
				stderr: wv,
			},
			"ok": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				config: goup.Config{PrintVersion: true, OnlyReleases: "ok"},
				stderr: wv,
			},
			"force outdated": {
				ctx:    context.Background(),
				in:     []string{fileOK},
				config: goup.Config{PrintVersion: true, OnlyReleases: fileOutdated, Verbose: true},
				stderr: wv + log.Prefix + noop + "\n",
			},
			"recursive": {
				ctx:    context.Background(),
				in:     []string{"./..."},
				config: goup.Config{PrintVersion: true, OnlyReleases: "ok"},
				stderr: wv,
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var stderr = new(strings.Builder)
			a := newApp(t, stderr)
			a.Config = tt.config
			are.Equal(a.Check(tt.ctx, tt.in), tt.out) // mismatch result
			are.Equal(stderr.String(), tt.stderr)     // mismatch logger
		})
	}
}

func TestApp_Check(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Error("expected no panic")
		}
	}()
	var a app.App
	is.New(t).True(a.Check(nil, nil))
}

func TestWithChecker(t *testing.T) {
	t.Parallel()
	are := is.New(t)

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		a, err := app.Open(version, app.WithChecker(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // mismatch error
		are.True(a == nil)                         // mismatch result
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		a, err := app.Open(version, app.WithChecker(goup.Check))
		are.Equal(err, nil) // mismatch error
		are.True(a != nil)  // mismatch result
	})
}

func TestWithLogger(t *testing.T) {
	t.Parallel()
	are := is.New(t)

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		a, err := app.Open(version, app.WithLogger(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // mismatch error
		are.True(a == nil)                         // mismatch result
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		b := new(strings.Builder)
		l := log.New(b, false)
		a, err := app.Open(version, app.WithLogger(l))
		are.Equal(err, nil) // mismatch error
		are.True(a != nil)  // mismatch result
	})
}

func TestWithParser(t *testing.T) {
	t.Parallel()
	are := is.New(t)

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		a, err := app.Open(version, app.WithParser(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // mismatch error
		are.True(a == nil)                         // mismatch result
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		a, err := app.Open(version, app.WithParser(mod.Parse))
		are.Equal(err, nil) // mismatch error
		are.True(a != nil)  // mismatch result
	})
}

type checker struct{}

// check implements the goup.Checker func.
func (p *checker) Check(_ context.Context, _ mod.Mod, conf goup.Config) chan goup.Message {
	var (
		oops = errors.New(oops)
		ch   = make(chan goup.Message)
	)
	go func() {
		defer close(ch)
		switch conf.OnlyReleases {
		case fileOutdated:
			ch <- goup.NewEntry(goup.InfoLevel, "%s", noop)
		case fileBuggy:
			ch <- goup.NewEntry(goup.DebugLevel, "%s", noop)
			ch <- goup.NewEntry(goup.WarnLevel, "%s", oops)
		case fileErr:
			ch <- goup.NewEntry(goup.ErrorLevel, "%s", oops)
		}
	}()

	return ch
}

type parser struct{}

// Parse implements the mod.Parser func.
func (p *parser) Parse(path string) (*mod.File, error) {
	if strings.HasPrefix(path, notFound) {
		return nil, errup.ErrMod
	}
	return &mod.File{}, nil
}

func newApp(t *testing.T, stderr io.Writer) *app.App {
	t.Helper()
	var (
		c = &checker{}
		p = &parser{}
	)
	a, err := app.Open(
		version,
		app.WithChecker(c.Check),
		app.WithLogger(log.New(stderr, false)),
		app.WithParser(p.Parse),
	)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

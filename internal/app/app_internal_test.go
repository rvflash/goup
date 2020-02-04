// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/pkg/goup"
)

const version = "v0.1.0"

func TestWithChecker(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		app, err := Open(version, WithChecker(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // expected error
		are.Equal(app, nil)                        // unexpected app
	})
	t.Run("ok", func(t *testing.T) {
		app, err := Open(version, WithChecker(goup.GoUp{}))
		are.Equal(err, nil)        // unexpected error
		are.True(app.check != nil) // expected app with checker
	})
}

func TestWithErrOutput(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		app, err := Open(version, WithErrOutput(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // expected error
		are.Equal(app, nil)                        // unexpected app
	})
	t.Run("ok", func(t *testing.T) {
		b := new(strings.Builder)
		app, err := Open(version, WithErrOutput(b))
		are.Equal(err, nil)         // unexpected error
		are.True(app.stderr != nil) // expected app with error log
	})
}

func TestWithOutput(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		app, err := Open(version, WithOutput(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // expected error
		are.Equal(app, nil)                        // unexpected app
	})
	t.Run("ok", func(t *testing.T) {
		b := new(strings.Builder)
		app, err := Open(version, WithOutput(b))
		are.Equal(err, nil)        // unexpected error
		are.True(app.stdin != nil) // expected app with log
	})
}

func TestWithParser(t *testing.T) {
	are := is.New(t)
	t.Run("default", func(t *testing.T) {
		app, err := Open(version, WithParser(nil))
		are.True(errors.Is(err, errup.ErrMissing)) // expected error
		are.Equal(app, nil)                        // unexpected app
	})
	t.Run("ok", func(t *testing.T) {
		app, err := Open(version, WithParser(mod.Parse))
		are.Equal(err, nil)        // unexpected error
		are.True(app.parse != nil) // expected app with parser
	})
}

// todo
/*
func TestOpen(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			version string
			opts    []Configurator
			out     *App
			err     error
		}{
			"default": {out: &App{}},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := Open(tt.version, tt.opts...)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(out, tt.out)           // mismatch application
		})
	}
}
*/

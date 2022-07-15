// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/log"
	"github.com/rvflash/goup/pkg/goup"
)

func TestPatterns(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  []string
			out string
		}{
			"Default": {},
			"Ok":      {in: []string{"a", "b", "", "d", " e "}, out: "a,b,d,e"},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			are.Equal(tt.out, patterns(tt.in...)) // mismatch result
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()
	var (
		are    = is.New(t)
		stderr = log.New(new(strings.Builder), false)
		dt     = map[string]struct {
			ctx    context.Context
			cnf    goup.Config
			args   []string
			stderr log.Printer
			err    error
		}{
			"default":      {err: errup.ErrMissing},
			"no context":   {stderr: stderr, err: errup.ErrMod},
			"context only": {ctx: context.Background(), stderr: stderr, err: errup.ErrMod},
			"ok":           {ctx: context.Background(), cnf: goup.Config{PrintVersion: true}, stderr: stderr},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := run(tt.ctx, tt.cnf, tt.args, tt.stderr)
			are.True(errors.Is(err, tt.err)) // mismatch error
		})
	}
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs_test

import (
	"errors"
	"testing"

	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/vcs"
)

func TestErrorf(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			name string
			args []interface{}
			err  error
			msg  string
		}{
			"default": {},
			"one": {
				name: "git",
				args: []interface{}{errup.ErrMissing},
				err:  errup.ErrMissing,
				msg:  "git: missing data",
			},
			"more": {
				name: "go-get",
				args: []interface{}{errup.ErrSystem, errup.ErrMissing},
				err:  errup.ErrSystem,
				msg:  "go-get: invalid VCS: missing data",
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := vcs.Errorf(tt.name, tt.args...)
			are.True(errors.Is(err, tt.err)) // mismatch error
			if tt.err != nil {
				are.Equal(err.Error(), tt.msg) // mismatch message
			}
		})
	}
}

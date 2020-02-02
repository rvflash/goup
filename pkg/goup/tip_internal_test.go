// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestTip(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			tip Tip
			err error
			msg string
		}{
			"default": {tip: &tip{}},
			//"default": {tip: &tip{err: newError()}},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			are.True(errors.Is(tt.tip.Err(), tt.err)) // mismatch error
			are.Equal(tt.tip.String(), tt.msg)        // mismatch message
		})
	}
}

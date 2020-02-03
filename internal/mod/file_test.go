// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	uperr "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
)

func TestOpen(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in     []string
			module string
			depLen int
			err    error
		}{
			"default":        {err: uperr.ErrMod},
			"invalid name":   {in: []string{"testdata"}, err: uperr.ErrMod},
			"invalid path":   {in: []string{"testdata", mod.Filename}, err: uperr.ErrMod},
			"invalid go.mod": {in: []string{"..", "..", "testdata", "gomod", "invalid", mod.Filename}, err: uperr.ErrMod},
			"valid go.mod": {
				in:     []string{"..", "..", "testdata", "gomod", "valid", mod.Filename},
				module: "github.com/rvflash/goup",
				depLen: 4,
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := mod.Parse(filepath.Join(tt.in...))
			are.True(errors.Is(err, tt.err)) // mismatch error
			if tt.err == nil {
				are.Equal(out.Module(), tt.module)            // mismatch module
				are.Equal(len(out.Dependencies()), tt.depLen) // mismatch number of dependencies
			}
		})
	}
}

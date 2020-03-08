// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package path_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/path"
)

func TestMatch(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			glob    string
			target  string
			matched bool
		}{
			"Default":      {},
			"No target":    {glob: "a/b/c"},
			"No glob":      {target: "a/b/c"},
			"Wrong glob":   {glob: "b/*", target: "a/b/c"},
			"Invalid glob": {glob: "a/b/c/d", target: "a/b/c"},
			"Glob":         {glob: "a/b/*", target: "a/b/c", matched: true},
			"Prefix":       {glob: "a", target: "a/b/c", matched: true},
			"Parts":        {glob: "a/b", target: "a/b/c", matched: true},
			"Exact":        {glob: "a/b/c", target: "a/b/c", matched: true},
			"Complex":      {glob: "b,d/e,f.*,a/b/c", target: "a/b/c", matched: true},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			are.Equal(tt.matched, path.Match(tt.glob, tt.target)) // mismatch result
		})
	}
}

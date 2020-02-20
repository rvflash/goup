// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goget

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestParseMetaGoImport(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			path   []string
			system string
			uri    string
			err    error
		}{
			"empty": {
				path: []string{"..", "..", "..", "testdata", "golden", "goget", "empty.html"},
				err:  io.EOF,
			},
			"default": {
				path:   []string{"..", "..", "..", "testdata", "golden", "goget", "default.html"},
				system: "git",
				uri:    "https://go.googlesource.com/mod",
			},
			"no-charset": {
				path:   []string{"..", "..", "..", "testdata", "golden", "goget", "no-charset.html"},
				system: "git",
				uri:    "https://github.com/golang/appengine",
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(filepath.Join(tt.path...))
			if err != nil {
				t.Fatal(err)
			}
			sys, uri, err := parseMetaGoImport(f)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(uri, tt.uri)           // mismatch url
			are.Equal(sys, tt.system)        // mismatch system
		})
	}
}

func TestCharsetReader(t *testing.T) {
	var (
		are   = is.New(t)
		input = strings.NewReader("")
		dt    = map[string]struct {
			in   string
			out  io.Reader
			fail bool
		}{
			"default": {fail: true},
			"invalid": {in: "invalid", fail: true},
			"utf-8":   {in: "utf-8", out: input},
			"ascii":   {in: "ascii", out: input},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := charsetReader(tt.in, input)
			are.Equal(err != nil, tt.fail) // mismatch error
			are.Equal(out, tt.out)         // mismatch result
		})
	}
}

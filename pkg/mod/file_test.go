// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod_test

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/pkg/mod"
)

const (
	numDep = 14
	d0     = "github.com/rvflash/elapsed"
	d1     = "github.com/rvflash/backoff"
	d3     = "github.com/rvflash/goup"
	v0     = "v1.1.1"
	v1     = "v2.2.2"
)

var (
	invalidGoMod = []string{"..", "..", "testdata", "golden", "invalid", mod.Filename}
	updateGoMod  = []string{"..", "..", "testdata", "golden", "update", mod.Filename}
	updatedGoMod = []string{"..", "..", "testdata", "golden", "updated", mod.Filename}
	validGoMod   = []string{"..", "..", "testdata", "golden", "valid", mod.Filename}
)

func TestFile_Name(t *testing.T) {
	var f mod.File
	is.New(t).Equal(f.Name(), "")
}

func TestFile_Module(t *testing.T) {
	var f mod.File
	is.New(t).Equal(f.Module(), "")
}

func TestFile_Format(t *testing.T) {
	are := is.New(t)
	name, cleanup := newTmpGoMod(t)
	defer cleanup()

	out, err := mod.Parse(name)
	are.NoErr(err) // parse error
	buf, err := out.Format()
	are.True(errors.Is(err, errup.ErrNotModified))
	got, err := ioutil.ReadFile(name)
	are.NoErr(err)                       // expected source file
	are.Equal(buf, got)                  // expected no change
	are.NoErr(out.UpdateReplace(d0, v0)) // update replace
	are.NoErr(out.UpdateRequire(d1, v1)) // update require
	got, err = out.Format()
	are.NoErr(err) // writing failed
	exp, err := ioutil.ReadFile(filepath.Join(updatedGoMod...))
	are.NoErr(err)      // missing expecting
	are.Equal(got, exp) // mismatch data
}

func TestFile_Format2(t *testing.T) {
	var f mod.File
	are := is.New(t)
	buf, err := f.Format()
	are.Equal(err, errup.ErrMod)
	are.Equal(buf, nil)
}

func TestFile_UpdateRequire(t *testing.T) {
	var f mod.File
	is.New(t).Equal(f.UpdateRequire(d0, v0), errup.ErrMod)
}

func TestFile_UpdateReplace(t *testing.T) {
	are := is.New(t)
	out, err := mod.Parse(filepath.Join(validGoMod...))
	are.NoErr(err)
	are.True(errors.Is(out.UpdateReplace(d1, v1), errup.ErrMissing))
}

func TestFile_UpdateReplace2(t *testing.T) {
	var f mod.File
	is.New(t).Equal(f.UpdateReplace(d0, v0), errup.ErrMod)
}

func TestOpen(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in     []string
			module string
			depLen int
			err    error
		}{
			"default":        {err: errup.ErrMod},
			"invalid name":   {in: []string{"testdata"}, err: errup.ErrMod},
			"invalid path":   {in: []string{"testdata", mod.Filename}, err: errup.ErrMod},
			"invalid go.mod": {in: invalidGoMod, err: errup.ErrMod},
			"valid go.mod":   {in: validGoMod, module: d3, depLen: numDep},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := mod.Parse(filepath.Join(tt.in...))
			are.True(errors.Is(err, tt.err)) // mismatch error
			if tt.err == nil {
				t.Log(out.Name())
				are.Equal(out.Module(), tt.module)            // mismatch module
				are.Equal(len(out.Dependencies()), tt.depLen) // mismatch number of dependencies
			}
		})
	}
}

func newTmpGoMod(t *testing.T) (name string, cleanup func()) {
	dir, err := ioutil.TempDir("", "goup")
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadFile(filepath.Join(updateGoMod...))
	if err != nil {
		log.Fatal(err)
	}
	name = filepath.Join(dir, "go.mod")
	err = ioutil.WriteFile(name, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return name, func() {
		_ = os.RemoveAll(dir)
	}
}

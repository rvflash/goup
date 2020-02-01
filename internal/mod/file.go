// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"

	"golang.org/x/mod/modfile"
)

// Filename is the name of Go Module file.
const Filename = "go.mod"

// Mod represents a Go Module file.
type Mod interface {
	Module() string
	Dependencies() []Module
}

// File is a go.mod file.
type File struct {
	path string
	mods []Module
}

// OpenFile tries to open a go.mod file.
func OpenFile(path string) (*File, error) {
	if filepath.Base(path) != Filename {
		return nil, errors.ErrMod
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrMod, err.Error())
	}
	f, err := modfile.Parse(path, b, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrMod, err.Error())
	}
	return &File{
		path: f.Module.Mod.Path,
		mods: dependencies(f),
	}, nil
}

// Module returns the name of the module.
func (f *File) Module() string {
	return f.path
}

// Dependencies returns the dependencies of the module.
func (f *File) Dependencies() []Module {
	return f.mods
}

func dependencies(f *modfile.File) []Module {
	if f == nil {
		return nil
	}
	m := make(map[string]Module)

	// Firstly we get the modules used to replace legacy ones.
	for _, r := range f.Replace {
		m[r.Old.Path] = &module{
			path:    r.New.Path,
			version: semver.New(r.New.Version),
		}
	}
	// Then those required.
	for _, r := range f.Require {
		_, ok := m[r.Mod.Path]
		if ok {
			// Use the replace dependency instead of this one.
			continue
		}
		m[r.Mod.Path] = &module{
			indirect: r.Indirect,
			path:     r.Mod.Path,
			version:  semver.New(r.Mod.Version),
		}
	}
	var (
		i  int
		rs = make([]Module, len(m))
	)
	for _, d := range m {
		rs[i] = d
		i++
	}
	return rs
}

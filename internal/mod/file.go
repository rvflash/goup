// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod

import (
	"io/ioutil"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/rvflash/goup"
	"github.com/rvflash/goup/internal/semver"
)

// Filename
const Filename = "go.mod"

// Mod
type Mod interface {
	Module() string
	Dependencies() []Module
}

// File
type File struct {
	path string
	mods []Module
}

// OpenFile
func OpenFile(path string) (*File, error) {
	if filepath.Base(path) != Filename {
		return nil, goup.ErrMod
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := modfile.Parse(path, b, nil)
	if err != nil {
		return nil, err
	}
	return &File{
		path: f.Module.Mod.Path,
		mods: dependencies(f),
	}, nil
}

// Module
func (f *File) Module() string {
	return f.path
}

// Dependencies
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

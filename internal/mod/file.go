// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"

	"golang.org/x/mod/modfile"
)

//go:generate mockgen -destination ../../testdata/mock/mod/file.go -source file.go

// Filename is the name of Go module file.
const Filename = "go.mod"

// Parser defined the interface used to parse a go.mod file.
type Parser func(path string) (*File, error)

// Mod represents a Go module file.
type Mod interface {
	// Module returns the name of the module.
	Module() string
	// Dependencies returns the dependencies of the module.
	Dependencies() []Module
	// UpdateRequire adds an update of this required module path to the given version.
	UpdateRequire(path, version string) error
	// UpdateRequire adds an update on the replacement of this module path to the given version.
	UpdateReplace(oldPath, newVersion string) error
	// UpdateAndSave applies any requested updates to the file.
	UpdateAndSave() error
}

// File is a go.mod file.
type File struct {
	mods []Module

	raw *modfile.File
	mu  sync.RWMutex
}

// Parse tries to open a go.mod file.
func Parse(path string) (*File, error) {
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
		raw:  f,
		mods: dependencies(f),
	}, nil
}

// Module implements the Mod interface.
func (f *File) Module() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.raw == nil || f.raw.Module == nil {
		// Avoids panic.
		return ""
	}
	return f.raw.Module.Mod.Path
}

// Dependencies implements the Mod interface.
func (f *File) Dependencies() []Module {
	return f.mods
}

// UpdateRequire implements the Mod interface.
func (f *File) UpdateRequire(path, version string) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.raw == nil {
		return errors.ErrMod
	}
	return f.raw.AddRequire(path, version)
}

// UpdateReplace implements the Mod interface.
func (f *File) UpdateReplace(oldPath, newVersion string) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	newPath, err := findNewPath(f.raw, oldPath)
	if err != nil {
		return err
	}
	return f.raw.AddReplace(oldPath, "", newPath, newVersion)
}

func findNewPath(f *modfile.File, oldPath string) (string, error) {
	if f == nil {
		return "", errors.ErrMod
	}
	for _, r := range f.Replace {
		if r.Old.Path == oldPath {
			return r.New.Path, nil
		}
	}
	return "", errors.ErrMissing
}

// UpdateAndSave implements the Mod interface.
func (f *File) UpdateAndSave() error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.raw == nil {
		return errors.ErrMod
	}
	buf, err := f.raw.Format()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.raw.Syntax.Name, buf, 0644)
}

// dependencies returns the list of modules in this go.mod file.
// Firstly we get the modules used to replace legacy ones.
// Then those required. We use the replace dependency instead of this required.
func dependencies(f *modfile.File) []Module {
	var m = make(map[string]Module)
	for _, r := range f.Replace {
		// Ignores local replacements. like:
		// => ../vendor/example.com/group/pkg
		_, err := os.Stat(r.New.Path)
		if !os.IsNotExist(err) {
			continue
		}
		m[r.Old.Path] = &module{
			path:        r.New.Path,
			replacement: true,
			version:     semver.New(r.New.Version),
		}
	}
	for _, r := range f.Require {
		_, ok := m[r.Mod.Path]
		if ok {
			// Ignores known dependency (replace statement or duplicate).
			continue
		}
		m[r.Mod.Path] = &module{
			indirect: r.Indirect,
			path:     r.Mod.Path,
			version:  semver.New(r.Mod.Version),
		}
	}
	for _, r := range f.Exclude {
		_, ok := m[r.Mod.Path]
		if !ok {
			// Ignores exclusion of any unused dependency.
			continue
		}
		m[r.Mod.Path].(*module).excludeVersion = semver.New(r.Mod.Version)
	}
	return modules(m)
}

// modules converts the map of modules to a slice.
func modules(m map[string]Module) []Module {
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

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package mod exposes methods to parse a go.mod file.
package mod

import "github.com/rvflash/goup/internal/semver"

//go:generate mockgen -destination ../../testdata/mock/mod/module.go -source module.go

// Module represents a dependency.
type Module interface {
	Indirect() bool
	Path() string
	Replacement() bool
	Version() semver.Tag
	ExcludeVersions() []semver.Tag
}

type module struct {
	indirect,
	replacement bool
	path     string
	excludes []semver.Tag
	version  *semver.Version
}

// ExcludeVersions implements the module interface.
func (m *module) ExcludeVersions() []semver.Tag {
	return m.excludes
}

// Indirect implements the module interface.
func (m *module) Indirect() bool {
	return m.indirect
}

// Path implements the module interface.
func (m *module) Path() string {
	return m.path
}

// Replacement implements the module interface.
func (m *module) Replacement() bool {
	return m.replacement
}

// Version implements the module interface.
func (m *module) Version() semver.Tag {
	return m.version
}

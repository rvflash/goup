// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod

import "github.com/rvflash/goup/internal/semver"

// Module
type Module interface {
	Indirect() bool
	Path() string
	Version() semver.Tag
}

type module struct {
	indirect bool
	path     string
	version  *semver.Version
}

// Indirect  implements the Module interface.
func (m *module) Indirect() bool {
	return m.indirect
}

// Path  implements the Module interface.
func (m *module) Path() string {
	return m.path
}

// Version implements the Module interface.
func (m *module) Version() semver.Tag {
	return m.version
}

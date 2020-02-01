// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"fmt"
	"strings"
)

type errVCS string

// Error
func (e errVCS) Error() string {
	return string(e)
}

const (
	// ErrMod
	ErrMod = errVCS("invalid go.mod")
	// Err1System
	ErrSystem = errVCS("invalid VCS")
	// ErrRepository
	ErrRepository = errVCS("invalid repository")
	// ErrFetch
	ErrFetch = errVCS("failed to list tags")
)

// Errorf
func Errorf(vcs string, v ...interface{}) error {
	const (
		base = "%s: %w"
		more = ": %s"
	)
	switch len(v) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf(base, append([]interface{}{vcs}, v...)...)
	default:
		return fmt.Errorf(base+strings.Repeat(more, len(v)-1), append([]interface{}{vcs}, v...)...)
	}
}

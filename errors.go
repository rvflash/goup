// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"fmt"
	"strings"
)

type errUp string

// Error implements the error interface.
func (e errUp) Error() string {
	return string(e)
}

const (
	// ErrExpectedTag
	ErrExpectedTag = errUp("release tag expected")
	// ErrMod
	ErrMod = errUp("invalid go.mod")
	// Err1System
	ErrSystem = errUp("invalid VCS")
	// ErrRepository
	ErrRepository = errUp("invalid repository")
	// ErrFetch
	ErrFetch = errUp("failed to list tags")
	// ErrFormat
	ErrFormat = "invalid format or charset"
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

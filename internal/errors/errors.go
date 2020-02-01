// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package errors

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
)

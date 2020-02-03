// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package errors exposes the errors used by the GoUp app.
package errors

import (
	"errors"
	"fmt"
)

// NewCharset returns a new charset's error.
func NewCharset(charset string) error {
	return errors.New("unsupported charset:" + charset)
}

// NewMissingData returns the data is missing.
func NewMissingData(name string) error {
	return fmt.Errorf("%s: %w", name, ErrMissing)
}

type errUp string

// Error implements the error interface.
func (e errUp) Error() string {
	return string(e)
}

const (
	// ErrExpectedTag is returned when the version is not a release tag.
	ErrExpectedTag = errUp("release tag expected")
	// ErrMissing is returned when the data is missing.
	ErrMissing = errUp("missing data")
	// ErrMod is returned when the go.mod file is invalid.
	ErrMod = errUp("invalid go.mod")
	// ErrSystem is returned when the VCS does not respond to the remote request.
	ErrSystem = errUp("invalid VCS")
	// ErrRepository is returned when the repository is invalid.
	ErrRepository = errUp("invalid repository")
	// ErrFetch is returned when the fetching of versions failed.
	ErrFetch = errUp("failed to list tags")
)

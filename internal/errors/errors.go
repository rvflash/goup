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
	return errors.New("unsupported charset: " + charset)
}

// NewSecurityIssue returns a ne security issue based on the given url.
func NewSecurityIssue(url string) error {
	return fmt.Errorf("unsecured call to %s cancelled: %w", url, ErrFetch)
}

// NewMissingData returns the data is missing.
func NewMissingData(name string) error {
	return fmt.Errorf("%s: %w", name, ErrMissing)
}

type upError string

// Error implements the error interface.
func (e upError) Error() string {
	return string(e)
}

const (
	// ErrExpectedTag is returned when the version is not a release tag.
	ErrExpectedTag = upError("release tag expected")
	// ErrFetch is returned when the fetching of versions failed.
	ErrFetch = upError("failed to list tags")
	// ErrMissing is returned when the data is missing.
	ErrMissing = upError("missing data")
	// ErrMod is returned when the go.mod file is invalid.
	ErrMod = upError("invalid go.mod")
	// ErrNotModified is returned when the file has not changed.
	ErrNotModified = upError("not modified")
	// ErrRepository is returned when the repository is invalid.
	ErrRepository = upError("invalid repository")
	// ErrSystem is returned when the VCS does not respond to the remote request.
	ErrSystem = upError("invalid VCS")
)

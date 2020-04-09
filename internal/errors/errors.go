// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package errors exposes the errors used by the GoUp app.
package errors

import (
	"errors"
	"fmt"
)

// Failure contains an check error.
type Failure struct {
	Mod string
	Err error
}

// Error implements the error interface.
func (e *Failure) Error() string {
	if e.Err == nil {
		return ""
	}
	return fmt.Sprintf("%s check failed: %s", e.Mod, e.Err)
}

// Unwrap returns the error's source of the failure.
func (e *Failure) Unwrap() error {
	return e.Err
}

// OutOfDate contains an update error.
type OutOfDate struct {
	Mod,
	OldVersion,
	NewVersion string
}

// Error implements the error interface.
func (e *OutOfDate) Error() string {
	return fmt.Sprintf("%s %s must be updated with %s", e.Mod, e.OldVersion, e.NewVersion)
}

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

type errUp string

// Error implements the error interface.
func (e errUp) Error() string {
	return string(e)
}

const (
	// ErrExpectedTag is returned when the version is not a release tag.
	ErrExpectedTag = errUp("release tag expected")
	// ErrFetch is returned when the fetching of versions failed.
	ErrFetch = errUp("failed to list tags")
	// ErrMissing is returned when the data is missing.
	ErrMissing = errUp("missing data")
	// ErrMod is returned when the go.mod file is invalid.
	ErrMod = errUp("invalid go.mod")
	// ErrSystem is returned when the VCS does not respond to the remote request.
	ErrSystem = errUp("invalid VCS")
	// ErrRepository is returned when the repository is invalid.
	ErrRepository = errUp("invalid repository")
)

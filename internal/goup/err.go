// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"fmt"

	"github.com/rvflash/goup/internal/mod"
)

// NewError
func NewError(module mod.Module, err error) *Error {
	return &Error{
		Module: module,
		Src:    err,
	}
}

// NewOrder
func NewOrder(module mod.Module, msg string) *Error {
	return &Error{
		Module: module,
		Msg:    msg,
	}
}

// Error
type Error struct {
	mod.Module
	Msg string
	Src error
}

// Err
func (e *Error) Err() error {
	if e.Module == nil {
		return nil
	}
	if e.Src == nil {
		return fmt.Errorf("%s %s must be updated with %s", e.Module.Path(), e.Module.Version().String(), e.Msg)
	}
	return fmt.Errorf("%s check failed: %w", e.Module.Path(), e.Src)
}

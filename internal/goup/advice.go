// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"fmt"

	"github.com/rvflash/goup/internal/mod"
)

// NewAdvice
func NewAdvice(module mod.Module, msg string) *Advice {
	return &Advice{
		Module: module,
		Msg:    msg,
	}
}

// Advice
type Advice struct {
	mod.Module
	Msg string
}

// String
func (a *Advice) String() string {
	if a.Module == nil {
		return ""
	}
	return fmt.Sprintf("%s %s could be updated with %s", a.Module.Path(), a.Module.Version().String(), a.Msg)
}

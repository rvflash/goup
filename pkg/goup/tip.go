// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"fmt"
)

// Tip exposes methods to retrieve an order (with Err method) or an advice.
type Tip interface {
	Err() error
	fmt.Stringer
}

type tip struct {
	err error
	msg string
}

// Err implements the Tip interface.
func (t *tip) Err() error {
	return t.err
}

// String implements the Tip interface.
func (t *tip) String() string {
	if t.err != nil {
		return t.err.Error()
	}
	return t.msg
}

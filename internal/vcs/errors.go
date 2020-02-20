// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs

import (
	"fmt"
	"strings"
)

// Errorf returns a new error message by combining various errors with the name of the VCS as prefix.
func Errorf(vcs string, a ...interface{}) error {
	const (
		base = "%s: %w"
		more = ": %s"
	)
	switch len(a) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf(base, append([]interface{}{vcs}, a...)...)
	default:
		return fmt.Errorf(base+strings.Repeat(more, len(a)-1), append([]interface{}{vcs}, a...)...)
	}
}

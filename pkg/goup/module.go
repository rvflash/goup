// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"fmt"

	"github.com/rvflash/goup/internal/mod"
)

func newError(module mod.Module, err error) error {
	if module == nil || err == nil {
		return nil
	}
	return fmt.Errorf("%s check failed: %w", module.Path(), err)
}

func newOrder(module mod.Module, msg string) error {
	if module == nil {
		return nil
	}
	return fmt.Errorf("%s %s must be updated with %s", module.Path(), module.Version().String(), msg)
}

func checked(module mod.Module) string {
	if module == nil {
		return ""
	}
	return fmt.Sprintf("%s %s is up to date", module.Path(), module.Version().String())
}

func skipped(module mod.Module) string {
	if module == nil {
		return ""
	}
	return fmt.Sprintf("%s %s update skipped", module.Path(), module.Version().String())
}

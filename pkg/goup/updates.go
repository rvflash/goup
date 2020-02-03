// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package goup provides methods to check updates on go.mod file and modules.
package goup

import (
	"sync"

	"github.com/rvflash/goup/internal/mod"
)

// Updates contains any information about go.mod updates.
type Updates struct {
	mod.Mod

	mu   sync.RWMutex
	tips []Tip
}

// Tips returns all updates messages.
func (up *Updates) Tips() []Tip {
	up.mu.RLock()
	defer up.mu.RLock()
	return up.tips
}

func (up *Updates) could(s string) {
	up.mu.Lock()
	up.tips = append(up.tips, &tip{msg: s})
	up.mu.Unlock()
}

func (up *Updates) must(err error) {
	up.mu.Lock()
	up.tips = append(up.tips, &tip{err: err})
	up.mu.Unlock()
}

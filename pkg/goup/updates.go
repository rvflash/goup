// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package goup provides methods to check updates on go.mod file and modules.
package goup

import (
	"sync"

	"github.com/rvflash/goup/internal/mod"
)

// updates contains any information about go.mod updates.
type updates struct {
	mod.Mod

	mu  sync.RWMutex
	rs  []Tip
	bad bool
}

func (up *updates) tips() []Tip {
	up.mu.RLock()
	defer up.mu.RLock()
	return up.rs
}

func (up *updates) could(s string) {
	up.mu.Lock()
	up.rs = append(up.rs, &tip{msg: s})
	up.mu.Unlock()
}

func (up *updates) must(err error) {
	if err == nil {
		return
	}
	up.mu.Lock()
	up.rs = append(up.rs, &tip{err: err})
	up.bad = true
	up.mu.Unlock()
}

func (up *updates) failed() bool {
	up.mu.RLock()
	defer up.mu.RUnlock()
	return up.bad
}

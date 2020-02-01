// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"context"
	"sync"
	"time"

	"github.com/rvflash/goup"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/git"
)

// Config
type Config struct {
	Fast            bool
	Major           bool
	MajorMinor      bool
	OnlyReleases    string
	ExcludeIndirect bool
	Timeout         time.Duration
	Verbose         bool
	Version         bool
}

// Check
func Check(parent context.Context, m mod.Mod, conf *Config) *Update {
	up := &Update{Mod: m}
	if parent == nil || m == nil {
		up.errs = []error{goup.ErrMod}
		return up
	}
	ctx, cancel := context.WithTimeout(parent, conf.Timeout)
	defer cancel()

	var w8 sync.WaitGroup
	for _, d := range m.Dependencies() {
		w8.Add(1)
		go func(d mod.Module) {
			defer w8.Done()
			adv, err := Status(ctx, d, conf)
			if err != nil {
				up.must(err)
				return
			}
			up.could(adv)
		}(d)
	}
	w8.Wait()

	return up
}

// Status
func Status(ctx context.Context, d mod.Module, conf *Config) (string, error) {
	if ctx == nil || d == nil || conf == nil {
		return "", goup.ErrRepository
	}
	if conf.ExcludeIndirect && d.Indirect() {
		return "", nil
	}
	for _, remote := range []vcs.Remote{git.New()} {
		if remote.CanFetch(d.Path()) {
			vs, err := remote.FetchContext(ctx, d.Path())
			if err != nil {
				e := Error{
					Module: d,
					Src:    err,
				}
				return "", NewError(d, err).Err()
			}
			adv := Advice{
				Module: d,
				Msg:    "",
			}
			return adv, nil
		}
	}
	return "", nil
}

// Update
type Update struct {
	mod.Mod

	mu      sync.RWMutex
	errs    []error
	advices []string
}

// Advices
func (up *Update) Advices() []string {
	up.mu.RLock()
	defer up.mu.RLock()
	return up.advices
}

func (up *Update) could(s string) {
	up.mu.Lock()
	up.advices = append(up.advices, s)
	up.mu.Unlock()
}

// Errs
func (up *Update) Errs() []error {
	up.mu.RLock()
	defer up.mu.RLock()
	return up.errs
}
func (up *Update) must(err error) {
	up.mu.Lock()
	up.errs = append(up.errs, err)
	up.mu.Unlock()
}

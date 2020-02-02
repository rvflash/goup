// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package goup provides methods to check updates on go.mod file and modules.
package goup

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/git"
	"github.com/rvflash/goup/internal/vcs/goget"
)

// Config is used as the settings of the GoUp application.
type Config struct {
	ExcludeIndirect bool
	Fast            bool
	Major           bool
	MajorMinor      bool
	Verbose         bool
	Version         bool
	OnlyReleases    string
	Timeout         time.Duration
}

const delta = 1

// File checks the given go.mod file based on this configuration.
func File(parent context.Context, file mod.Mod, conf Config) *Updates {
	up := &Updates{Mod: file}
	if parent == nil || file == nil {
		up.must(errors.ErrMod)
		return up
	}
	ctx, cancel := context.WithTimeout(parent, conf.Timeout)
	defer cancel()

	var w8 sync.WaitGroup
	for _, d := range file.Dependencies() {
		w8.Add(delta)
		go func(d mod.Module) {
			defer w8.Done()
			adv, err := Module(ctx, d, conf)
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

// Module checks the version of the given module based on this configuration.
func Module(ctx context.Context, dep mod.Module, conf Config) (string, error) {
	if ctx == nil || dep == nil {
		return "", errors.ErrRepository
	}
	if conf.ExcludeIndirect && dep.Indirect() {
		return skipped(dep), nil
	}
	var (
		gitVCS     = git.New()
		httpClient = vcs.NewHTTPClient(conf.Timeout)
	)
	for _, remote := range []vcs.System{
		goget.New(httpClient, gitVCS),
		gitVCS,
	} {
		if !remote.CanFetch(dep.Path()) {
			continue
		}
		vs, err := remote.FetchPath(ctx, dep.Path())
		if err != nil {
			return "", newError(dep, err)
		}
		v, ok := latest(vs, dep, conf)
		if !ok {
			return checked(dep), nil
		}
		if semver.Compare(dep.Version(), v) < 0 {
			return "", newOrder(dep, v.String())
		}
		err = onlyTag(dep, conf.OnlyReleases)
		if err != nil {
			return "", newError(dep, err)
		}
		return checked(dep), nil
	}
	return "", newError(dep, errors.ErrSystem)
}

func latest(versions semver.Tags, dep mod.Module, conf Config) (semver.Tag, bool) {
	var v semver.Tag
	switch {
	case conf.Major:
		v = semver.Latest(versions)
	case conf.MajorMinor:
		v = semver.LatestMinor(dep.Version().Major(), versions)
	default:
		v = semver.LatestPatch(dep.Version().MajorMinor(), versions)
	}
	return v, v != nil
}

const sep = ","

func onlyTag(d mod.Module, paths string) error {
	if d == nil {
		return nil
	}
	for _, path := range strings.Split(paths, sep) {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		if strings.Contains(d.Path(), path) && !d.Version().IsTag() {
			return errors.ErrExpectedTag
		}
	}
	return nil
}

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

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"context"
	"path"
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
	ExcludeIndirect  bool
	ForceUpdate      bool
	Major            bool
	MajorMinor       bool
	PrintVersion     bool
	Strict           bool
	Verbose          bool
	InsecurePatterns string
	OnlyReleases     string
	Timeout          time.Duration
}

// Check is the default checker of go.mod file.
func Check(ctx context.Context, file mod.Mod, conf Config) []Tip {
	return New(conf).CheckFile(ctx, file, conf)
}

// Checker must be implemented to check updates on go mod file or module.
type Checker func(ctx context.Context, file mod.Mod, conf Config) []Tip

// Setter defines the interface used to set settings.
type Setter func(u *GoUp)

// SetGit set the VCS git.
func SetGit(git vcs.System) Setter {
	return func(u *GoUp) {
		u.git = git
	}
}

// SetGoGet sets the VCS go-get.
func SetGoGet(goGet vcs.System) Setter {
	return func(u *GoUp) {
		u.goGet = goGet
	}
}

// New returns a new instance of GoUp with the default dependencies checkers.
func New(conf Config, sets ...Setter) *GoUp {
	var (
		u          = new(GoUp)
		httpClient = vcs.NewHTTPClient(conf.Timeout, conf.InsecurePatterns)
		gitVCS     = git.New(httpClient)
	)
	sets = append([]Setter{
		SetGit(gitVCS),
		SetGoGet(goget.New(httpClient, gitVCS)),
	}, sets...)
	for _, set := range sets {
		set(u)
	}
	return u
}

// GoUp allows to check a go.mod file with each of these dependencies.
type GoUp struct {
	git, goGet vcs.System
}

const delta = 1

// CheckFile verifies the given go.mod file based on this configuration.
// It's implements the Checker interface.
func (u *GoUp) CheckFile(parent context.Context, file mod.Mod, conf Config) []Tip {
	up := &updates{Mod: file}
	if !u.ready(parent) || file == nil {
		up.must(errors.ErrMod)
		return up.tips()
	}
	ctx, cancel := context.WithTimeout(parent, conf.Timeout)
	defer cancel()
	var w8 sync.WaitGroup
	for _, d := range file.Dependencies() {
		w8.Add(delta)
		go func(d mod.Module) {
			defer w8.Done()
			adv, err := u.CheckModule(ctx, d, conf)
			if err != nil {
				up.must(err)
				return
			}
			up.could(adv)
		}(d)
	}
	w8.Wait()

	return up.tips()
}

// CheckModule checks the version of the given module based on this configuration.
func (u *GoUp) CheckModule(ctx context.Context, dep mod.Module, conf Config) (string, error) {
	if !u.ready(ctx) || dep == nil {
		return "", errors.ErrRepository
	}
	if conf.ExcludeIndirect && dep.Indirect() {
		return skipped(dep), nil
	}
	for _, system := range []vcs.System{u.goGet, u.git} {
		if !system.CanFetch(dep.Path()) {
			continue
		}
		vs, err := system.FetchPath(ctx, dep.Path())
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

func (u *GoUp) ready(ctx context.Context) bool {
	return ctx != nil && u.goGet != nil && u.git != nil
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

func onlyTag(d mod.Module, globs string) error {
	var matched bool
	for _, glob := range strings.Split(globs, sep) {
		if glob = strings.TrimSpace(glob); glob == "" {
			continue
		}
		matched, _ = path.Match(glob, d.Path())
		if matched && !d.Version().IsTag() {
			return errors.ErrExpectedTag
		}
	}
	return nil
}

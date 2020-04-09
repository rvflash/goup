// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"context"
	"sync"
	"time"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/internal/path"
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

// Checker must be implemented to check updates on go.mod file or module.
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

// GoUp allows to check a go.dep file with each of these dependencies.
type GoUp struct {
	git, goGet vcs.System
}

const delta = 1

// CheckFile verifies the given go.mod file based on this configuration.
// Regarding the configuration, it applies any update advices to the go.mod file.
// If one or more checks failed, the update is cancelled.
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
			switch e := err.(type) {
			case nil:
				up.could(adv)
			case *errors.OutOfDate:
				if !conf.ForceUpdate {
					up.must(err)
					return
				}
				if d.Replacement() {
					up.must(file.UpdateReplace(e.Mod, e.NewVersion))
				} else {
					up.must(file.UpdateRequire(e.Mod, e.NewVersion))
				}
			default:
				up.must(err)
			}
		}(d)
	}
	w8.Wait()

	if !up.failed() && conf.ForceUpdate {
		up.must(file.UpdateAndSave())
	}
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
			return "", &errors.Failure{Mod: dep.Path(), Err: err}
		}
		x, ok := dep.ExcludeVersion()
		if ok {
			vs = vs.Not(x)
		}
		v, ok := latest(vs, dep, conf)
		if !ok {
			return checked(dep), nil
		}
		if semver.Compare(dep.Version(), v) < 0 {
			return "", &errors.OutOfDate{Mod: dep.Path(), OldVersion: dep.Version().String(), NewVersion: v.String()}
		}
		err = onlyTag(dep, conf.OnlyReleases)
		if err != nil {
			return "", &errors.Failure{Mod: dep.Path(), Err: err}
		}
		return checked(dep), nil
	}
	return "", &errors.Failure{Mod: dep.Path(), Err: errors.ErrSystem}
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

func onlyTag(d mod.Module, globs string) error {
	if path.Match(globs, d.Path()) && !d.Version().IsTag() {
		return errors.ErrExpectedTag
	}
	return nil
}

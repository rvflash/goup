// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package goup provides methods to check updates on go.mod file and modules.
package goup

import (
	"context"
	"errors"
	"github.com/rvflash/workr"
	"io/ioutil"
	"sync/atomic"
	"time"

	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/path"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/git"
	"github.com/rvflash/goup/internal/vcs/goget"
	"github.com/rvflash/goup/pkg/mod"
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

// Checker must be implemented to checkFile updates on go.mod file or module.
type Checker func(ctx context.Context, file mod.Mod, conf Config) chan Message

// Check is the default checker of go.mod file.
func Check(ctx context.Context, file mod.Mod, conf Config) chan Message {
	chk := newGoUp(conf)
	go chk.checkFile(ctx, file)
	return chk.log
}

// newGoUp returns a new instance of GoUp with the default dependencies checkers.
func newGoUp(conf Config, sets ...setter) *goUp {
	var (
		u = &goUp{
			Config: conf,
			log:    make(chan Message),
		}
		httpClient = vcs.NewHTTPClient(conf.Timeout, conf.InsecurePatterns)
		gitVCS     = git.New(httpClient)
	)
	sets = append([]setter{
		setGit(gitVCS),
		setGoGet(goget.New(httpClient, gitVCS)),
	}, sets...)
	for _, set := range sets {
		set(u)
	}
	return u
}

type goUp struct {
	Config
	git, goGet vcs.System
	log        chan Message
}

const (
	delta = 1
	perm  = 0644
)

// checkFile verifies the given go.mod file based on this configuration.
// Regarding the configuration, it applies any update advices to the go.mod file.
// If one or more checks failed, the update is cancelled.
func (e *goUp) checkFile(parent context.Context, file mod.Mod) {
	defer close(e.log)
	if !e.ready(parent) || file == nil {
		e.log <- newError(errup.ErrMod, file)
		return
	}
	ctx, cancel := context.WithTimeout(parent, e.Timeout)
	defer cancel()
	bad := e.checkDependencies(ctx, file)
	if !e.ForceUpdate {
		return
	}
	if bad > 0 {
		e.log <- newError(errup.ErrNotModified, file)
		return
	}
	err := updateFile(file)
	if err != nil {
		e.log <- newError(err, file)
	}
}

// Errors are internally managed with a dedicated channel.
// So we only return each task as succeeded and eventually the number of fails.
func (e *goUp) checkDependencies(parent context.Context, file mod.Mod) uint64 {
	grp, ctx := workr.WithContext(parent)
	var bad uint64
	for _, d := range file.Dependencies() {
		dep := d
		grp.Go(func() (err error) {
			log := e.checkDependency(ctx, dep)
			v, ok := log.OutDated()
			if !ok || !e.ForceUpdate {
				if log.Level() < InfoLevel {
					atomic.AddUint64(&bad, delta)
				}
				e.log <- log
				return nil
			}
			if dep.Replacement() {
				err = file.UpdateReplace(dep.Path(), v)
			} else {
				err = file.UpdateRequire(dep.Path(), v)
			}
			if err != nil {
				atomic.AddUint64(&bad, delta)
				e.log <- newFailure(err, dep)
			} else {
				e.log <- newUpdate(dep, v)
			}
			return nil
		})
	}
	_ = grp.Wait()
	return bad
}

// checkDependency checks the version of the given module based on this configuration.
func (e *goUp) checkDependency(ctx context.Context, dep mod.Module) *Entry {
	if e.ExcludeIndirect && dep.Indirect() {
		return newSkip(dep)
	}
	for _, system := range []vcs.System{e.goGet, e.git} {
		if !system.CanFetch(dep.Path()) {
			continue
		}
		vs, err := system.FetchPath(ctx, dep.Path())
		if err != nil {
			return newFailure(err, dep)
		}
		x, ok := dep.ExcludeVersion()
		if ok {
			vs = vs.Not(x)
		}
		v, ok := latest(vs, dep, e.Config.Major, e.Config.MajorMinor)
		if !ok {
			return newCheck(dep)
		}
		if semver.Compare(dep.Version(), v) < 0 {
			return newOutOfDate(dep, v.String())
		}
		err = onlyTag(dep, e.OnlyReleases)
		if err != nil {
			return newFailure(err, dep)
		}
		return newCheck(dep)
	}
	return newFailure(errup.ErrSystem, dep)
}

func updateFile(file mod.Mod) error {
	buf, err := file.Format()
	if err != nil {
		if !errors.Is(err, errup.ErrNotModified) {
			return err
		}
		return nil
	}
	return ioutil.WriteFile(file.Name(), buf, perm)
}

func (e *goUp) ready(ctx context.Context) bool {
	return ctx != nil && e.log != nil && e.goGet != nil && e.git != nil
}

func latest(versions semver.Tags, dep mod.Module, major, majorMinor bool) (semver.Tag, bool) {
	var v semver.Tag
	switch {
	case major:
		v = semver.Latest(versions)
	case majorMinor:
		v = semver.LatestMinor(dep.Version().Major(), versions)
	default:
		v = semver.LatestPatch(dep.Version().MajorMinor(), versions)
	}
	return v, v != nil
}

func onlyTag(d mod.Module, globs string) error {
	if path.Match(globs, d.Path()) && !d.Version().IsTag() {
		return errup.ErrExpectedTag
	}
	return nil
}

// setter defines the interface used to set settings.
type setter func(u *goUp)

// setGit set the VCS git.
func setGit(git vcs.System) setter {
	return func(u *goUp) {
		u.git = git
	}
}

// setGoGet sets the VCS go-get.
func setGoGet(goGet vcs.System) setter {
	return func(u *goUp) {
		u.goGet = goGet
	}
}

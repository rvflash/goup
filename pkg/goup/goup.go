// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

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

// Checker must be implemented to check updates on go mod file or module.
type Checker func(ctx context.Context, file mod.Mod, conf Config) []Tip

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

const delta = 1

// CheckFile verifies the given go.mod file based on this configuration.
func CheckFile(parent context.Context, file mod.Mod, conf Config) []Tip {
	up := &updates{Mod: file}
	if parent == nil || file == nil {
		up.must(errors.ErrMod)
		return nil
	}
	ctx, cancel := context.WithTimeout(parent, conf.Timeout)
	defer cancel()

	var w8 sync.WaitGroup
	for _, d := range file.Dependencies() {
		w8.Add(delta)
		go func(d mod.Module) {
			defer w8.Done()
			adv, err := CheckModule(ctx, d, conf)
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
func CheckModule(ctx context.Context, dep mod.Module, conf Config) (string, error) {
	if ctx == nil || dep == nil {
		return "", errors.ErrRepository
	}
	if conf.ExcludeIndirect && dep.Indirect() {
		return skipped(dep), nil
	}
	var (
		gitVCS     = git.New()
		httpClient = vcs.NewHTTPClient(conf.Timeout, conf.InsecurePatterns)
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

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/pkg/mod"
	mockMod "github.com/rvflash/goup/testdata/mock/mod"
	mockVCS "github.com/rvflash/goup/testdata/mock/vcs"

	"go.uber.org/mock/gomock"
)

func TestGoUp_CheckDependency(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var (
		sy1 = newSystem(ctrl, semver.Tags{semver.New(v0)}, nil)
		are = is.New(t)
		dt  = map[string]struct {
			system vcs.System
			ctx    context.Context
			module mod.Module
			cnf    Config
			level  Level
			format string
		}{
			"skip indirect": {
				system: sy1,
				ctx:    ctx,
				module: newModule(ctrl, true),
				cnf:    Config{ExcludeIndirect: true},
				level:  DebugLevel,
				format: "update skipped",
			},
			"not matches vcs": {
				system: newNoSystem(ctrl),
				ctx:    ctx,
				module: newModule(ctrl, true),
				level:  ErrorLevel,
				format: "check failed",
			},
			"ok": {
				system: sy1,
				ctx:    ctx,
				module: newModule(ctrl, false),
				cnf:    Config{ExcludeIndirect: true},
				level:  DebugLevel,
				format: "up to date",
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			u := newGoUp(tt.cnf, setGoGet(tt.system), setGit(tt.system))
			e := u.checkDependency(tt.ctx, tt.module)
			are.Equal(tt.level, e.Level())                    // mismatch level
			are.True(strings.Contains(e.Format(), tt.format)) // mismatch format
		})
	}
}

func TestUpdateFile(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	are := is.New(t)
	dir, err := ioutil.TempDir("", "goup")
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	are.NoErr(err)

	dt := map[string]struct {
		file    mod.Mod
		err     error
		updated bool
	}{
		"Not modified": {file: newTmpMod(ctrl, filepath.Join(dir, "t0"), errup.ErrNotModified)},
		"Failed":       {file: newTmpMod(ctrl, filepath.Join(dir, "t1"), errup.ErrFetch), err: errup.ErrFetch},
		"Default":      {file: newTmpMod(ctrl, filepath.Join(dir, "t2"), nil), updated: true},
	}
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			err := updateFile(tt.file)
			are.True(errors.Is(err, tt.err))                  // mismatch error
			are.Equal(tt.updated, fileExists(tt.file.Name())) // mismatch file "created"
		})
	}
}

func TestLatest(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		res = semver.Tags{
			semver.New("v0.1.2"),
			semver.New("v0.1.3"),
			semver.New("v1.1.2"),
			semver.New("v2.1.2"),
			semver.New("v0.2.2"),
			semver.New("v0.2.3"),
		}
		are = is.New(t)
		dt  = map[string]struct {
			in  semver.Tags
			dep mod.Module
			cnf Config
			out semver.Tag
			ok  bool
		}{
			"default": {in: res, dep: newVer(ctrl, "v0.1.2"), out: semver.New("v0.1.3"), ok: true},
			"major": {
				in:  res,
				dep: mockMod.NewMockModule(ctrl),
				cnf: Config{Major: true},
				out: semver.New("v2.1.2"),
				ok:  true,
			},
			"major+minor": {
				in:  res,
				dep: newVer(ctrl, "v0.1.2"),
				cnf: Config{MajorMinor: true},
				out: semver.New("v0.2.3"),
				ok:  true,
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			out, ok := latest(tt.in, tt.dep, tt.cnf.Major, tt.cnf.MajorMinor)
			are.Equal(out, tt.out) // mismatch tag
			are.Equal(ok, tt.ok)   // mismatch found
		})
	}
}

func TestOnlyTag(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			dep   mod.Module
			paths string
			err   error
		}{
			"Default":              {dep: newTag(ctrl, "")},
			"Invalid glob pattern": {dep: newTag(ctrl, "v1.0.0-b42"), paths: "another.com"},
			"Valid glob pattern": {
				dep:   newTag(ctrl, "v1.0.0-b42"),
				paths: "example.com/*/*",
				err:   errup.ErrExpectedTag,
			},
			"Skip": {dep: newTag(ctrl, "v1.0.0"), paths: "test,,example.com"},
			"Valid prefix": {
				dep:   newTag(ctrl, "v1.0.0-b42"),
				paths: "example.com",
				err:   errup.ErrExpectedTag,
			},
			"Ok": {dep: newTag(ctrl, "v1.0.0"), paths: "example.com/pkg/*"},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			err := onlyTag(tt.dep, tt.paths)
			are.Equal(err, tt.err) // mismatch error
		})
	}
}

func newModule(ctrl *gomock.Controller, indirect bool) *mockMod.MockModule {
	m := mockMod.NewMockModule(ctrl)
	m.EXPECT().Path().Return(repoName).AnyTimes()
	m.EXPECT().Version().Return(semver.New(v0)).AnyTimes()
	m.EXPECT().Indirect().Return(indirect).AnyTimes()
	m.EXPECT().ExcludeVersions().Return(nil).AnyTimes()
	return m
}

func newNoSystem(ctrl *gomock.Controller) *mockVCS.MockSystem {
	m := mockVCS.NewMockSystem(ctrl)
	m.EXPECT().CanFetch(gomock.Any()).Return(false).AnyTimes()
	return m
}

func newSystem(ctrl *gomock.Controller, tags semver.Tags, err error) *mockVCS.MockSystem {
	m := mockVCS.NewMockSystem(ctrl)
	m.EXPECT().CanFetch(gomock.Any()).Return(true).AnyTimes()
	m.EXPECT().FetchPath(gomock.Any(), gomock.Any()).Return(tags, err).AnyTimes()
	return m
}

func newTag(ctrl *gomock.Controller, v string) *mockMod.MockModule {
	d := mockMod.NewMockModule(ctrl)
	d.EXPECT().Path().Return(repoName).Times(oneTime)
	d.EXPECT().Version().Return(semver.New(v)).AnyTimes()
	return d
}

func newVer(ctrl *gomock.Controller, v string) *mockMod.MockModule {
	d := mockMod.NewMockModule(ctrl)
	d.EXPECT().Version().Return(semver.New(v)).AnyTimes()
	return d
}

func newTmpMod(ctrl *gomock.Controller, name string, err error) *mockMod.MockMod {
	m := mockMod.NewMockMod(ctrl)
	if err != nil {
		m.EXPECT().Format().Return(nil, err).Times(oneTime)
		m.EXPECT().Name().Return(name).AnyTimes()
	} else {
		m.EXPECT().Format().Return([]byte("foo bar"), nil).Times(oneTime)
		m.EXPECT().Name().Return(name).MinTimes(oneTime)
	}

	return m
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

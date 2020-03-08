// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
	mock_mod "github.com/rvflash/goup/testdata/mock/mod"
	mock_vcs "github.com/rvflash/goup/testdata/mock/vcs"
)

func TestCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	is.New(t).Equal(Check(ctx, newMockMod(ctrl, nil), Config{}), nil)
}

func TestGoUp_CheckFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var (
		sy1 = newMockSystem(ctrl, semver.Tags{semver.New(release)}, nil)
		are = is.New(t)
		dt  = map[string]struct {
			system vcs.System
			ctx    context.Context
			file   mod.Mod
			cnf    Config
			res    []Tip
		}{
			"default":       {system: sy1, res: newErrTip(errup.ErrMod)},
			"no context":    {system: sy1, file: mock_mod.NewMockMod(ctrl), res: newErrTip(errup.ErrMod)},
			"no file":       {system: sy1, ctx: ctx, res: newErrTip(errup.ErrMod)},
			"no dependency": {system: sy1, ctx: ctx, file: newMockMod(ctrl, nil)},
			"up to date": {
				system: sy1,
				ctx:    ctx,
				file: newMockMod(ctrl, []mod.Module{
					newMockModule(ctrl, false),
				}),
				res: newTip("example.com/group/go v1.0.2 is up to date"),
			},
			"outdated": {
				system: newMockSystem(ctrl, semver.Tags{semver.New("v1.0.3")}, nil),
				ctx:    ctx,
				file: newMockMod(ctrl, []mod.Module{
					newMockModule(ctrl, false),
				}),
				res: newErrTip(errors.New("example.com/group/go v1.0.2 must be updated with v1.0.3")),
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := New(Config{}, SetGoGet(tt.system), SetGit(tt.system)).CheckFile(tt.ctx, tt.file, tt.cnf)
			are.Equal(res, tt.res) // mismatch result
		})
	}
}

func TestGoUp_CheckModule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var (
		sy1 = newMockSystem(ctrl, semver.Tags{semver.New(release)}, nil)
		are = is.New(t)
		dt  = map[string]struct {
			system vcs.System
			ctx    context.Context
			module mod.Module
			cnf    Config
			out    string
			err    error
		}{
			"default":    {system: sy1, err: errup.ErrRepository},
			"no context": {system: sy1, module: mock_mod.NewMockModule(ctrl), err: errup.ErrRepository},
			"no module":  {system: sy1, ctx: ctx, err: errup.ErrRepository},
			"skip indirect": {
				system: sy1,
				ctx:    ctx,
				module: newMockModule(ctrl, true),
				cnf:    Config{ExcludeIndirect: true},
				out:    "example.com/group/go v1.0.2 update skipped",
			},
			"not matches vcs": {
				system: newMockNoSystem(ctrl),
				ctx:    ctx,
				module: newMockModule(ctrl, true),
				err:    errup.ErrSystem,
			},
			"ok": {
				system: sy1,
				ctx:    ctx,
				module: newMockModule(ctrl, false),
				cnf:    Config{ExcludeIndirect: true},
				out:    "example.com/group/go v1.0.2 is up to date",
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			u := New(Config{}, SetGoGet(tt.system), SetGit(tt.system))
			out, err := u.CheckModule(tt.ctx, tt.module, tt.cnf)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(out, tt.out)           // mismatch result
		})
	}
}

func newMockNoSystem(ctrl *gomock.Controller) *mock_vcs.MockSystem {
	m := mock_vcs.NewMockSystem(ctrl)
	m.EXPECT().CanFetch(gomock.Any()).Return(false).AnyTimes()
	return m
}

func newMockSystem(ctrl *gomock.Controller, tags semver.Tags, err error) *mock_vcs.MockSystem {
	m := mock_vcs.NewMockSystem(ctrl)
	m.EXPECT().CanFetch(gomock.Any()).Return(true).AnyTimes()
	m.EXPECT().FetchPath(gomock.Any(), gomock.Any()).Return(tags, err).AnyTimes()
	return m
}

func newMockModule(ctrl *gomock.Controller, indirect bool) *mock_mod.MockModule {
	m := mock_mod.NewMockModule(ctrl)
	m.EXPECT().Path().Return(repoName).AnyTimes()
	m.EXPECT().Version().Return(semver.New(release)).AnyTimes()
	m.EXPECT().Indirect().Return(indirect).AnyTimes()
	return m
}

func TestCheckModule(t *testing.T) {

}

func TestLatest(t *testing.T) {
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
			"default": {in: res, dep: newDep(ctrl, "v0.1.2"), out: semver.New("v0.1.3"), ok: true},
			"major": {
				in:  res,
				dep: mock_mod.NewMockModule(ctrl),
				cnf: Config{Major: true},
				out: semver.New("v2.1.2"),
				ok:  true,
			},
			"major+minor": {
				in:  res,
				dep: newDep(ctrl, "v0.1.2"),
				cnf: Config{MajorMinor: true},
				out: semver.New("v0.2.3"),
				ok:  true,
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, ok := latest(tt.in, tt.dep, tt.cnf)
			are.Equal(out, tt.out) // mismatch tag
			are.Equal(ok, tt.ok)   // mismatch found
		})
	}
}

func TestOnlyTag(t *testing.T) {
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
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			err := onlyTag(tt.dep, tt.paths)
			are.Equal(err, tt.err) // mismatch error
		})
	}
}

const oneTime = 1

func newDep(ctrl *gomock.Controller, v string) *mock_mod.MockModule {
	d := mock_mod.NewMockModule(ctrl)
	d.EXPECT().Version().Return(semver.New(v)).Times(oneTime)
	return d
}

func newErrTip(err error) []Tip {
	return []Tip{&tip{err: err}}
}

func newMockMod(ctrl *gomock.Controller, deps []mod.Module) *mock_mod.MockMod {
	m := mock_mod.NewMockMod(ctrl)
	m.EXPECT().Dependencies().Return(deps).Times(oneTime)
	return m
}

func newTip(msg string) []Tip {
	return []Tip{&tip{msg: msg}}
}

func newTag(ctrl *gomock.Controller, v string) *mock_mod.MockModule {
	d := mock_mod.NewMockModule(ctrl)
	d.EXPECT().Path().Return(repoName).Times(oneTime)
	d.EXPECT().Version().Return(semver.New(v)).AnyTimes()
	return d
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/internal/semver"
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
	const v = "v1.0.0"
	var (
		sy1 = newMockSystem(ctrl, semver.Tags{semver.New(v)}, nil)
		are = is.New(t)
		dt  = map[string]struct {
			ctx  context.Context
			file mod.Mod
			cnf  Config
			res  []Tip
		}{
			"default":       {res: newErrTip(errors.ErrMod)},
			"no context":    {file: mock_mod.NewMockMod(ctrl), res: newErrTip(errors.ErrMod)},
			"no file":       {ctx: ctx, res: newErrTip(errors.ErrMod)},
			"no dependency": {ctx: ctx, file: newMockMod(ctrl, nil)},
			"all good": {
				ctx: ctx,
				file: newMockMod(ctrl, []mod.Module{
					newMockModule(ctrl, "gitlab.lan/group/pkg", "v1.0.0", false),
				}),
				res: newTip("gitlab.lan/group/pkg v1.0.0 is up to date"),
			},
		}
		u = New(Config{}, SetGoGet(sy1))
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			res := u.CheckFile(tt.ctx, tt.file, tt.cnf)
			are.Equal(res, tt.res) // mismatch result
		})
	}
}

func newMockSystem(ctrl *gomock.Controller, tags semver.Tags, err error) *mock_vcs.MockSystem {
	m := mock_vcs.NewMockSystem(ctrl)
	m.EXPECT().CanFetch(gomock.Any()).Return(true).AnyTimes()
	m.EXPECT().FetchPath(gomock.Any(), gomock.Any()).Return(tags, err).AnyTimes()
	return m
}

func newMockModule(ctrl *gomock.Controller, path, version string, indirect bool) *mock_mod.MockModule {
	m := mock_mod.NewMockModule(ctrl)
	m.EXPECT().Path().Return(path).AnyTimes()
	m.EXPECT().Version().Return(semver.New(version)).AnyTimes()
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
			"default": {dep: mock_mod.NewMockModule(ctrl)},
			"invalid glob pattern": {
				dep:   newTag(ctrl, "v1.0.0-b42", 1),
				paths: "example.com",
			},
			"valid glob pattern": {
				dep:   newTag(ctrl, "v1.0.0-b42", 1),
				paths: "example.com/*/*",
				err:   errors.ErrExpectedTag,
			},
			"skip": {dep: newTag(ctrl, "v1.0.0", 2), paths: "test,,example.com"},
			"ok":   {dep: newTag(ctrl, "v1.0.0", 1), paths: "example.com/pkg/*"},
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

func newTag(ctrl *gomock.Controller, v string, times int) *mock_mod.MockModule {
	d := mock_mod.NewMockModule(ctrl)
	d.EXPECT().Path().Return("example.com/pkg/go").Times(times)
	d.EXPECT().Version().Return(semver.New(v)).AnyTimes()
	return d
}

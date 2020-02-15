// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/internal/semver"
	mock_mod "github.com/rvflash/goup/testdata/mock/mod"
)

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
			"only": {
				dep:   newTag(ctrl, "example.com/pkg/go", "v1.0.0-b42", 1),
				paths: "example.com",
				err:   errors.ErrExpectedTag,
			},
			"skip": {dep: newTag(ctrl, "example.com/pkg/go", "v1.0.0", 2), paths: "test,,example.com"},
			"ok":   {dep: newTag(ctrl, "example.com/pkg/go", "v1.0.0", 1), paths: "example.com"},
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

func newTag(ctrl *gomock.Controller, path, v string, times int) *mock_mod.MockModule {
	d := mock_mod.NewMockModule(ctrl)
	d.EXPECT().Path().Return(path).Times(times)
	d.EXPECT().Version().Return(semver.New(v)).AnyTimes()
	return d
}

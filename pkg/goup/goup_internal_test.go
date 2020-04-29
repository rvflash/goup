// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"

	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/pkg/mod"
	mockMod "github.com/rvflash/goup/testdata/mock/mod"
)

func TestUpdateFile(t *testing.T) {
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
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			err := updateFile(tt.file)
			are.True(errors.Is(err, tt.err))                  // mismatch error
			are.Equal(tt.updated, fileExists(tt.file.Name())) // mismatch file "created"
		})
	}
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
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, ok := latest(tt.in, tt.dep, tt.cnf.Major, tt.cnf.MajorMinor)
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

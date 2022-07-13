// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE dep.

package goup

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/pkg/mod"
	mockMod "github.com/rvflash/goup/testdata/mock/mod"
)

const (
	v0       = "v0.0.0"
	v1       = "v0.0.1"
	repoName = "example.com/group/go"
	oneTime  = 1
)

func TestNewEntry(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Error("expected no panic")
		}
	}()
	var e *Entry
	e.Args()
	e.Format()
	e.Level()
	e.OutDated()
}

func TestNewCheck(t *testing.T) {
	t.Parallel()
	var (
		dep  mod.Module
		are  = is.New(t)
		ctrl = gomock.NewController(t)
	)
	defer ctrl.Finish()

	are.Equal(newCheck(dep), nil) // mismatch default
	dep = newDep(ctrl)
	msg := newCheck(dep)
	are.Equal(msg.Level(), DebugLevel)                        // mismatch level
	are.True(strings.Contains(msg.Format(), "is up to date")) // mismatch message
	are.Equal(len(msg.Args()), 2)                             // expected dep and version
	_, ok := msg.OutDated()
	are.True(!ok) // not outdated
}

func TestNewError(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		ctrl = gomock.NewController(t)
		dt   = map[string]struct {
			// inputs
			err  error
			file mod.Mod
			// outputs
			msg string
			len int
		}{
			"Default":       {},
			"Without error": {file: &mockMod.MockMod{}},
			"Without file":  {err: errors.ErrExpectedTag},
			"Complete": {
				file: newMod(ctrl),
				err:  errors.ErrExpectedTag,
				msg:  errors.ErrExpectedTag.Error(),
				len:  1,
			},
		}
	)
	defer ctrl.Finish()

	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			msg := newError(tt.err, tt.file)
			are.Equal(msg.Level(), ErrorLevel)               // mismatch level
			are.True(strings.Contains(msg.Format(), tt.msg)) // mismatch message
			are.Equal(len(msg.Args()), tt.len)               // mismatch len
			_, ok := msg.OutDated()
			are.True(!ok) // not outdated
		})
	}
}

func TestNewFailure(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		ctrl = gomock.NewController(t)
		dt   = map[string]struct {
			// inputs
			err error
			dep mod.Module
			// outputs
			msg string
			len int
		}{
			"Default":       {},
			"Without error": {dep: &mockMod.MockModule{}},
			"Without dep":   {err: errors.ErrExpectedTag},
			"Complete": {
				dep: newDep(ctrl),
				err: errors.ErrExpectedTag,
				msg: "check failed",
				len: 2,
			},
		}
	)
	defer ctrl.Finish()

	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			msg := newFailure(tt.err, tt.dep)
			are.Equal(msg.Level(), ErrorLevel)               // mismatch level
			are.True(strings.Contains(msg.Format(), tt.msg)) // mismatch message
			are.Equal(len(msg.Args()), tt.len)               // mismatch len
			_, ok := msg.OutDated()
			are.True(!ok) // not outdated
		})
	}
}

func TestNewSkip(t *testing.T) {
	t.Parallel()
	var (
		dep  mod.Module
		are  = is.New(t)
		ctrl = gomock.NewController(t)
	)
	defer ctrl.Finish()

	are.Equal(newSkip(dep), nil) // mismatch default
	dep = newDep(ctrl)
	msg := newSkip(dep)
	are.Equal(msg.Level(), DebugLevel)                         // mismatch level
	are.True(strings.Contains(msg.Format(), "update skipped")) // mismatch message
	are.Equal(len(msg.Args()), 2)                              // expected dep and version
	_, ok := msg.OutDated()
	are.True(!ok) // not outdated
}

func TestNewUpdate(t *testing.T) {
	t.Parallel()
	var (
		dep  mod.Module
		are  = is.New(t)
		ctrl = gomock.NewController(t)
	)
	defer ctrl.Finish()

	are.Equal(newUpdate(dep, v1), nil) // mismatch default
	dep = newDep(ctrl)
	msg := newUpdate(dep, v1)
	are.Equal(msg.Level(), InfoLevel)                           // mismatch level
	are.True(strings.Contains(msg.Format(), "will be updated")) // mismatch message
	are.Equal(len(msg.Args()), 3)                               // expected dep, old and new versions
	_, ok := msg.OutDated()
	are.True(!ok) // not outdated
}

func TestNewOutOfDate(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		ctrl = gomock.NewController(t)
	)
	defer ctrl.Finish()

	are.Equal(newOutOfDate(nil, v1), nil) // mismatch default
	var (
		dep = newDep(ctrl)
		msg = newOutOfDate(dep, v1)
	)
	are.Equal(msg.Level(), WarnLevel)                           // mismatch level
	are.True(strings.Contains(msg.Format(), "must be updated")) // mismatch message
	are.Equal(len(msg.Args()), 3)                               // expected dep, old and new versions
	v, ok := msg.OutDated()
	are.True(ok)     // outdated
	are.Equal(v, v1) // new version mismatch
}

func newMod(ctrl *gomock.Controller) *mockMod.MockMod {
	m := mockMod.NewMockMod(ctrl)
	m.EXPECT().Module().Return(repoName).Times(oneTime)
	return m
}

func newDep(ctrl *gomock.Controller) *mockMod.MockModule {
	d := mockMod.NewMockModule(ctrl)
	d.EXPECT().Path().Return(repoName).Times(oneTime)
	d.EXPECT().Version().Return(semver.New(v0)).AnyTimes()
	return d
}

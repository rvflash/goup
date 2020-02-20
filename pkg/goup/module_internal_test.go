// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"

	"github.com/rvflash/goup/internal/mod"
)

const (
	release  = "v1.0.2"
	repoName = "example.com/group/go"
)

func TestNewError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		err = errors.New("oops")
		are = is.New(t)
		dt  = map[string]struct {
			mod     mod.Module
			in, out error
		}{
			"default":       {},
			"missing error": {mod: newMockModule(ctrl, false)},
			"missing mod":   {in: err},
			"ok":            {mod: newMockModule(ctrl, false), in: err, out: err},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := newError(tt.mod, tt.in)
			are.True(errors.Is(out, tt.out)) // mismatch result
			if tt.out != nil {
				are.True(strings.Contains(out.Error(), tt.mod.Path())) // mismatch path
			}
		})
	}
}

func TestNewOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			mod  mod.Module
			in   string
			fail bool
		}{
			"default":         {},
			"missing message": {mod: newMockModule(ctrl, false), fail: true},
			"ok":              {mod: newMockModule(ctrl, false), in: release, fail: true},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := newOrder(tt.mod, tt.in)
			are.Equal(out != nil, tt.fail) // mismatch error
			if out != nil {
				are.True(strings.Contains(out.Error(), tt.mod.Path()))             // mismatch path
				are.True(strings.Contains(out.Error(), tt.mod.Version().String())) // mismatch version
				are.True(strings.Contains(out.Error(), tt.in))                     // mismatch release
			}
		})
	}
}

func TestChecked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		m   mod.Module
		are = is.New(t)
	)
	are.Equal(checked(m), "") // mismatch default
	m = newMockModule(ctrl, false)
	are.Equal(checked(m), "example.com/group/go v1.0.2 is up to date") // mismatch result
}

func TestSkipped(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		m   mod.Module
		are = is.New(t)
	)
	are.Equal(skipped(m), "") // mismatch default
	m = newMockModule(ctrl, false)
	are.Equal(skipped(m), "example.com/group/go v1.0.2 update skipped") // mismatch result
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package semver_test

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/semver"
)

type tag string

func (t tag) String() string {
	return string(t)
}

func TestCompare(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			v, w fmt.Stringer
			out  int
		}{
			"default": {},
			">":       {v: tag("v1.2.3"), w: tag("v1.2.2"), out: 1},
			"<":       {v: tag("v0.2.3"), w: tag("v1.2.2"), out: -1},
			"<=":      {v: tag("v0.2.3-12"), w: tag("v0.2.3"), out: -1},
			"=":       {v: tag("v0.2.3+12"), w: tag("v0.2.3")},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := semver.Compare(tt.v, tt.w)
			are.Equal(out, tt.out) // mismatch result
		})
	}
}

func TestLatest(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  semver.Tags
			out semver.Tag
		}{
			"default": {},
			"one":     {in: semver.Tags{v3}, out: v3},
			"latest":  {in: tags(), out: v0},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := semver.Latest(tt.in)
			are.Equal(out, tt.out) // mismatch result
		})
	}
}

func TestLatestMinor(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in    semver.Tags
			major string
			out   semver.Tag
		}{
			"default":   {},
			"undefined": {in: tags()},
			"unknown":   {in: tags(), major: "v8"},
			"v0":        {in: tags(), major: "v0", out: v4},
			"v1":        {in: tags(), major: "v1", out: v5},
			"v2":        {in: tags(), major: "v2", out: v0},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := semver.LatestMinor(tt.major, tt.in)
			are.Equal(out, tt.out) // mismatch result
		})
	}
}

func TestLatestPatch(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in         semver.Tags
			majorMinor string
			out        semver.Tag
		}{
			"default":   {},
			"undefined": {in: tags()},
			"unknown":   {in: tags(), majorMinor: "v0.1"},
			"v0":        {in: tags(), majorMinor: "v0.2", out: v4},
			"v1":        {in: tags(), majorMinor: "v1.2", out: v5},
			"v2":        {in: tags(), majorMinor: "v2.2", out: v0},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := semver.LatestPatch(tt.majorMinor, tt.in)
			are.Equal(out, tt.out) // mismatch result
		})
	}
}

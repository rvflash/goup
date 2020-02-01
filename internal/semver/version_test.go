// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package semver_test

import (
	"sort"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/semver"
)

var (
	v0 = semver.New("v2.2.13")
	v1 = semver.New("v2.2.12-beta")
	v2 = semver.New("v2.2.12")
	v3 = semver.New("v1.2.12+incompatible")
	v4 = semver.New("v0.2.12")
	v5 = semver.New("v1.2.13")
)

func tags() semver.Tags {
	return semver.Tags{v0, v1, v2, v3, v4, v5}
}

func TestNew(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in         string
			out        string
			canonical  string
			tag        bool
			valid      bool
			major      string
			majorMinor string
			prerelease string
			build      string
		}{
			"1.2.0": {in: "1.2.0"},
			"gopls/v1.2.12": {
				in:         "gopls/v1.2.12",
				out:        "gopls/v1.2.12",
				canonical:  "v1.2.12",
				major:      "v1",
				majorMinor: "v1.2",
				valid:      true,
				tag:        true,
			},
			"v1.2.12": {
				in:         "v1.2.12",
				out:        "v1.2.12",
				canonical:  "v1.2.12",
				major:      "v1",
				majorMinor: "v1.2",
				valid:      true,
				tag:        true,
			},
			"v0.2.1-0.20200121190230-accd165b1659": {
				in:         "v0.2.1-0.20200121190230-accd165b1659",
				out:        "v0.2.1-0.20200121190230-accd165b1659",
				canonical:  "v0.2.1-0.20200121190230-accd165b1659",
				major:      "v0",
				majorMinor: "v0.2",
				prerelease: "-0.20200121190230-accd165b1659",
				valid:      true,
			},
			"v1.3.14+1904": {
				in:         "v1.3.14+1904",
				out:        "v1.3.14+1904",
				canonical:  "v1.3.14",
				major:      "v1",
				majorMinor: "v1.3",
				build:      "+1904",
				valid:      true,
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := semver.New(tt.in)
			are.Equal(out.IsTag(), tt.tag)             // mismatch tag
			are.Equal(out.IsValid(), tt.valid)         // mismatch validity
			are.Equal(out.Major(), tt.major)           // mismatch major
			are.Equal(out.MajorMinor(), tt.majorMinor) // mismatch major minor
			are.Equal(out.Canonical(), tt.canonical)   // mismatch canonical
			are.Equal(out.Prerelease(), tt.prerelease) // mismatch prerelease
			are.Equal(out.Build(), tt.build)           // mismatch build
			are.Equal(out.String(), tt.out)            // mismatch raw
		})
	}
}

func TestTags_Len(t *testing.T) {
	var (
		are  = is.New(t)
		list semver.Tags
	)
	list = tags()
	sort.Sort(list)
	are.Equal(list.Len(), 6)
	are.Equal(list, semver.Tags{v4, v3, v5, v1, v2, v0})
}

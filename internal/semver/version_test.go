// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package semver_test

import (
	"testing"

	"github.com/matryer/is"

	"github.com/rvflash/goup/internal/semver"
)

func TestNew(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in         string
			out        string
			canonical  string
			valid      bool
			major      string
			majorMinor string
			prerelease string
			build      string
		}{
			"1.2.0": {in: "1.2.0"},
			"v1.2.0": {
				in:         "v1.2.0",
				out:        "v1.2.0",
				canonical:  "v1.2.0",
				major:      "v1",
				majorMinor: "v1.2",
				valid:      true,
			},
			"v1.2.12": {
				in:         "v1.2.12",
				out:        "v1.2.12",
				canonical:  "v1.2.12",
				major:      "v1",
				majorMinor: "v1.2",
				valid:      true,
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

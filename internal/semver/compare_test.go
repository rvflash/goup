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

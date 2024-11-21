// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod

import (
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/semver"
)

const (
	indirect = true
	path     = "github.com/rvflash/goup"
	version  = "v1.0.1"
)

func TestModule_Version(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		var (
			mod = module{}
			are = is.New(t)
		)
		are.Equal(mod.Indirect(), false)         // mismatch indirect
		are.Equal(mod.Path(), "")                // mismatch path
		are.True(!mod.Replacement())             // mismatch replacement
		are.Equal(mod.Version(), nil)            // mismatch version
		are.Equal(0, len(mod.ExcludeVersions())) // unexpected exclude version
	})

	t.Run("valued", func(t *testing.T) {
		t.Parallel()
		var (
			v   = semver.New(version)
			mod = module{
				excludes:    []semver.Tag{v},
				indirect:    indirect,
				path:        path,
				replacement: true,
				version:     v,
			}
			are = is.New(t)
		)
		x := mod.ExcludeVersions()
		are.Equal(1, len(x))                // expected exclude version
		are.Equal(x[0], v)                  // mismatch exclude version
		are.Equal(mod.Indirect(), indirect) // mismatch indirect
		are.Equal(mod.Path(), path)         // mismatch path
		are.True(mod.Replacement())         // mismatch replacement
		are.Equal(mod.Version(), v)         // mismatch version
	})
}

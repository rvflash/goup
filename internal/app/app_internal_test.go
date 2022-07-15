// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

const (
	root    = "/rv"
	modFile = root + "/go.mod"
	numFile = 2
)

func TestFilePath(t *testing.T) {
	t.Parallel()
	are := is.New(t)
	are.Equal(filepath.ToSlash(filePath(root)), modFile)    // mismatch root only
	are.Equal(filepath.ToSlash(filePath(modFile)), modFile) // mismatch path
}

func TestWalkPath(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		res  []string
		bad  = []string{"..", "..", "testdata", "not_exists"}
		tree = []string{"..", "..", "testdata", "golden", "tree"}
	)
	res = walkPath(".")
	are.True(len(res) == 0) // unexpected result
	res = walkPath(filepath.Join(bad...))
	are.True(len(res) == 1)                                   // mismatch not found
	are.Equal(res[0], filepath.Join(bad...))                  // mismatch not found result
	are.Equal(len(walkPath(filepath.Join(tree...))), numFile) // mismatch result
}

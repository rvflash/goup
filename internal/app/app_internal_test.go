// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"testing"

	"github.com/matryer/is"
)

const (
	root    = "/rv"
	modFile = "/rv/go.mod"
	badPath = "../../testdata/not_found"
	treeDir = "../../testdata/gomod/tree"
	numFile = 2
)

func TestFilePath(t *testing.T) {
	are := is.New(t)
	are.Equal(filePath(root), modFile)    // mismatch root only
	are.Equal(filePath(modFile), modFile) // mismatch path
}

func TestWalkPath(t *testing.T) {
	var (
		are = is.New(t)
		res []string
	)
	res = walkPath(".")
	are.True(len(res) == 0) // unexpected result
	res = walkPath(badPath)
	are.True(len(res) == 1)                    // mismatch not found
	are.Equal(res[0], badPath)                 // mismatch not found result
	are.Equal(len(walkPath(treeDir)), numFile) // mismatch result
}

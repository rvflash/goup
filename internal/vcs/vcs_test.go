// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs_test

import (
	"net/url"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/vcs"
)

func TestRepoPath(t *testing.T) {
	var (
		are = is.New(t)
		uri *url.URL
	)
	are.Equal(vcs.RepoPath(uri), "") // mismatch default

	var err error
	uri, err = url.Parse("https://google.golang.org/appengine?go-get=1")
	are.NoErr(err)
	are.Equal(vcs.RepoPath(uri), "google.golang.org/appengine") // mismatch result
}

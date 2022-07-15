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

func TestIsSecureScheme(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  string
			out bool
		}{
			"default": {},
			"unknown": {in: "oops"},
			"http":    {in: vcs.HTTP},
			"git":     {in: vcs.Git},
			"https":   {in: vcs.HTTPS, out: true},
			"ssh+git": {in: vcs.SSHGit, out: true},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			are.Equal(vcs.IsSecureScheme(tt.in), tt.out) // mismatch result
		})
	}
}

func TestRepoPath(t *testing.T) {
	t.Parallel()
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

func TestURLScheme(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  string
			out string
		}{
			"default": {},
			"unknown": {in: "oops"},
			"git":     {in: vcs.Git, out: "git://"},
			"http":    {in: vcs.HTTP, out: "http://"},
			"https":   {in: vcs.HTTPS, out: "https://"},
			"ssh+git": {in: vcs.SSHGit, out: "ssh://git@"},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			are.Equal(vcs.URLScheme(tt.in), tt.out) // mismatch result
		})
	}
}

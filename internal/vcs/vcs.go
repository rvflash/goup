// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package vcs exposes interfaces and methods implemented by a VCS.
package vcs

//go:generate mockgen -destination ../../testdata/mock/vcs/vcs.go -source vcs.go

import (
	"context"
	"net/http"
	"net/url"

	"github.com/rvflash/goup/internal/semver"
)

// System must be implemented by any VCS.
type System interface {
	CanFetch(path string) bool
	FetchPath(ctx context.Context, path string) (semver.Tags, error)
	FetchURL(ctx context.Context, url string) (semver.Tags, error)
}

// HTTPClient must be implemented by any HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientChooser must be implemented to return a HTTP Client for the given rawURL.
type ClientChooser interface {
	ClientFor(path string) HTTPClient
}

// RepoPath returns the path of the repository based on its URL.
func RepoPath(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.Host + u.Path
}

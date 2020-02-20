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

// Client must be implemented by any HTTP client.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientChooser must be implemented to return a HTTPClient Client for the given rawURL.
type ClientChooser interface {
	ClientFor(path string) Client
	AllowInsecure(path string) bool
}

// List of supported schemes
const (
	HTTPS  = "https"
	HTTP   = "http"
	SSHGit = "ssh+git"
	Git    = "git"
)

// IsSecureScheme returns true if the given scheme is marked as secure.
func IsSecureScheme(s string) bool {
	switch s {
	case HTTPS, SSHGit:
		return true
	default:
		return false
	}
}

// RepoPath returns the path of the repository based on its URL.
func RepoPath(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.Host + u.Path
}

// URLScheme returns the protocol scheme to use as prefix path for this scheme.
func URLScheme(name string) string {
	switch name {
	case HTTPS, HTTP, Git:
		return name + "://"
	case SSHGit:
		return "ssh://git@"
	default:
		return ""
	}
}

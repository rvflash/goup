// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package vcs exposes interfaces and methods implemented by a VCS.
package vcs

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/rvflash/goup/internal/semver"
)

// System must be implemented by any VCS.
type System interface {
	CanFetch(path string) bool
	FetchPath(ctx context.Context, path string) (semver.Tags, error)
	FetchURL(ctx context.Context, url string) (semver.Tags, error)
}

// HTTPClient represents a HTTP Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	keepAlive       = 30 * time.Second
	continueTimeout = 1 * time.Second
)

// NewHTTPClient returns a new instance of HTTP client.
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: keepAlive,
			}).DialContext,
			TLSHandshakeTimeout:   timeout,
			ResponseHeaderTimeout: timeout,
			ExpectContinueTimeout: continueTimeout,
		},
	}
}

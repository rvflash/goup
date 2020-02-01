// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/rvflash/goup/internal/semver"
)

// System
type System interface {
	CanFetch(path string) bool
	FetchPath(ctx context.Context, path string) (semver.Tags, error)
}

// HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTPClient
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   timeout,
			ResponseHeaderTimeout: timeout,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

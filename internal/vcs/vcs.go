// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package vcs exposes interfaces and methods implemented by a VCS.
package vcs

//go:generate mockgen -destination ../../testdata/mock/vcs/vcs.go -source vcs.go

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

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

// ClientChooser must be implemented to return a HTTP Client for the given URL.
type ClientChooser interface {
	ClientFor(path string) HTTPClient
}

const (
	comma = ","
	https = "https"
)

// NewHTTPClient creates a new instance of Client.
func NewHTTPClient(timeout time.Duration, goInsecure string) *HTTP {
	skipSec := strings.Split(goInsecure, comma)
	return &HTTP{
		insecure: newInsecureHTTPClient(timeout),
		secure:   newSecureHTTPClient(timeout),
		skipSec:  skipSec,
	}
}

// HTTP allows to communicate over HTTP or HTTPS.
type HTTP struct {
	secure,
	insecure *http.Client
	skipSec []string
}

// ClientFor returns the HTTP client to use for this URL.
func (c *HTTP) ClientFor(path string) HTTPClient {
	if c.AllowInsecure(path) {
		return c.insecure
	}
	return c.secure
}

// AllowInsecure returns true if this URL allows insecure request.
func (c *HTTP) AllowInsecure(name string) bool {
	if name == "" {
		return false
	}
	var matched bool
	for _, pattern := range c.skipSec {
		if matched, _ = path.Match(pattern, name); matched {
			return true
		}
	}
	return false
}

// newInsecureHTTPClient returns a HTTP client that allows plain HTTP and skips HTTPS validation.
func newInsecureHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// newSecureHTTPClient returns a HTTP client that rejects redirection from HTTPS to HTTP and validate HTTPS.
func newSecureHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 && via[0].URL.Scheme == https && req.URL.Scheme != https {
				lastURL := via[len(via)-1].URL
				return fmt.Errorf("redirected from secure URL %s to insecure URL %s", lastURL, req.URL)
			}
			return nil
		},
	}
}

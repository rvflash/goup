// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

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

// ClientFor returns the HTTP client to use for this rawURL.
func (c *HTTP) ClientFor(path string) HTTPClient {
	if c.AllowInsecure(path) {
		return c.insecure
	}
	return c.secure
}

// AllowInsecure returns true if this rawURL allows insecure request.
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
				return fmt.Errorf("redirected from secure rawURL %s to insecure rawURL %s", lastURL, req.URL)
			}
			return nil
		},
	}
}

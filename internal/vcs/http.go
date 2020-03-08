// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/rvflash/goup/internal/path"
)

const https = "https"

// NewHTTPClient creates a new instance of Client.
func NewHTTPClient(timeout time.Duration, insecurePaths string) *HTTPClient {
	return &HTTPClient{
		insecure:      newInsecureHTTPClient(timeout),
		secure:        newSecureHTTPClient(timeout),
		insecurePaths: insecurePaths,
	}
}

// HTTPClient allows to communicate over HTTPClient or HTTPS.
type HTTPClient struct {
	secure,
	insecure *http.Client
	insecurePaths string
}

// ClientFor returns the HTTPClient client to use for this rawURL.
func (c *HTTPClient) ClientFor(path string) Client {
	if c.AllowInsecure(path) {
		return c.insecure
	}
	return c.secure
}

// AllowInsecure returns true if this rawURL allows insecure request.
func (c *HTTPClient) AllowInsecure(target string) bool {
	return path.Match(c.insecurePaths, target)
}

// newInsecureHTTPClient returns a HTTPClient client that allows plain HTTPClient and skips HTTPS validation.
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

// newSecureHTTPClient returns a HTTPClient client that rejects redirection from HTTPS to HTTPClient and validate HTTPS.
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

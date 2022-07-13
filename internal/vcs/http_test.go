// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/vcs"
)

const repo = "google.golang.org/grpc"

func TestHTTP_AllowInsecure(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			repo       string
			goInsecure string
			out        bool
		}{
			"default":        {},
			"secure":         {repo: repo},
			"still secure":   {repo: repo, goInsecure: "github.com/*"},
			"insecure like":  {repo: repo, goInsecure: "google.golang.org/*", out: true},
			"insecure match": {repo: repo, goInsecure: repo, out: true},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cli := vcs.NewHTTPClient(time.Second, tt.goInsecure)
			are.Equal(cli.AllowInsecure(tt.repo), tt.out)                             // mismatch result
			are.Equal(cli.ClientFor(tt.repo).(*http.Client).Transport != nil, tt.out) // mismatch client
		})
	}
}

func TestHTTP_ClientFor(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		req *http.Request
		via []*http.Request
	)
	// Default security (not insecure but not secure only)
	cli := vcs.NewHTTPClient(time.Second, "")

	// Default
	are.NoErr(cli.ClientFor(repo).(*http.Client).CheckRedirect(req, via)) // expected no error

	// HTTPS > HTTPS
	var err error
	req, err = http.NewRequest("GET", "https://"+repo, nil)
	are.NoErr(err)
	var req2 *http.Request
	req2, err = http.NewRequest("GET", "https://"+repo+"/new", nil)
	are.NoErr(err)
	are.NoErr(cli.ClientFor(repo).(*http.Client).CheckRedirect(req, []*http.Request{req2})) // expected no error

	// HTTPS > HTTPClient
	req2, err = http.NewRequest("GET", "http://"+repo+"/new", nil)
	are.NoErr(err)
	are.True(cli.ClientFor(repo).(*http.Client).CheckRedirect(req2, []*http.Request{req}) != nil) // expected error
}

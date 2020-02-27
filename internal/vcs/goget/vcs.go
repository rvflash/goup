// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package goget provides methods to deal with go get as VCS.
package goget

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/git"
)

// Name is the name of this VCS.
const Name = "go-get"

// VCS is a go-get version control system.
// We use go-get to retrieve the remote's properties behind a package.
type VCS struct {
	http vcs.ClientChooser
	git  vcs.System
}

// New creates a new instance of VCS.
func New(client vcs.ClientChooser, git vcs.System) *VCS {
	return &VCS{
		http: client,
		git:  git,
	}
}

// CanFetch implements the vcs.VCS interface.
func (s *VCS) CanFetch(path string) bool {
	if path == "" {
		return false
	}
	for _, service := range []string{"github.com", "gitlab", "bitbucket"} {
		if strings.Contains(path, service) {
			return false
		}
	}
	return true
}

// FetchPath implements the vcs.VCS interface.
func (s *VCS) FetchPath(ctx context.Context, path string) (semver.Tags, error) {
	system, remote, err := s.vcsByPath(ctx, path)
	if err != nil {
		return nil, err
	}
	return s.fetchURL(ctx, system, remote)
}

// FetchURL implements the vcs.VCS interface.
func (s *VCS) FetchURL(ctx context.Context, url string) (semver.Tags, error) {
	system, remote, err := s.vcsByURL(ctx, url)
	if err != nil {
		return nil, err
	}
	return s.fetchURL(ctx, system, remote)
}

func (s *VCS) fetchURL(ctx context.Context, system, url string) (semver.Tags, error) {
	if s.git == nil {
		return nil, errors.ErrSystem
	}
	switch system {
	case git.Name:
		return s.git.FetchURL(ctx, url)
	default:
		return nil, vcs.Errorf(system, errors.ErrSystem)
	}
}

func (s *VCS) vcsByPath(ctx context.Context, path string) (name, remote string, err error) {
	if path == "" {
		return "", "", errors.ErrRepository
	}
	for _, scheme := range []string{vcs.HTTPS, vcs.HTTP} {
		name, remote, err = s.vcsByURL(ctx, vcs.URLScheme(scheme)+path)
		if err == nil {
			break
		}
	}
	return
}

func (s *VCS) vcsByURL(ctx context.Context, url string) (name, remote string, err error) {
	if ctx == nil || s.http == nil {
		return "", "", errors.ErrSystem
	}
	if url == "" {
		return "", "", errors.ErrRepository
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	setQuery(req.URL)

	// Security check
	if !vcs.IsSecureScheme(req.URL.Scheme) && !s.http.AllowInsecure(vcs.RepoPath(req.URL)) {
		return "", "", errors.NewSecurityIssue(req.URL.String())
	}
	var resp *http.Response
	resp, err = s.http.ClientFor(vcs.RepoPath(req.URL)).Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	return parseMetaGoImport(resp.Body)
}

const (
	body    = "body"
	head    = "head"
	meta    = "meta"
	name    = "name"
	attr    = "go-import"
	content = "content"

	// supported charsets
	utf8  = "utf-8"
	ascii = "ascii"
)

func parseMetaGoImport(r io.Reader) (vcs, url string, err error) {
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	d.Strict = false
	var t xml.Token
	for {
		t, err = d.RawToken()
		if err != nil {
			break
		}
		if e, ok := t.(xml.StartElement); ok && strings.EqualFold(e.Name.Local, body) {
			break
		}
		if e, ok := t.(xml.EndElement); ok && strings.EqualFold(e.Name.Local, head) {
			break
		}
		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, meta) {
			continue
		}
		if attrValue(e.Attr, name) != attr {
			continue
		}
		if f := strings.Fields(attrValue(e.Attr, content)); len(f) == 3 {
			vcs = f[1]
			url = f[2]
			break
		}
	}
	return
}

func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case utf8, ascii:
		return input, nil
	default:
		return nil, vcs.Errorf(Name, errors.NewCharset(charset))
	}
}

const enabled = "1"

func setQuery(u *url.URL) {
	if u == nil {
		return
	}
	values, _ := url.ParseQuery(u.RawQuery)
	values.Set(Name, enabled)
	u.RawQuery = values.Encode()
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package git provides methods to handle Git.
package git

import (
	"context"
	"net/url"
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
)

const (
	// Name is the name of this VCS.
	Name = "git"
	// Ext of a git repository.
	Ext = ".git"
)

// VCS is a Git PrintVersion Control VCS.
type VCS struct {
	client  vcs.ClientChooser
	storage storage.Storer
}

// New returns a new instance of VCS.
func New(client vcs.ClientChooser) *VCS {
	return &VCS{
		client:  client,
		storage: memory.NewStorage(),
	}
}

// CanFetch implements the vcs.VCS interface.
func (s *VCS) CanFetch(_ string) bool {
	return true
}

const oneRef = 1

// FetchPath implements the vcs.VCS interface.
func (s *VCS) FetchPath(ctx context.Context, path string) (semver.Tags, error) {
	if !s.ready(ctx) {
		return nil, errors.ErrSystem
	}
	if path == "" {
		return nil, errors.ErrRepository
	}
	var c = make(chan *reference, oneRef)
	go func() {
		c <- s.fetchWithRetry(path)
	}()
	return tags(ctx, c)
}

// FetchURL implements the vcs.VCS interface.
func (s *VCS) FetchURL(ctx context.Context, url string) (semver.Tags, error) {
	if !s.ready(ctx) {
		return nil, errors.ErrSystem
	}
	var c = make(chan *reference, oneRef)
	go func() {
		c <- s.fetch(url)
	}()
	return tags(ctx, c)
}

func (s *VCS) fetchWithRetry(path string) (ref *reference) {
	for _, t := range []transport{
		// Secure
		{scheme: vcs.HTTPS},
		{scheme: vcs.SSHGit, extension: Ext},
		// Insecure
		{scheme: vcs.Git, extension: Ext},
		{scheme: vcs.HTTP},
	} {
		ref = s.fetch(t.rawURL(path))
		if ref.err == nil {
			break
		}
	}
	return
}

func (s *VCS) fetch(rawURL string) *reference {
	ref := new(reference)
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		ref.err = vcs.Errorf(Name, errors.ErrRepository, err)
		return ref
	}
	// Security check
	if !vcs.IsSecureScheme(u.Scheme) && !s.client.AllowInsecure(vcs.RepoPath(u)) {
		ref.err = vcs.Errorf(Name, errors.ErrRepository, errors.NewSecurityIssue(u.String()))
		return ref
	}
	rem := git.NewRemote(s.storage, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{u.String()},
	})
	// Retrieves the releases list of the repository.
	var res []*plumbing.Reference
	res, err = rem.List(&git.ListOptions{})
	if err != nil {
		ref.err = vcs.Errorf(Name, errors.ErrFetch, err)
		return ref
	}
	// Filters to keep only tag.
	var n plumbing.ReferenceName
	for _, r := range res {
		n = r.Name()
		if n.IsTag() {
			ref.list = append(ref.list, semver.New(n.Short()))
		}
	}
	return ref
}

func (s *VCS) ready(ctx context.Context) bool {
	return ctx != nil && s.storage != nil && s.client != nil
}

type reference struct {
	list semver.Tags
	err  error
}

func tags(ctx context.Context, c chan *reference) (semver.Tags, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ref := <-c:
		if err := ref.err; err != nil {
			return nil, err
		}
		return ref.list, nil
	}
}

const (
	// example.com/group/pkg, so with 2 slashes: 3 parts
	stdNumPart = 3
	slash      = "/"
)

type transport struct {
	scheme    string
	extension string
}

func (t transport) rawURL(uri string) string {
	p := strings.Split(uri, slash)
	if len(p) > stdNumPart {
		// Works around with sub-packages.
		uri = path.Join(p[:stdNumPart]...)
	}
	return vcs.URLScheme(t.scheme) + uri + t.extension
}

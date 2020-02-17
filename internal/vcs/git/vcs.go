// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package git provides methods to handle Git.
package git

import (
	"context"
	"net/url"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
)

// Name is the name of this VCS.
const Name = "git"

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
		{protocol: "https://"},
		{protocol: "http://"},
	} {
		ref = s.fetch(t.rawURL(path))
		if ref.err == nil {
			break
		}
	}
	return
}

func (s *VCS) fetch(rawURL string) (ref *reference) {
	ref = new(reference)
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		ref.err = vcs.Errorf(Name, errors.ErrRepository, err)
		return
	}
	// Override http(s) default protocol to use one dedicated to this package (insecure?).
	// For now, it is not possible with git, the HTTP client is embedded and global set :/
	rem := git.NewRemote(s.storage, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{u.String()},
	})
	// Retrieves the releases list of the repository.
	var res []*plumbing.Reference
	res, err = rem.List(&git.ListOptions{})
	if err != nil {
		ref.err = vcs.Errorf(Name, errors.ErrFetch, err)
		return
	}
	// Filters to keep only tag.
	var n plumbing.ReferenceName
	for _, r := range res {
		n = r.Name()
		if n.IsTag() {
			ref.list = append(ref.list, semver.New(n.Short()))
		}
	}
	return
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

type transport struct {
	protocol  string
	extension string
}

func (t transport) rawURL(path string) string {
	return t.protocol + path + t.extension
}

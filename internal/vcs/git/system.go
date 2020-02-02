// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package git provides methods to handle Git.
package git

import (
	"context"

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

// System is a Git Version Control System.
type System struct {
	storage storage.Storer
}

type reference struct {
	list semver.Tags
	err  error
}

// New returns a new instance of System.
func New() *System {
	return &System{storage: memory.NewStorage()}
}

// CanFetch implements the vcs.System interface.
func (s *System) CanFetch(_ string) bool {
	return true
}

// FetchPath implements the vcs.System interface.
func (s *System) FetchPath(ctx context.Context, path string) (semver.Tags, error) {
	if ctx == nil || s.storage == nil {
		return nil, errors.ErrSystem
	}
	if path == "" {
		return nil, errors.ErrRepository
	}
	var c = make(chan *reference, 1)
	go func() {
		c <- s.fetchWithRetry(path)
	}()
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

// FetchURL implements the vcs.System interface.
func (s *System) FetchURL(ctx context.Context, url string) (semver.Tags, error) {
	if ctx == nil || s.storage == nil {
		return nil, errors.ErrSystem
	}
	if url == "" {
		return nil, errors.ErrRepository
	}
	var c = make(chan *reference, 1)
	go func() {
		c <- s.fetch(url)
	}()
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

// URL builds and returns a URL based on transport data and the given path.
func (t transport) URL(path string) string {
	return t.protocol + path + t.extension
}

func (s *System) fetchWithRetry(path string) (ref *reference) {
	for _, t := range []transport{
		// {protocol: "git://", extension: ".git"},
		// {protocol: "ssh://git@"},
		{protocol: "https://"},
		{protocol: "http://"},
	} {
		ref = s.fetch(t.URL(path))
		if ref.err == nil {
			break
		}
	}
	return
}

func (s *System) fetch(url string) (ref *reference) {
	ref = new(reference)
	rem := git.NewRemote(s.storage, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})
	res, err := rem.List(&git.ListOptions{})
	if err != nil {
		ref.err = vcs.Errorf(Name, errors.ErrFetch, err)
		return
	}
	var n plumbing.ReferenceName
	for _, r := range res {
		n = r.Name()
		if n.IsTag() {
			ref.list = append(ref.list, semver.New(n.Short()))
		}
	}
	return
}

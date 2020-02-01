// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package git

import (
	"context"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/rvflash/goup"
	"github.com/rvflash/goup/internal/semver"
)

// System
type System struct {
	storage storage.Storer
}

// New
func New() *System {
	return &System{storage: memory.NewStorage()}
}

// CanCanFetch implements the VCS interface.
func (s *System) CanFetch(_ string) bool {
	return true
}

type reference struct {
	list semver.Tags
	err  error
}

// FetchContext implements the VCS interface.
func (s *System) FetchContext(ctx context.Context, path string) (semver.Tags, error) {
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

func (s *System) fetchWithRetry(path string) (ref *reference) {
	for _, protocol := range []string{"https://", "http://"} {
		ref = s.fetch(protocol, path)
		if ref.err == nil {
			break
		}
	}
	return
}

func (s *System) fetch(protocol, path string) (ref *reference) {
	ref = new(reference)
	if s.storage == nil {
		ref.err = goup.Errorf("git", goup.ErrSystem)
		return
	}
	if protocol == "" || path == "" {
		ref.err = goup.Errorf("git", goup.ErrRepository)
		return
	}
	rem := git.NewRemote(s.storage, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{protocol + path},
	})
	res, err := rem.List(&git.ListOptions{})
	if err != nil {
		ref.err = goup.Errorf("git", goup.ErrFetch, err)
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

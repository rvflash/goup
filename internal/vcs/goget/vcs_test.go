// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goget_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/goget"
	mock_vcs "github.com/rvflash/goup/testdata/mock/vcs"
)

const (
	pkgName  = "golang.org/x/mod"
	repoURL  = "https://go.googlesource.com/mod"
	tagValue = "v0.1.2"
)

func TestVCS_CanFetch(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  string
			out bool
		}{
			"Default":               {in: ""},
			"Ignore Bitbucket":      {in: "bitbucket.org/repo/all"},
			"Ignore private Gitlab": {in: "gitlab.example.lan/group/pkg"},
			"Ignore public Gitlab":  {in: "gitlab.com/group/pkg"},
			"Ignore Github":         {in: "github.com/golang/mock"},
			"Incomplete":            {in: "golang.org", out: true},
			"Ok":                    {in: pkgName, out: true},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			s := goget.New(mock_vcs.NewMockClientChooser(ctrl), mock_vcs.NewMockSystem(ctrl))
			are.Equal(s.CanFetch(tt.in), tt.out) // mismatch fetch
		})
	}
}

func TestVCS_FetchPath(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			cli  vcs.ClientChooser
			git  vcs.System
			ctx  context.Context
			path string
			res  semver.Tags
			err  error
		}{
			"default":     {err: errors.ErrRepository},
			"http only":   {cli: mock_vcs.NewMockClientChooser(ctrl), err: errors.ErrRepository},
			"system only": {git: mock_vcs.NewMockSystem(ctrl), err: errors.ErrRepository},
			"missing context": {
				cli: mock_vcs.NewMockClientChooser(ctrl),
				git: mock_vcs.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
			"missing path": {
				cli: mock_vcs.NewMockClientChooser(ctrl),
				git: mock_vcs.NewMockSystem(ctrl),
				ctx: context.Background(),
				err: errors.ErrRepository,
			},
			"ok": {
				cli:  newMockClientChooser(ctrl, nil),
				git:  newMockSystem(ctrl),
				ctx:  context.Background(),
				path: pkgName,
				res:  semver.Tags{semver.New(tagValue)},
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			s := goget.New(tt.cli, tt.git)
			res, err := s.FetchPath(tt.ctx, tt.path)
			are.Equal(err, tt.err) // mismatch error
			are.Equal(res, tt.res) // mismatch result
		})
	}
}

func TestVCS_FetchURL(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			cli vcs.ClientChooser
			git vcs.System
			ctx context.Context
			uri string
			res semver.Tags
			err error
		}{
			"default":     {err: errors.ErrSystem},
			"client only": {cli: mock_vcs.NewMockClientChooser(ctrl), err: errors.ErrSystem},
			"system only": {git: mock_vcs.NewMockSystem(ctrl), err: errors.ErrSystem},
			"missing context": {
				cli: mock_vcs.NewMockClientChooser(ctrl),
				git: mock_vcs.NewMockSystem(ctrl),
				err: errors.ErrSystem,
			},
			"missing url": {
				cli: mock_vcs.NewMockClientChooser(ctrl),
				git: mock_vcs.NewMockSystem(ctrl),
				ctx: context.Background(),
				err: errors.ErrRepository,
			},
			"missing system": {
				cli: newMockClientChooser(ctrl, nil),
				ctx: context.Background(),
				uri: repoURL,
				err: errors.ErrSystem,
			},
			"invalid calls": {
				cli: newMockClientChooser(ctrl, errors.ErrFetch),
				git: mock_vcs.NewMockSystem(ctrl),
				ctx: context.Background(),
				uri: repoURL,
				err: errors.ErrFetch,
			},
			"ok": {
				cli: newMockClientChooser(ctrl, nil),
				git: newMockSystem(ctrl),
				ctx: context.Background(),
				uri: repoURL,
				res: semver.Tags{semver.New(tagValue)},
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			s := goget.New(tt.cli, tt.git)
			res, err := s.FetchURL(tt.ctx, tt.uri)
			are.Equal(err, tt.err) // mismatch error
			are.Equal(res, tt.res) // mismatch result
		})
	}
}

const oneTime = 1

func newMockClientChooser(ctrl *gomock.Controller, err error) *mock_vcs.MockClientChooser {
	var (
		c = mock_vcs.NewMockClientChooser(ctrl)
		m = &mockClient{err: err}
	)
	if err != nil {
		c.EXPECT().ClientFor(gomock.Any()).Return(m).AnyTimes()
	} else {
		c.EXPECT().ClientFor(gomock.Any()).Return(m).Times(oneTime)
	}
	return c
}

type mockClient struct {
	err error
}

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	if c.err != nil {
		return nil, c.err
	}
	name := []string{"..", "..", "..", "testdata", "golden", "goget", "default.html"}
	b, err := os.Open(filepath.Join(name...))
	if err != nil {
		return nil, err
	}
	resp := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       b,
		Request:    req,
	}
	return resp, nil
}

func newMockSystem(ctrl *gomock.Controller) *mock_vcs.MockSystem {
	var (
		c = mock_vcs.NewMockSystem(ctrl)
		v = semver.New(tagValue)
	)
	c.EXPECT().FetchURL(gomock.Any(), repoURL).Return(semver.Tags{v}, nil).Times(oneTime)
	return c
}

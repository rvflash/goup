// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package git_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/git"
	mock_vcs "github.com/rvflash/goup/testdata/mock/vcs"
)

const (
	pkgName = "github.com/src-d/go-git"
	repoURL = "https://github.com/src-d/go-git"
)

func TestVCS_CanFetch(t *testing.T) {
	var s git.VCS
	is.New(t).True(s.CanFetch("")) // always true.
}

func TestVCS_FetchPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			cli vcs.ClientChooser
			ctx context.Context
			in  string
			err error
		}{
			"default":         {err: errup.ErrSystem},
			"missing context": {cli: mock_vcs.NewMockClientChooser(ctrl), in: pkgName, err: errup.ErrSystem},
			"missing path":    {cli: mock_vcs.NewMockClientChooser(ctrl), ctx: ctx, err: errup.ErrRepository},
			"ok":              {cli: newMockClientChooser(ctrl), ctx: ctx, in: pkgName},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			s := git.New(tt.cli)
			_, err := s.FetchPath(tt.ctx, tt.in)
			are.True(errors.Is(err, tt.err)) // mismatch error
		})
	}
}

func TestVCS_FetchURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			cli vcs.ClientChooser
			ctx context.Context
			in  string
			err error
		}{
			"default":         {err: errup.ErrSystem},
			"missing context": {cli: mock_vcs.NewMockClientChooser(ctrl), in: repoURL, err: errup.ErrSystem},
			"missing url":     {cli: mock_vcs.NewMockClientChooser(ctrl), ctx: ctx, err: errup.ErrRepository},
			"ok":              {cli: newMockClientChooser(ctrl), ctx: ctx, in: repoURL},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			s := git.New(tt.cli)
			_, err := s.FetchURL(tt.ctx, tt.in)
			are.True(errors.Is(err, tt.err)) // mismatch error
		})
	}
}

func newMockClientChooser(ctrl *gomock.Controller) *mock_vcs.MockClientChooser {
	c := mock_vcs.NewMockClientChooser(ctrl)
	c.EXPECT().ClientFor(pkgName).Return(&http.Client{}).AnyTimes()
	c.EXPECT().AllowInsecure(pkgName).Return(false).AnyTimes()
	return c
}

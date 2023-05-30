// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package git_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/git"
	mockvcs "github.com/rvflash/goup/testdata/mock/vcs"
)

const (
	hostname  = "github.com"
	pkgName   = hostname + "/src-d/go-git"
	repoURL   = "https://" + pkgName
	unsafeURL = "http://" + pkgName
)

func TestVCS_CanFetch(t *testing.T) {
	t.Parallel()
	var s git.VCS
	is.New(t).True(s.CanFetch("")) // always true.
}

func TestVCS_FetchPath(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			cli  vcs.ClientChooser
			auth vcs.BasicAuthentifier
			ctx  context.Context
			in   string
			err  error
		}{
			"Default":              {err: errup.ErrSystem},
			"Missing authentifier": {cli: mockvcs.NewMockClientChooser(ctrl), err: errup.ErrSystem},
			"Missing context": {
				cli:  mockvcs.NewMockClientChooser(ctrl),
				auth: mockvcs.NewMockBasicAuthentifier(ctrl),
				in:   pkgName,
				err:  errup.ErrSystem,
			},
			"Missing path": {
				cli:  mockvcs.NewMockClientChooser(ctrl),
				auth: mockvcs.NewMockBasicAuthentifier(ctrl),
				ctx:  ctx,
				err:  errup.ErrRepository,
			},
			"OK": {
				cli:  newMockClientChooser(ctrl),
				auth: newMockBasicAuthentifier(ctrl),
				ctx:  ctx,
				in:   pkgName,
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			s := git.New(tt.cli, tt.auth)
			_, err := s.FetchPath(tt.ctx, tt.in)
			are.True(errors.Is(err, tt.err)) // mismatch error
		})
	}
}

func TestVCS_FetchURL(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			cli  vcs.ClientChooser
			auth vcs.BasicAuthentifier
			ctx  context.Context
			in   string
			err  error
		}{
			"Default":              {err: errup.ErrSystem},
			"Missing authentifier": {cli: mockvcs.NewMockClientChooser(ctrl), err: errup.ErrSystem},
			"Missing context": {
				cli:  mockvcs.NewMockClientChooser(ctrl),
				auth: mockvcs.NewMockBasicAuthentifier(ctrl),
				in:   repoURL,
				err:  errup.ErrSystem,
			},
			"Missing url": {
				cli:  mockvcs.NewMockClientChooser(ctrl),
				auth: mockvcs.NewMockBasicAuthentifier(ctrl),
				ctx:  ctx,
				err:  errup.ErrRepository,
			},
			"Invalid": {
				cli:  newMockClientChooser(ctrl),
				auth: mockvcs.NewMockBasicAuthentifier(ctrl),
				ctx:  ctx,
				in:   unsafeURL,
				err:  errup.ErrRepository,
			},
			"OK": {
				cli:  mockvcs.NewMockClientChooser(ctrl),
				auth: newMockBasicAuthentifier(ctrl),
				ctx:  ctx,
				in:   repoURL,
			},
		}
	)
	for name, ts := range dt {
		tt := ts
		t.Run(name, func(t *testing.T) {
			s := git.New(tt.cli, tt.auth)
			_, err := s.FetchURL(tt.ctx, tt.in)
			are.True(errors.Is(err, tt.err)) // mismatch error
		})
	}
}

func newMockClientChooser(ctrl *gomock.Controller) *mockvcs.MockClientChooser {
	c := mockvcs.NewMockClientChooser(ctrl)
	c.EXPECT().AllowInsecure(pkgName).Return(false).AnyTimes()
	return c
}

func newMockBasicAuthentifier(ctrl *gomock.Controller) *mockvcs.MockBasicAuthentifier {
	m := mockvcs.NewMockBasicAuthentifier(ctrl)
	m.EXPECT().BasicAuth(hostname).Return(nil).AnyTimes()
	return m
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goget_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/semver"
	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/goget"
	"github.com/rvflash/goup/internal/vcs/mock"
)

func TestVCS_CanFetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  string
			out bool
		}{
			"default":    {in: ""},
			"invalid":    {in: "org"},
			"incomplete": {in: "golang"},
			"wrong":      {in: "github.com/golang/mock"},
			"complete":   {in: "golang.org", out: true},
			"package":    {in: "golang.org/x/mod", out: true},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			vcs := goget.New(mock.NewMockHTTPClient(ctrl), mock.NewMockSystem(ctrl))
			out := vcs.CanFetch(tt.in)
			are.Equal(out, tt.out) // mismatch fetch
		})
	}
}

func TestVCS_FetchPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		are = is.New(t)
		dt  = map[string]struct {
			cli  vcs.HTTPClient
			git  vcs.System
			ctx  context.Context
			path string
			res  []semver.Tags
			err  error
		}{
			"default":         {err: errors.ErrSystem},
			"client only":     {cli: mock.NewMockHTTPClient(ctrl), err: errors.ErrSystem},
			"system only":     {git: mock.NewMockSystem(ctrl), err: errors.ErrSystem},
			"missing context": {cli: mock.NewMockHTTPClient(ctrl), git: mock.NewMockSystem(ctrl), err: errors.ErrSystem},
			"missing path": {
				cli: mock.NewMockHTTPClient(ctrl),
				ctx: context.Background(),
				git: mock.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
			"both protocol failed": {
				cli: mock.NewMockHTTPClient(ctrl),
				ctx: context.Background(),
				git: mock.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
			"https failed": {
				cli: mock.NewMockHTTPClient(ctrl),
				ctx: context.Background(),
				git: mock.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
			"context done": {
				cli: mock.NewMockHTTPClient(ctrl),
				ctx: context.Background(),
				git: mock.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
			"bad context #1": {
				cli: mock.NewMockHTTPClient(ctrl),
				ctx: context.Background(),
				git: mock.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
			"unknown vcs": {
				cli: mock.NewMockHTTPClient(ctrl),
				ctx: context.Background(),
				git: mock.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
			"ok": {
				cli: mock.NewMockHTTPClient(ctrl),
				ctx: context.Background(),
				git: mock.NewMockSystem(ctrl),
				err: errors.ErrRepository,
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			vcs := goget.New(tt.cli, tt.git)
			res, err := vcs.FetchPath(tt.ctx, tt.path)
			are.Equal(err, tt.err) // mismatch error
			are.Equal(res, tt.res) // mismatch result
		})
	}
}

func TestVCS_FetchURL(t *testing.T) {

}

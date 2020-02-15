// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package git

import (
	"context"
	"errors"
	"testing"

	"github.com/matryer/is"
)

const (
	ext  = ".git"
	repo = "github.com/rvflash/goup"
)

func TestTags(t *testing.T) {
	are := is.New(t)
	ref := make(chan *reference, oneRef)
	ctx := context.Background()
	ctxCancel, cancel := context.WithCancel(ctx)
	cancel()

	// Job cancelled.
	_, err := tags(ctxCancel, ref)
	are.Equal(err, context.Canceled)

	// Job failed.
	oops := errors.New("oops")
	ref <- &reference{err: oops}
	_, err = tags(ctx, ref)
	are.Equal(err, oops)

	// Job done.
	ref <- &reference{}
	_, err = tags(ctx, ref)
	are.NoErr(err)
}

func TestTransport_RawURL(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  transport
			out string
		}{
			"default": {out: repo},
			"ssh":     {in: transport{protocol: "ssh://git@"}, out: "ssh://git@github.com/rvflash/goup"},
			"https":   {in: transport{protocol: "https://"}, out: "https://github.com/rvflash/goup"},
			"http":    {in: transport{protocol: "http://", extension: ext}, out: "http://github.com/rvflash/goup.git"},
			//?"git":     {in: transport{protocol: "git://", extension: ext}, out: "git@github.com:rvflash/goup.git"},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			are.Equal(tt.in.rawURL(repo), tt.out) // mismatch result
		})
	}
}

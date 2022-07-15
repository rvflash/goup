// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/pkg/goup"
	"github.com/rvflash/goup/pkg/mod"
)

func TestCheck(t *testing.T) {
	t.Parallel()
	var (
		are  = is.New(t)
		file mod.Mod
	)
	ch := goup.Check(context.Background(), file, goup.Config{})
	msg := <-ch
	are.Equal(msg.Level(), goup.ErrorLevel) // mismatch error
}

func TestEntry_OutDated(t *testing.T) {
	t.Parallel()
	var (
		are = is.New(t)
		msg = &goup.Entry{}
	)
	are.Equal(goup.ErrorLevel, msg.Level()) // mismatch default level
	are.Equal("", msg.Format())             // mismatch default format
	are.Equal(nil, msg.Args())              // mismatch default args
	v, ok := msg.OutDated()
	are.True(!ok)    // not outdated
	are.Equal("", v) // no new version expected
}

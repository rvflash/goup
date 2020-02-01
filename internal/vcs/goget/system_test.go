// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goget_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rvflash/goup/internal/vcs"
	"github.com/rvflash/goup/internal/vcs/git"
	"github.com/rvflash/goup/internal/vcs/goget"
)

func TestNew(t *testing.T) {
	s := goget.New(vcs.NewHTTPClient(3*time.Second), git.New())
	rs, err := s.FetchPath(context.Background(), "golang.org/x/tools")
	fmt.Println(err)
	fmt.Println(rs[5].String())

}

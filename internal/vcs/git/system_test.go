// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package git_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rvflash/gomod/internal/vcs/git"
)

func TestNew(t *testing.T) {
	r := git.New()
	rs, err := r.FetchContext(context.Background(), "github.com/Zenika/MARCEL")
	fmt.Println(err)
	fmt.Println(rs)
}

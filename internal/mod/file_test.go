// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mod_test

import (
	"fmt"
	"testing"

	"github.com/rvflash/goup/internal/mod"
)

func TestOpen(t *testing.T) {
	f, err := mod.OpenFile("/mnt/data/go/src/github.com/rvflash/goup/go.mod")
	fmt.Println(f.Module())
	fmt.Println(f.Dependencies()[0].Path())
	fmt.Println(f.Dependencies()[0].Version())
	fmt.Println(f.Dependencies()[0].Indirect())
	fmt.Println(err)
}

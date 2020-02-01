// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package vcs

import (
	"context"

	"github.com/rvflash/goup/internal/semver"
)

// Remote
type Remote interface {
	CanFetch(path string) bool
	FetchContext(ctx context.Context, path string) (semver.Tags, error)
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package path provides methods to check if a path matches some glob patterns.
package path

import (
	"path"
	"strings"
)

const (
	comma = ","
	slash = "/"
)

// Match reports whether any path prefix of target matches one of the glob patterns.
// globs is a comma-separated list of glob patterns.
func Match(globs, target string) (matched bool) {
	if target == "" {
		return
	}
	var (
		src []string
		dst = strings.Split(target, slash)
	)
	for _, glob := range strings.Split(globs, comma) {
		if glob = strings.TrimSpace(glob); glob == "" {
			continue
		}
		src = strings.Split(glob, slash)
		if len(src) > len(dst) {
			continue
		}
		matched, _ = path.Match(glob, path.Join(dst[:len(src)]...))
		if matched {
			return
		}
	}
	return
}

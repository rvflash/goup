// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package semver

import (
	"fmt"
	"sort"

	"golang.org/x/mod/semver"
)

// Latest
func Latest(versions Tags) Tag {
	switch len(versions) {
	case 0:
		return nil
	case 1:
		return versions[0]
	default:
		sort.Sort(versions)
		return versions[len(versions)-1]
	}
}

// LatestMinor
func LatestMinor(major string, versions Tags) Tag {
	if major == "" || len(versions) == 0 {
		return nil
	}
	sort.Sort(versions)

	var latest Tag
	for _, v := range versions {
		if v.Major() == major {
			latest = v
		}
	}
	return latest
}

// LatestPatch
func LatestPatch(majorMinor string, versions Tags) Tag {
	if majorMinor == "" || len(versions) == 0 {
		return nil
	}
	sort.Sort(versions)

	var latest Tag
	for _, v := range versions {
		if v.MajorMinor() == majorMinor {
			latest = v
		}
	}
	return latest
}

// Compare
func Compare(v, w fmt.Stringer) int {
	if v == nil || w == nil {
		return 0
	}
	return semver.Compare(v.String(), w.String())
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package semver implements comparison of semantic version strings.
package semver

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

// Tags represents a list of versions that can be sorted.
type Tags []Tag

// Not removes from the list of tags the given tag.
func (t Tags) Not(w fmt.Stringer) Tags {
	key := func(t Tags) int {
		for k, v := range t {
			if Compare(v, w) == 0 {
				return k
			}
		}
		return -1
	}
	i := key(t)
	if i < 0 {
		return t
	}
	return append(t[:i], t[i+1:]...)
}

// Len implements the sort interface.New
func (t Tags) Len() int {
	return len(t)
}

// Less implements the sort interface.
func (t Tags) Less(i, j int) bool {
	return Compare(t[i], t[j]) < 0
}

// Swap implements the sort interface.
func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Tag represents a version.
type Tag interface {
	Build() string
	Canonical() string
	IsTag() bool
	IsValid() bool
	Major() string
	MajorMinor() string
	Prerelease() string
	fmt.Stringer
}

// New returns a new instance of PrintVersion.
func New(version string) *Version {
	var (
		v = new(Version)
		s = trim(version)
	)
	if !semver.IsValid(s) {
		return v
	}
	v.raw = version
	v.canonical = semver.Canonical(s)
	v.major = semver.Major(s)
	v.majorMinor = semver.MajorMinor(s)
	v.prerelease = semver.Prerelease(s)
	v.build = semver.Build(s)
	return v
}

func trim(s string) string {
	p := strings.LastIndex(s, "/")
	if p > -1 {
		return s[p+1:]
	}
	return s
}

// Version represents a semantic version.
type Version struct {
	raw        string
	canonical  string
	major      string
	majorMinor string
	prerelease string
	build      string
}

// Build implements the Tag interface.
func (v Version) Build() string {
	return v.build
}

// Canonical implements the Tag interface.
func (v Version) Canonical() string {
	return v.canonical
}

// IsValid implements the Tag interface.
func (v Version) IsValid() bool {
	return v.canonical != ""
}

// IsTag implements the Tag interface.
func (v Version) IsTag() bool {
	return v.canonical != "" && v.build == "" && v.prerelease == ""
}

// Major implements the Tag interface.
func (v Version) Major() string {
	return v.major
}

// MajorMinor implements the Tag interface.
func (v Version) MajorMinor() string {
	return v.majorMinor
}

// Prerelease implements the Tag interface.
func (v Version) Prerelease() string {
	return v.prerelease
}

// String implements the Tag interface.
func (v Version) String() string {
	return v.raw
}

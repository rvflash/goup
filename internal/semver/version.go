// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package semver

import (
	"fmt"

	"golang.org/x/mod/semver"
)

// IsTag
func IsTag(version string) bool {
	return New(version).IsTag()
}

// Tags
type Tags []Tag

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

// Tag
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

// New
func New(version string) *Version {
	v := new(Version)
	if !semver.IsValid(version) {
		return v
	}
	v.raw = version
	v.canonical = semver.Canonical(version)
	v.major = semver.Major(version)
	v.majorMinor = semver.MajorMinor(version)
	v.prerelease = semver.Prerelease(version)
	v.build = semver.Build(version)
	return v
}

// Version
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
	return v.build == "" && v.prerelease == ""
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

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matryer/is"

	"github.com/rvflash/goup/internal/mod"
)

const (
	release  = "v1.0.2"
	repoName = "example.com/group/go"
)

func TestChecked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		m   mod.Module
		are = is.New(t)
	)
	are.Equal(checked(m), "") // mismatch default
	m = newMockModule(ctrl, false)
	are.Equal(checked(m), "example.com/group/go v1.0.2 is up to date") // mismatch result
}

func TestSkipped(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		m   mod.Module
		are = is.New(t)
	)
	are.Equal(skipped(m), "") // mismatch default
	m = newMockModule(ctrl, false)
	are.Equal(skipped(m), "example.com/group/go v1.0.2 update skipped") // mismatch result
}

func TestUpdated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		m   mod.Module
		are = is.New(t)
	)
	are.Equal(updated(m), "") // mismatch default
	m = newMockModule(ctrl, false)
	are.Equal(updated(m), "example.com/group/go v1.0.2 was updated") // mismatch result
}

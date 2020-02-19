// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	uperrs "github.com/rvflash/goup/internal/errors"
)

const charset = "utf8"

func TestNewCharset(t *testing.T) {
	is.New(t).Equal(uperrs.NewCharset(charset).Error(), "unsupported charset: "+charset)
}

func TestNewSecurityIssue(t *testing.T) {
	is.New(t).Equal(
		uperrs.NewSecurityIssue("http://example.com").Error(),
		"unsecured call to http://example.com cancelled: failed to list tags",
	)
}

func TestNewMissingData(t *testing.T) {
	var (
		err = uperrs.NewMissingData(charset)
		are = is.New(t)
	)
	are.True(errors.Is(err, uperrs.ErrMissing))      // wrong error kind
	are.True(strings.Contains(err.Error(), charset)) // missing source
}

func TestErrUp_Error(t *testing.T) {
	is.New(t).Equal(uperrs.ErrRepository.Error(), "invalid repository")
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	errup "github.com/rvflash/goup/internal/errors"
)

const charset = "utf8"

func TestFailure_Error(t *testing.T) {
	var (
		are = is.New(t)
		msg = "oops"
		err = errors.New(msg)
	)
	t.Run("Default", func(t *testing.T) {
		var e errup.Failure
		are.Equal("", e.Error()) // mismatch error message
		are.NoErr(e.Unwrap())    // mismatch error
	})
	t.Run("Ok", func(t *testing.T) {
		e := &errup.Failure{
			Mod: charset,
			Err: err,
		}
		are.True(errors.Is(e, err) && strings.Contains(e.Error(), charset)) // mismatch error message
		are.Equal(err, e.Unwrap())                                          // mismatch error
	})
}

func TestOutOfDate_Error(t *testing.T) {
	const (
		v1 = "v0.0.1"
		v2 = "v0.0.2"
	)
	err := &errup.OutOfDate{
		Mod:        charset,
		OldVersion: v1,
		NewVersion: v2,
	}
	is.New(t).Equal(charset+" "+v1+" must be updated with "+v2, err.Error()) // mismatch error message
}

func TestNewCharset(t *testing.T) {
	is.New(t).Equal(errup.NewCharset(charset).Error(), "unsupported charset: "+charset)
}

func TestNewSecurityIssue(t *testing.T) {
	is.New(t).Equal(
		errup.NewSecurityIssue("http://example.com").Error(),
		"unsecured call to http://example.com cancelled: failed to list tags",
	)
}

func TestNewMissingData(t *testing.T) {
	var (
		err = errup.NewMissingData(charset)
		are = is.New(t)
	)
	are.True(errors.Is(err, errup.ErrMissing))       // wrong error kind
	are.True(strings.Contains(err.Error(), charset)) // missing source
}

func TestErrUp_Error(t *testing.T) {
	is.New(t).Equal(errup.ErrRepository.Error(), "invalid repository")
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package errors_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/errors"
)

func TestErrUp_Error(t *testing.T) {
	is.New(t).Equal(errors.ErrRepository.Error(), "invalid repository")
}

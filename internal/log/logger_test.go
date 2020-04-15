// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package log_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/log"
)

const (
	t0 = "foo"
	prefixedT0Len = 10
)

func TestLogger_SetVerbose(t *testing.T) {
	f, cleanup := newFile(t)
	t.Cleanup(func() {
		err := cleanup()
		if err != nil {
			t.Error(err)
		}
	})
	// Default
	var (
		are = is.New(t)
		l = log.New(f, false)
	)
	l.Debugf(t0)
	_, size := readFile(t, f)
	are.Equal(size, 0) // expected no content
	// Enabled
	l.SetVerbose(true)
	l.Debugf(t0)
	_, size = readFile(t, f)
	are.Equal(size, prefixedT0Len) // expected content
	// Disabled
	l.SetVerbose(false)
	l.Debugf(t0)
	_, size = readFile(t, f)
	are.Equal(size, prefixedT0Len) // expected no new content
}

func readFile(t *testing.T, f *os.File) ([]byte, int) {
	buf, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return buf, len(buf)
}

func newFile(t *testing.T) (*os.File, func() error) {
	f, err := ioutil.TempFile("", "goup")
	if err != nil {
		t.Fatal(err)
	}
	return f, func() error {
		_ = f.Close()
		return os.Remove(f.Name())
	}
}

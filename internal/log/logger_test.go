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
	d0            = "foo"
	prefixedD0Len = 10

	fDebugName = "Debugf"
	fErrorName = "Errorf"
	fInfoName  = "Infof"
	fWarnName  = "Warnf"
)

func TestNew(t *testing.T) {
	// /dev/null
	l := log.New(nil, false)
	callf(t, l, fDebugName, d0)
	callf(t, l, fErrorName, d0)
	callf(t, l, fInfoName, d0)
	callf(t, l, fWarnName, d0)
}

func TestDevNull(t *testing.T) {
	l := log.DevNull()
	callf(t, l, fDebugName, d0)
	callf(t, l, fErrorName, d0)
	callf(t, l, fInfoName, d0)
	callf(t, l, fWarnName, d0)
}
func TestLogger_Debugf(t *testing.T) {
	const (
		pattern  = "%s, %s (%s)"
		blank    = "goup: \n"
		noArgs   = "goup: foo\n"
		withArgs = "goup: foo, 9 (90.2)\n"
	)
	var (
		args = []interface{}{d0, "9", 90.2}
		are  = is.New(t)
		dt   = map[string]struct {
			method  string
			verbose bool
			in      string
			args    []interface{}
			out     string
		}{
			"Default Debugf":              {method: fDebugName},
			"Default Errorf":              {method: fErrorName, out: blank},
			"Default Infof":               {method: fInfoName, out: blank},
			"Default Warnf":               {method: fWarnName, out: blank},
			"Debugf without args":         {method: fDebugName, in: d0},
			"Errorf without args":         {method: fErrorName, in: d0, out: noArgs},
			"Infof without args":          {method: fInfoName, in: d0, out: noArgs},
			"Warnf without args":          {method: fWarnName, in: d0, out: noArgs},
			"Debugf with args":            {method: fDebugName, in: d0, args: args},
			"Errorf with args":            {method: fErrorName, in: pattern, args: args, out: withArgs},
			"Infof with args":             {method: fInfoName, in: pattern, args: args, out: withArgs},
			"Warnf with args":             {method: fWarnName, in: pattern, args: args, out: withArgs},
			"Verbose Debugf without args": {method: fDebugName, verbose: true, in: d0, out: noArgs},
			"Verbose Debugf with args":    {method: fDebugName, verbose: true, in: pattern, args: args, out: withArgs},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, cleanup := newFile(t)
			defer func() {
				err := cleanup()
				if err != nil {
					t.Error(err)
				}
			}()
			l := log.New(out, false)
			l.SetVerbose(tt.verbose)
			callf(t, l, tt.method, tt.in, tt.args...)
			s, _ := readFile(t, out)
			are.Equal(tt.out, s) // mismatch result
		})
	}
}

func TestLogger_SetVerbose(t *testing.T) {
	f, cleanup := newFile(t)
	defer func() {
		err := cleanup()
		if err != nil {
			t.Error(err)
		}
	}()
	// Default
	var (
		are = is.New(t)
		l   = log.New(f, false)
	)
	l.Debugf(d0)
	_, size := readFile(t, f)
	are.Equal(size, 0) // expected no content
	// Enabled
	l.SetVerbose(true)
	l.Debugf(d0)
	_, size = readFile(t, f)
	are.Equal(size, prefixedD0Len) // expected content
	// Disabled
	l.SetVerbose(false)
	l.Debugf(d0)
	_, size = readFile(t, f)
	are.Equal(size, prefixedD0Len) // expected no new content
}

func callf(t *testing.T, w log.Printer, method string, format string, args ...interface{}) {
	switch method {
	case fDebugName:
		w.Debugf(format, args...)
	case fErrorName:
		w.Errorf(format, args...)
	case fInfoName:
		w.Infof(format, args...)
	case fWarnName:
		w.Warnf(format, args...)
	default:
		t.Fatal("unknown method named:" + method)
	}
}

func readFile(t *testing.T, f *os.File) (string, int) {
	buf, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Error(err)
		return "", 0
	}
	return string(buf), len(buf)
}

func newFile(t *testing.T) (*os.File, func() error) {
	f, err := ioutil.TempFile("", log.Prefix)
	if err != nil {
		t.Fatal(err)
	}
	return f, func() error {
		_ = f.Close()
		return os.Remove(f.Name())
	}
}

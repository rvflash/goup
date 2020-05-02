// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package main provides the Go Up command line tool.
package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rvflash/goup/internal/app"
	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/log"
	"github.com/rvflash/goup/internal/signal"
	"github.com/rvflash/goup/pkg/goup"
)

// Filled by the CI when building.
var buildVersion string

const (
	comma      = ","
	errorCode  = 1
	goInsecure = "GOINSECURE"
	goPrivate  = "GOPRIVATE"
	timeout    = time.Minute
)

func main() {
	var (
		c = goup.Config{
			InsecurePatterns: patterns(os.Getenv(goInsecure), os.Getenv(goPrivate)),
		}
		s = "exclude indirect modules"
		l = log.New(os.Stderr, isatty.IsTerminal(os.Stderr.Fd()))
	)
	flag.BoolVar(&c.ExcludeIndirect, "i", false, s)
	s = "exit on first error occurred"
	flag.BoolVar(&c.Strict, "s", false, s)
	s = "ensure to have the latest major version"
	flag.BoolVar(&c.Major, "M", false, s)
	s = "ensure to have the latest couple major with minor version"
	flag.BoolVar(&c.MajorMinor, "m", false, s)
	s = "comma-separated list of glob patterns to match the repository paths where to force tag usage."
	flag.StringVar(&c.OnlyReleases, "r", "", s)
	s = "maximum time duration"
	flag.DurationVar(&c.Timeout, "t", timeout, s)
	s = "force the update of the go.mod file as advised"
	flag.BoolVar(&c.ForceUpdate, "f", false, s)
	s = "print version"
	flag.BoolVar(&c.PrintVersion, "V", false, s)
	s = "verbose output"
	flag.BoolVar(&c.Verbose, "v", false, s)
	flag.Parse()
	l.SetVerbose(c.Verbose)

	err := run(signal.Background(), c, flag.Args(), l)
	if err != nil {
		if err != errors.ErrMod {
			l.Errorf(err.Error())
		}
		os.Exit(errorCode)
	}
}

func patterns(v ...string) string {
	var a []string
	for _, s := range v {
		if s = strings.TrimSpace(s); s != "" {
			a = append(a, s)
		}
	}
	return strings.Join(a, comma)
}

func run(ctx context.Context, cnf goup.Config, args []string, out log.Printer) error {
	a, err := app.Open(
		buildVersion,
		app.WithLogger(out),
	)
	if err != nil {
		return err
	}
	a.Config = cnf

	if a.Check(ctx, args) {
		return errors.ErrMod
	}
	return nil
}

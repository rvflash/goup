// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/rvflash/goup/internal/app"
	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/signal"
	"github.com/rvflash/goup/pkg/goup"
)

// Filled by the CI when building.
var buildVersion string

const (
	errorCode = 1
	timeout   = 10 * time.Second
)

func main() {
	var (
		c goup.Config
		s = "exclude indirect modules"
	)
	flag.BoolVar(&c.ExcludeIndirect, "i", false, s)
	s = "exit on first error occurred"
	flag.BoolVar(&c.Strict, "s", false, s)
	s = "ensure to have the latest major version"
	flag.BoolVar(&c.Major, "M", false, s)
	s = "ensure to have the latest couple major with minor version"
	flag.BoolVar(&c.MajorMinor, "m", false, s)
	s = "comma separated list of repositories (or part of), where forcing tag usage"
	flag.StringVar(&c.OnlyReleases, "r", "", s)
	s = "maximum time duration"
	flag.DurationVar(&c.Timeout, "t", timeout, s)
	// todo next release
	//s = "force the update of the go.mod file as advised"
	//flag.BoolVar(&c.Force, "f", false, s)
	s = "verbose output"
	flag.BoolVar(&c.Verbose, "v", false, s)
	s = "print version"
	flag.BoolVar(&c.PrintVersion, "V", false, s)
	flag.Parse()

	err := run(signal.Background(), c, flag.Args(), os.Stdout, os.Stderr)
	if err != nil {
		if err != errors.ErrMod {
			log.SetPrefix(app.LogPrefix)
			log.Println(err)
		}
		os.Exit(errorCode)
	}
}

func run(ctx context.Context, cnf goup.Config, args []string, stdout, stderr io.Writer) error {
	a, err := app.Open(
		buildVersion,
		app.WithErrOutput(stderr),
		app.WithOutput(stdout),
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

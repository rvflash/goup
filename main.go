// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/rvflash/goup/internal/app"
	"github.com/rvflash/goup/internal/signal"
)

// Filled by the CI when building.
var buildVersion string

const (
	errorCode = 1
	timeout   = 10 * time.Second
)

func main() {
	a, err := app.Open(buildVersion)
	if err != nil {
		log.SetPrefix(app.LogPrefix)
		log.Fatal(err)
	}
	s := "exclude indirect modules"
	flag.BoolVar(&a.ExcludeIndirect, "i", false, s)
	s = "exit on first error occurred"
	flag.BoolVar(&a.Fast, "f", false, s)
	s = "ensure to have the latest major version"
	flag.BoolVar(&a.Major, "M", false, s)
	s = "ensure to have the latest couple major with minor version"
	flag.BoolVar(&a.MajorMinor, "m", false, s)
	s = "comma separated list of repositories (or part of), where forcing tag usage"
	flag.StringVar(&a.OnlyReleases, "r", "", s)
	s = "maximum time duration"
	flag.DurationVar(&a.Timeout, "t", timeout, s)
	// todo in the next release
	//s = "update the go.mod file as advised"
	//flag.DurationVar(&a.Update, "u", timeout, s)
	s = "verbose output"
	flag.BoolVar(&a.Verbose, "v", false, s)
	s = "print version"
	flag.BoolVar(&a.Version, "V", false, s)
	flag.Parse()

	ctx := signal.Background()
	if !a.Check(ctx, flag.Args()) {
		os.Exit(errorCode)
	}
}

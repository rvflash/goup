// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"log"
	"os"
	"time"
)

// New
func New(version string) *App {
	return &App{
		buildVersion: version,
		stdin:        log.New(os.Stdin, "goup: ", 0),
		stderr:       log.New(os.Stderr, "goup: ", 0),
	}
}

// App
type App struct {
	Fast            bool
	Major           bool
	MajorMinor      bool
	OnlyReleases    string
	ExcludeIndirect bool
	Timeout         time.Duration
	Verbose         bool
	Version         bool

	stdin, stderr *log.Logger
	buildVersion  string
}

// Check
func (a *App) Check(paths []string) bool {
	if a.Version {
		a.errorf("version %s", a.buildVersion)
	}
	return false
}

func (a *App) errorf(format string, v ...interface{}) {
	if a.stderr == nil {
		return
	}
	a.stderr.Printf(format, v...)
}

func (a *App) printf(format string, v ...interface{}) {
	if a.stdin == nil || !a.Verbose {
		return
	}
	a.stdin.Printf(format, v...)
}

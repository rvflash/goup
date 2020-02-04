// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package app manages the GoUp app.
package app

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/pkg/goup"
)

// LogPrefix is the prefix used when logging.
const LogPrefix = "goup: "

// Configurator defines the interface used to set settings.
type Configurator func(a *App) error

// WithChecker defines the go module parser to use.
// By default, it's a new instance of the goup package.
func WithChecker(f goup.Checker) Configurator {
	return func(a *App) error {
		if f == nil {
			return errors.NewMissingData("checker")
		}
		a.check = f
		return nil
	}
}

// WithErrOutput defines the receiver of update orders.
// By default it is the standard error file descriptor.
func WithErrOutput(w io.Writer) Configurator {
	return func(a *App) error {
		if w == nil {
			return errors.NewMissingData("err output")
		}
		a.stderr = log.New(w, LogPrefix, 0)
		return nil
	}
}

// WithParser defines the go module parser to use.
// By default, the Parse method from the internal parse package.
func WithParser(f mod.Parser) Configurator {
	return func(a *App) error {
		if f == nil {
			return errors.NewMissingData("parser")
		}
		a.parse = f
		return nil
	}
}

// WithOutput defines the receiver of update advices.
// By default it is the standard input file descriptor.
func WithOutput(w io.Writer) Configurator {
	return func(a *App) error {
		if w == nil {
			return errors.NewMissingData("output")
		}
		a.stdin = log.New(w, LogPrefix, 0)
		return nil
	}
}

// Open tries to create a new instance of App.
func Open(version string, opts ...Configurator) (*App, error) {
	a := &App{
		buildVersion: version,
	}
	opts = append([]Configurator{
		WithErrOutput(os.Stderr),
		WithOutput(os.Stdin),
		WithParser(mod.Parse),
		WithChecker(goup.GoUp{}),
	}, opts...)
	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// App represents an application.
type App struct {
	goup.Config

	check         goup.Checker
	parse         mod.Parser
	stdin, stderr *log.Logger
	buildVersion  string
}

// Check launches the analyses of given paths.
func (a *App) Check(ctx context.Context, paths []string) (failure bool) {
	if a.PrintVersion {
		a.errorf("version %s\n", a.buildVersion)
		if len(paths) == 0 {
			// Default behavior: without specified path, nothing else to do.
			return
		}
	}
	for _, path := range checkPaths(paths) {
		f, err := a.parse(path)
		if err != nil {
			return true
		}
		for _, tip := range a.check.CheckFile(ctx, f, a.Config).Tips() {
			err := tip.Err()
			if err != nil {
				a.errorf("%s: %s\n", f.Module(), err.Error())
				failure = true
			} else {
				a.printf("%s: %s\n", f.Module(), tip.String())
			}
		}
	}
	return failure
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

const (
	currentDir = "."
	recursive  = "./..."
)

func checkPaths(paths []string) []string {
	switch len(paths) {
	case 0:
		return []string{filePath(currentDir)}
	case 1:
		if paths[0] == recursive {
			return walkPath(currentDir)
		}
	}
	for k, v := range paths {
		paths[k] = filePath(v)
	}
	return paths
}

func filePath(path string) string {
	if filepath.Base(path) == mod.Filename {
		return path
	}
	return filepath.Join(path, mod.Filename)
}

func walkPath(root string) []string {
	var res []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) == mod.Filename {
			res = append(res, path)
		}
		return nil
	})
	if err != nil {
		return []string{root}
	}
	return res
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package app manages the GoUp app.
package app

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rvflash/goup/internal/errors"
	"github.com/rvflash/goup/internal/log"
	"github.com/rvflash/goup/pkg/goup"
	"github.com/rvflash/goup/pkg/mod"
)

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

// WithLogger defines the logger used to print events.
// By default we use a NullLogger.
func WithLogger(l log.Printer) Configurator {
	return func(a *App) error {
		if l == nil {
			return errors.NewMissingData("output")
		}
		a.logger = l
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

// Open tries to create a new instance of App.
func Open(version string, opts ...Configurator) (*App, error) {
	a := &App{
		buildVersion: version,
	}
	opts = append([]Configurator{
		WithLogger(log.NullLogger),
		WithParser(mod.Parse),
		WithChecker(goup.Check),
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

	check        goup.Checker
	parse        mod.Parser
	logger       log.Printer
	buildVersion string
}

// Check launches the analyses of given paths.
func (a *App) Check(ctx context.Context, paths []string) (failure bool) {
	if !a.ready(ctx) {
		a.logger.Errorf(context.Canceled.Error())
		return true
	}
	if a.PrintVersion {
		a.logger.Infof("version %s", a.buildVersion)
		if len(paths) == 0 {
			// Without explicit path: nothing else to do.
			return false
		}
	}
	for _, path := range checkPaths(paths) {
		f, err := a.parse(path)
		if err != nil {
			a.logger.Errorf(err.Error())
			return true
		}
		for msg := range a.check(ctx, f, a.Config) {
			switch msg.Level() {
			case goup.DebugLevel:
				a.logger.Debugf(msg.Format(), msg.Args()...)
			case goup.InfoLevel:
				a.logger.Infof(msg.Format(), msg.Args()...)
			case goup.WarnLevel:
				a.logger.Warnf(msg.Format(), msg.Args()...)
				failure = true
			default:
				a.logger.Errorf(msg.Format(), msg.Args()...)
				failure = true
			}
		}
	}
	return failure
}

func (a *App) ready(ctx context.Context) bool {
	return ctx != nil && a.check != nil && a.parse != nil && a.logger != nil
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

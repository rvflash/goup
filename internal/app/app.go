// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package app

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/rvflash/goup/internal/mod"
	"github.com/rvflash/goup/pkg/goup"
)

const prefix = "goup: "

// New
func New(version string) *App {
	return &App{
		buildVersion: version,
		stdin:        log.New(os.Stdin, prefix, 0),
		stderr:       log.New(os.Stderr, prefix, 0),
	}
}

// App
type App struct {
	goup.Config

	stdin, stderr *log.Logger
	buildVersion  string
}

// Check
func (a *App) Check(ctx context.Context, paths []string) bool {
	if a.Version {
		a.errorf("version %s\n", a.buildVersion)
	}
	var errorExit bool
	for _, path := range checkPaths(paths) {
		f, err := mod.OpenFile(path)
		if err != nil {
			return true
		}
		for _, tip := range goup.File(ctx, f, a.Config).Tips() {
			err := tip.Err()
			if err != nil {
				a.errorf("%s: %s\n", f.Module(), err.Error())
				errorExit = true
			} else {
				a.printf("%s: %s\n", f.Module(), tip.String())
			}
		}
	}
	return errorExit
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

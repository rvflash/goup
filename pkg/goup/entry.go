// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"github.com/rvflash/goup/pkg/mod"
)

// Level defines the log level.
type Level uint8

// List of available levels.
const (
	// ErrorLevel is used for errors that should definitely be noted.
	ErrorLevel Level = iota
	// WarnLevel designates non-critical entries that deserve eyes and can be recovered.
	WarnLevel
	// InfoLevel notifies changes to the user.
	InfoLevel
	// DebugLevel only enabled with verbose mode enabled. Allows to follow decisions.
	DebugLevel
)

// Message exposes entry properties.
type Message interface {
	Level() Level
	Format() string
	Args() []interface{}
	OutDated() (newVersion string, ok bool)
}

func newEntry(level Level, format string, a ...interface{}) *entry {
	return &entry{
		Kind:    level,
		Message: format,
		Data:    a,
	}
}

type entry struct {
	Kind    Level
	Message string
	Data    []interface{}
}

// Args implements the Message interface.
func (e *entry) Args() []interface{} {
	if e == nil {
		return nil
	}
	return e.Data
}

// Level implements the Message interface.
func (e *entry) Level() Level {
	if e == nil {
		return ErrorLevel
	}
	return e.Kind
}

// Format implements the Message interface.
func (e *entry) Format() string {
	if e == nil {
		return ""
	}
	return e.Message
}

const newVersionPos = 2

// OutDated implements the Message interface.
func (e *entry) OutDated() (newVersion string, ok bool) {
	if e == nil || e.Level() != WarnLevel || len(e.Args()) != newVersionPos+1 {
		return
	}
	return e.Args()[newVersionPos].(string), true
}

func newCheck(dep mod.Module) *entry {
	if dep == nil {
		// Avoids panic.
		return nil
	}
	return newEntry(DebugLevel, "%s: %s is up to date", dep.Path(), dep.Version().String())
}

func newError(err error, file mod.Mod) *entry {
	if err == nil || file == nil {
		// Avoids panic.
		return nil
	}
	return newEntry(ErrorLevel, "%s: "+err.Error(), file.Module())
}

func newFailure(err error, dep mod.Module) *entry {
	if dep == nil || err == nil {
		// Avoids panic.
		return nil
	}
	return newEntry(ErrorLevel, "%s: check failed: %s", dep.Path(), err)
}

func newSkip(dep mod.Module) *entry {
	if dep == nil {
		// Avoids panic.
		return nil
	}
	return newEntry(DebugLevel, "%s: %s update skipped: indirect", dep.Path(), dep.Version().String())
}

func newUpdate(dep mod.Module, newVersion string) *entry {
	if dep == nil {
		// Avoids panic.
		return nil
	}
	return newEntry(InfoLevel, "%s: %s will be updated to %s", dep.Path(), dep.Version().String(), newVersion)
}

func newOutOfDate(dep mod.Module, newVersion string) *entry {
	if dep == nil {
		// Avoids panic.
		return nil
	}
	return newEntry(WarnLevel, "%s: %s must be updated to %s", dep.Path(), dep.Version().String(), newVersion)
}

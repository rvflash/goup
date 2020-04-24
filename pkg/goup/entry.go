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

// Message exposes Entry properties.
type Message interface {
	Args() []interface{}
	Format() string
	Level() Level
	OutDated() (newVersion string, ok bool)
}

// NewEntry returns a new Entry.
func NewEntry(level Level, format string, a ...interface{}) *Entry {
	return &Entry{
		Kind:    level,
		Message: format,
		Data:    a,
	}
}

// Entry represents a message.
type Entry struct {
	Kind    Level
	Message string
	Data    []interface{}
}

// Args implements the Message interface.
func (e *Entry) Args() []interface{} {
	if e == nil {
		return nil
	}
	return e.Data
}

// Format implements the Message interface.
func (e *Entry) Format() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Level implements the Message interface.
func (e *Entry) Level() Level {
	if e == nil {
		return ErrorLevel
	}
	return e.Kind
}

// OutDated implements the Message interface.
func (e *Entry) OutDated() (newVersion string, ok bool) {
	const newVersionPos = 2
	if e == nil || e.Level() != WarnLevel || len(e.Args()) != newVersionPos+1 {
		return
	}
	return e.Args()[newVersionPos].(string), true
}

func newCheck(dep mod.Module) *Entry {
	if dep == nil {
		return nil
	}
	return NewEntry(DebugLevel, "%s: %s is up to date", dep.Path(), dep.Version().String())
}

func newError(err error, file mod.Mod) *Entry {
	if err == nil || file == nil {
		return nil
	}
	return NewEntry(ErrorLevel, "%s: "+err.Error(), file.Module())
}

func newFailure(err error, dep mod.Module) *Entry {
	if err == nil || dep == nil {
		return nil
	}
	return NewEntry(ErrorLevel, "%s: check failed: %s", dep.Path(), err)
}

func newSkip(dep mod.Module) *Entry {
	if dep == nil {
		return nil
	}
	return NewEntry(DebugLevel, "%s: %s update skipped: indirect", dep.Path(), dep.Version().String())
}

func newUpdate(dep mod.Module, newVersion string) *Entry {
	if dep == nil {
		return nil
	}
	return NewEntry(InfoLevel, "%s: %s will be updated to %s", dep.Path(), dep.Version().String(), newVersion)
}

func newOutOfDate(dep mod.Module, newVersion string) *Entry {
	if dep == nil {
		return nil
	}
	return NewEntry(WarnLevel, "%s: %s must be updated to %s", dep.Path(), dep.Version().String(), newVersion)
}

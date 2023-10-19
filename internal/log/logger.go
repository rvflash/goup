// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package log exposes the interface to implement as logger and provides some implementations.
package log

import (
	"io"
	"log"

	"github.com/fatih/color"
)

// Prefix is the prefix used when logging.
const Prefix = "goup: "

// Printer must be implemented to act as a logger for Go Up.
type Printer interface {
	// Debugf logs a message at level Debug on the standard logger.
	Debugf(format string, args ...interface{})
	// Errorf logs a message at level Error on the standard logger.
	Errorf(format string, args ...interface{})
	// Infof logs a message at level Info on the standard logger.
	Infof(format string, args ...interface{})
	// Warnf logs a message at level Warn on the standard logger.
	Warnf(format string, args ...interface{})
}

// DevNull is the default logger only used to mock the logger interface and do nothing else.
func DevNull() *Logger { return &Logger{} }

// New returns a new instance of a Logger.
func New(out io.Writer, tty bool) *Logger {
	if out == nil {
		// Avoids panics.
		return DevNull()
	}
	return &Logger{
		stderr: log.New(out, Prefix, 0),
		cyan:   sPrintFunc(tty, color.FgCyan),
		green:  sPrintFunc(tty, color.FgGreen),
		red:    sPrintFunc(tty, color.FgRed),
		yellow: sPrintFunc(tty, color.FgHiYellow),
	}
}

// Logger is the logger.
type Logger struct {
	stderr  *log.Logger
	verbose bool
	cyan,
	green,
	red,
	yellow func(a ...interface{}) string
}

// SetVerbose enabled verbose mode.
func (l *Logger) SetVerbose(ok bool) {
	l.verbose = ok
}

// Debugf implements the Printer interface.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if !l.verbose {
		// Avoid panics
		return
	}
	l.printf(format, l.cyan, args...)
}

// Errorf implements the Printer interface.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.printf(format, l.red, args...)
}

// Infof implements the Printer interface.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.printf(format, l.green, args...)
}

// Warnf implements the Printer interface.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.printf(format, l.yellow, args...)
}

func (l *Logger) printf(format string, color func(a ...interface{}) string, args ...interface{}) {
	if l.stderr == nil || color == nil {
		// /dev/null
		return
	}
	if len(args) == 0 {
		// No argument, the entire message is colored
		l.stderr.Printf(color(format))
		return
	}
	l.stderr.Printf(format, colors(color, args)...)
}

func colors(f func(a ...interface{}) string, args []interface{}) []interface{} {
	res := make([]interface{}, len(args))
	for k, v := range args {
		res[k] = f(v)
	}
	return res
}

func sPrintFunc(tty bool, values ...color.Attribute) func(a ...interface{}) string {
	c := color.New(values...)
	if !tty {
		c.DisableColor()
	}
	return c.SprintFunc()
}

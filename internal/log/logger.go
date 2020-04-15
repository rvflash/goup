// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package log

import (
	"io"
	"log"

	"github.com/logrusorgru/aurora"
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

// NullLogger is the default logger only used to mock the logger interface and do nothing.
var NullLogger = &Logger{}

// New returns a new instance of a Logger.
func New(out io.Writer, tty bool) *Logger {
	return &Logger{
		stderr: log.New(out, Prefix, 0),
		colors: aurora.NewAurora(tty),
	}
}

// Logger
type Logger struct {
	stderr  *log.Logger
	colors  aurora.Aurora
	verbose bool
}

// SetVerbose enabled verbose mode.
func (l *Logger) SetVerbose(ok bool) {
	l.verbose = ok
}

// Debugf implements the Printer interface.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l == nil || !l.verbose {
		// Avoid panics
		return
	}
	l.printf(format, l.colors.Cyan, args...)
}

// Errorf implements the Printer interface.
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l == nil {
		// Avoid panics
		return
	}
	l.printf(format, l.colors.Red, args...)
}

// Infof implements the Printer interface.
func (l *Logger) Infof(format string, args ...interface{}) {
	if l == nil {
		// Avoid panics
		return
	}
	l.printf(format, l.colors.Green, args...)
}

// Warnf implements the Printer interface.
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l == nil {
		// Avoid panics
		return
	}
	l.printf(format, l.colors.Yellow, args...)
}

func (l *Logger) printf(format string, color func(arg interface{}) aurora.Value, args ...interface{}) {
	if color == nil {
		// No color
		l.stderr.Printf(format, args...)
		return
	}
	if len(args) == 0 {
		// No argument, the entire message is colored
		l.stderr.Printf(color(format).String())
		return
	}
	l.stderr.Printf(format, colors(color, args)...)
}

func colors(f func(arg interface{}) aurora.Value, args []interface{}) []interface{} {
	if f == nil || len(args) == 0 {
		return nil
	}
	for k, v := range args {
		args[k] = f(v)
	}
	return args
}

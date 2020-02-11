// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package signal provides methods to listen OS signal and attaches them to context.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Background listens the OS signals Interrupt (CTRL+C) or SIGTERM.
// It's useful when you want to execute something before stopping the application.
func Background() context.Context {
	ctx := context.Background()
	return listen(ctx, os.Interrupt, syscall.SIGTERM)
}

func listen(parent context.Context, sig ...os.Signal) context.Context {
	// Extends the given context with a new cancellation behavior.
	ctx, cancel := context.WithCancel(parent)
	// Listens the given signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()
	return ctx
}

// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package signal_test

import (
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rvflash/goup/internal/signal"
)

func TestBackground(t *testing.T) {
	var (
		c   int32
		are = is.New(t)
		ctx = signal.Background()
	)
	go func() {
		<-ctx.Done()
		atomic.AddInt32(&c, 1)
	}()

	are.NoErr(syscall.Kill(syscall.Getpid(), syscall.SIGTERM))
	time.Sleep(100 * time.Millisecond)
	are.Equal(atomic.LoadInt32(&c), int32(1))
}

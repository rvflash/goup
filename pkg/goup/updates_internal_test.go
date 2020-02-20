// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package goup

import (
	"errors"
	"sync"
	"testing"

	"github.com/matryer/is"
)

func TestUpdates_Tips(t *testing.T) {
	var (
		rs    updates
		w8    sync.WaitGroup
		are   = is.New(t)
		oops  = errors.New("oops")
		sorry = errors.New("sorry my bad")
		could = map[string]struct{}{
			"Willow":                 {},
			"Star Wars":              {},
			"Harry Potter":           {},
			"Narnia: Prince Caspian": {},
		}
		must = []error{oops, sorry}
	)
	for s := range could {
		w8.Add(delta)
		go func(s string) {
			rs.could(s)
			w8.Done()
		}(s)
	}
	for _, err := range must {
		w8.Add(delta)
		go func(err error) {
			rs.must(err)
			w8.Done()
		}(err)
	}
	w8.Wait()
	are.Equal(len(rs.tips()), len(could)+len(must)) // mismatch tips size

	var (
		n   int
		err error
	)
	for _, t := range rs.tips() {
		err = t.Err()
		if err != nil {
			are.True(errors.Is(err, oops) || errors.Is(err, sorry)) // unexpected error
			are.Equal(t.String(), t.Err().Error())                  // mismatch message
			n++
		} else {
			_, ok := could[t.String()]
			are.True(ok) // unexpected message
		}
	}
	are.Equal(n, len(must)) // mismatch number of errors
}

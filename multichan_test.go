/**
 * Copyright 2019 Innodev LLC. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package multichan

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultichan(t *testing.T) {
	mc := NewMultichan()

	item := "foobar"
	routines := 40
	rcvChan := make(chan bool)
	for i := 0; i < routines; i++ {
		go func() {
			val, err := mc.Receive(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, item, val)
			rcvChan <- true
		}()
	}
	for mc.listeners < int64(routines)-4 {
		time.Sleep(50 * time.Millisecond)
	}
	mc.Send(item)
	for i := 0; i < routines; i++ {
		select {
		case <-rcvChan:
		case <-time.After(20 * time.Second):
			t.FailNow()
		}

	}

}

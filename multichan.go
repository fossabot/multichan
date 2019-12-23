/**
 * Copyright 2019 Innodev LLC. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package multichan

import (
	"context"
	"sync"
)

type Multichan struct {
	data      interface{}
	mux       *sync.Mutex
	listeners int64
	received  bool
	once      *sync.Once
	ch        chan interface{}
}

func NewMultichan() *Multichan {
	return &Multichan{
		received: false,
		once:     &sync.Once{},
		mux:      &sync.Mutex{},
		ch:       make(chan interface{}),
	}
}

func (mc *Multichan) release(left int64) {
	for i := int64(0); i < left; i++ {
		mc.ch <- mc.data
	}
}

func (mc *Multichan) Receive(ctx context.Context) (interface{}, error) {
	mc.mux.Lock()
	if mc.received {
		mc.mux.Unlock()
		return mc.data, nil
	}
	mc.listeners++
	mc.mux.Unlock()
	var val interface{}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case val = <-mc.ch:
	}

	mc.mux.Lock()
	mc.listeners--
	if !mc.received {
		mc.data = val
		mc.received = true
		mc.mux.Unlock()
		go mc.release(mc.listeners)
		return val, nil
	}
	mc.mux.Unlock()
	return val, nil
}

// Send gives data for consumption on everything receiving, calls after the first call do nothing
func (mc *Multichan) Send(val interface{}) {
	mc.once.Do(func() {
		mc.ch <- val
	})

}

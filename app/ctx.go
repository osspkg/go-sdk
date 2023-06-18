/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package app

import "context"

type (
	_ctx struct {
		ctx    context.Context
		cancel context.CancelFunc
	}

	//Context model for force close application
	Context interface {
		Close()
		Context() context.Context
		Done() <-chan struct{}
	}
)

func NewContext() Context {
	ctx, cancel := context.WithCancel(context.Background())

	return &_ctx{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Close context close method
func (v *_ctx) Close() {
	v.cancel()
}

// Context general context
func (v *_ctx) Context() context.Context {
	return v.ctx
}

// Done context close wait channel
func (v *_ctx) Done() <-chan struct{} {
	return v.ctx.Done()
}

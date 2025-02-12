// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-18, by liasica

package pool

import (
	"bytes"
	"sync"
)

// buffer pool to reduce GC
var buffers = sync.Pool{
	// New is called when a new instance is needed
	New: func() any {
		return new(bytes.Buffer)
	},
}

// GetBuffer fetches a buffer from the pool
func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *bytes.Buffer) {
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16
	if buf != nil && buf.Cap() < maxSize {
		buf.Reset()
		buffers.Put(buf)
	}
}

type Pool[T any] struct {
	sync.Pool
}

func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.Pool.Put(x)
}

func NewPool[T any](f func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() any {
				return f()
			},
		},
	}
}

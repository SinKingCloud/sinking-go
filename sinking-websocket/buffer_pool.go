package sinking_websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

var sharedWriteBufferPools sync.Map

type writeBufferPool struct {
	size int
	pool sync.Pool
}

func sharedWriteBufferPool(size int) websocket.BufferPool {
	size = resolvedWriteBufferSize(size)

	if pool, ok := sharedWriteBufferPools.Load(size); ok {
		return pool.(websocket.BufferPool)
	}

	pool := &writeBufferPool{
		size: size,
	}
	actual, _ := sharedWriteBufferPools.LoadOrStore(size, pool)
	return actual.(websocket.BufferPool)
}

func (pool *writeBufferPool) Get() interface{} {
	if value := pool.pool.Get(); value != nil {
		return value
	}
	return make([]byte, pool.size)
}

func (pool *writeBufferPool) Put(value interface{}) {
	buffer, ok := value.([]byte)
	if !ok {
		return
	}
	if cap(buffer) < pool.size {
		return
	}
	pool.pool.Put(buffer[:pool.size])
}

func resolvedWriteBufferSize(size int) int {
	if size > 0 {
		return size
	}
	return defaultWriteBufferSize
}

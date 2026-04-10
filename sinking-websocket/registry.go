package sinking_websocket

import (
	"hash/maphash"
	"sync"
	"sync/atomic"
)

type registryShard struct {
	mu    sync.RWMutex
	items map[string]*Connection
}

type registryEntry struct {
	id         string
	connection *Connection
}

type Registry struct {
	once   sync.Once
	seed   maphash.Seed
	size   int64
	shards []registryShard
}

func NewRegistry() *Registry {
	registry := &Registry{}
	registry.init()
	return registry
}

func (registry *Registry) Load(id string) (*Connection, bool) {
	shard := registry.shard(id)
	shard.mu.RLock()
	connection, ok := shard.items[id]
	shard.mu.RUnlock()
	return connection, ok
}

func (registry *Registry) Store(id string, connection *Connection) {
	if connection == nil {
		registry.Delete(id)
		return
	}

	var replaced *Connection

	shard := registry.shard(id)
	shard.mu.Lock()
	replaced = shard.items[id]
	if replaced == nil {
		atomic.AddInt64(&registry.size, 1)
	}
	shard.items[id] = connection
	shard.mu.Unlock()

	if replaced != nil && replaced != connection {
		_ = replaced.Close()
	}
}

func (registry *Registry) Delete(id string) bool {
	return registry.delete(id, nil)
}

func (registry *Registry) DeleteIfMatch(id string, connection *Connection) bool {
	if connection == nil {
		return false
	}
	return registry.delete(id, connection)
}

func (registry *Registry) Range(visitor func(id string, connection *Connection) bool) {
	if visitor == nil {
		return
	}

	registry.init()

	for i := range registry.shards {
		shard := &registry.shards[i]

		shard.mu.RLock()
		snapshot := make([]registryEntry, 0, len(shard.items))
		for id, connection := range shard.items {
			snapshot = append(snapshot, registryEntry{
				id:         id,
				connection: connection,
			})
		}
		shard.mu.RUnlock()

		for _, entry := range snapshot {
			if !visitor(entry.id, entry.connection) {
				return
			}
		}
	}
}

func (registry *Registry) Len() int {
	if registry == nil {
		return 0
	}
	return int(atomic.LoadInt64(&registry.size))
}

func (registry *Registry) Close() {
	registry.Range(func(id string, connection *Connection) bool {
		registry.DeleteIfMatch(id, connection)
		return true
	})
}

func (registry *Registry) delete(id string, expected *Connection) bool {
	var removed *Connection

	shard := registry.shard(id)
	shard.mu.Lock()
	current, ok := shard.items[id]
	if !ok {
		shard.mu.Unlock()
		return true
	}
	if expected != nil && current != expected {
		shard.mu.Unlock()
		return false
	}
	delete(shard.items, id)
	atomic.AddInt64(&registry.size, -1)
	removed = current
	shard.mu.Unlock()

	if removed != nil {
		_ = removed.Close()
	}

	return true
}

func (registry *Registry) init() {
	if registry == nil {
		return
	}

	registry.once.Do(func() {
		registry.seed = maphash.MakeSeed()
		registry.shards = make([]registryShard, defaultRegistryShards)
		for i := range registry.shards {
			registry.shards[i].items = make(map[string]*Connection)
		}
	})
}

func (registry *Registry) shard(id string) *registryShard {
	registry.init()
	index := maphash.String(registry.seed, id) % uint64(len(registry.shards))
	return &registry.shards[index]
}

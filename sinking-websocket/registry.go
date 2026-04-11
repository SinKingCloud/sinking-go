package sinking_websocket

import (
	"encoding/json"
	"hash/maphash"
	"runtime"
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

var registryEntriesPool = sync.Pool{
	New: func() interface{} {
		return make([]registryEntry, 0, 64)
	},
}

func acquireRegistryEntries(sizeHint int) []registryEntry {
	entries := registryEntriesPool.Get().([]registryEntry)
	if cap(entries) < sizeHint {
		return make([]registryEntry, 0, sizeHint)
	}
	return entries[:0]
}

func releaseRegistryEntries(entries []registryEntry) {
	if cap(entries) > 4096 {
		return
	}
	registryEntriesPool.Put(entries[:0])
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
		snapshot := acquireRegistryEntries(len(shard.items))
		for id, connection := range shard.items {
			snapshot = append(snapshot, registryEntry{
				id:         id,
				connection: connection,
			})
		}
		shard.mu.RUnlock()

		for _, entry := range snapshot {
			if !visitor(entry.id, entry.connection) {
				releaseRegistryEntries(snapshot)
				return
			}
		}

		releaseRegistryEntries(snapshot)
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

func (registry *Registry) Broadcast(messageType int, payload []byte) (BroadcastResult, error) {
	prepared, err := PrepareMessage(messageType, payload)
	if err != nil {
		return BroadcastResult{}, err
	}
	return registry.BroadcastPrepared(prepared), nil
}

func (registry *Registry) BroadcastJSON(value interface{}) (BroadcastResult, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return BroadcastResult{}, err
	}
	return registry.Broadcast(TextMessage, payload)
}

func (registry *Registry) BroadcastPrepared(message *PreparedMessage) BroadcastResult {
	if registry == nil || message == nil {
		return BroadcastResult{}
	}

	registry.init()

	var result BroadcastResult
	failed := acquireRegistryEntries(0)

	for i := range registry.shards {
		shard := &registry.shards[i]

		shard.mu.RLock()
		for id, connection := range shard.items {
			switch err := connection.TrySendPrepared(message); err {
			case nil:
				result.Queued++
			case ErrSendQueueFull:
				result.Dropped++
			case ErrConnectionClosed:
				result.Closed++
				failed = append(failed, registryEntry{
					id:         id,
					connection: connection,
				})
			default:
				result.Closed++
				failed = append(failed, registryEntry{
					id:         id,
					connection: connection,
				})
			}
		}
		shard.mu.RUnlock()
	}

	for _, entry := range failed {
		registry.DeleteIfMatch(entry.id, entry.connection)
	}

	releaseRegistryEntries(failed)
	return result
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
	defaultGroupShards := func() int {
		target := runtime.GOMAXPROCS(0) * 4
		if target <= 0 {
			target = runtime.NumCPU() * 8
		}
		if target <= 0 {
			target = 8
		}
		shards := 1
		for shards < target {
			shards <<= 1
		}
		return shards
	}
	registry.once.Do(func() {
		registry.seed = maphash.MakeSeed()
		registry.shards = make([]registryShard, defaultGroupShards())
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

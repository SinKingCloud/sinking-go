package sinking_websocket

import (
	"hash/maphash"
	"sync"
)

// NewConnections 初始化连接池。
func NewConnections() *ConnectionPool {
	pool := &ConnectionPool{}
	pool.init()
	return pool
}

type connectionShard struct {
	conn map[string]*Connection
	lock sync.RWMutex
}

// ConnectionPool 管理当前在线的 websocket 连接。
type ConnectionPool struct {
	once   sync.Once
	seed   maphash.Seed
	shards []connectionShard
}

// init 初始化连接池分片。
func (connections *ConnectionPool) init() {
	connections.once.Do(func() {
		connections.seed = maphash.MakeSeed()
		connections.shards = make([]connectionShard, defaultConnectionPoolShards)
		for i := range connections.shards {
			connections.shards[i].conn = make(map[string]*Connection)
		}
	})
}

// shard 返回 key 所在的连接池分片。
func (connections *ConnectionPool) shard(key string) *connectionShard {
	connections.init()
	index := maphash.String(connections.seed, key) % uint64(len(connections.shards))
	return &connections.shards[index]
}

// Get 获取长连接对象。
func (connections *ConnectionPool) Get(key string) *Connection {
	shard := connections.shard(key)
	shard.lock.RLock()
	defer shard.lock.RUnlock()
	return shard.conn[key]
}

// GetAll 获取所有长连接对象。
func (connections *ConnectionPool) GetAll() map[string]*Connection {
	connections.init()

	conn := make(map[string]*Connection)
	for i := range connections.shards {
		shard := &connections.shards[i]
		shard.lock.RLock()
		for k, v := range shard.conn {
			conn[k] = v
		}
		shard.lock.RUnlock()
	}
	return conn
}

// Set 设置长连接对象。
func (connections *ConnectionPool) Set(key string, conn *Connection) {
	var oldConn *Connection

	shard := connections.shard(key)
	shard.lock.Lock()
	oldConn = shard.conn[key]
	shard.conn[key] = conn
	shard.lock.Unlock()

	if oldConn != nil && oldConn != conn {
		_ = oldConn.Close()
	}
}

// Delete 删除长连接对象。
// 传入连接实例时，仅在 key 当前绑定的连接与实例一致时才删除，
// 可避免新连接被旧连接的关闭事件误删。
func (connections *ConnectionPool) Delete(key string, expected ...*Connection) bool {
	var expectedConn *Connection
	if len(expected) > 0 {
		expectedConn = expected[0]
	}
	shard := connections.shard(key)
	shard.lock.Lock()
	current, exists := shard.conn[key]
	if !exists {
		shard.lock.Unlock()
		return true
	}
	if expectedConn != nil && current != expectedConn {
		shard.lock.Unlock()
		return false
	}
	delete(shard.conn, key)
	shard.lock.Unlock()
	if current == nil || current.IsClosed() {
		return true
	}
	if err := current.Close(); err != nil {
		return false
	}
	return true
}

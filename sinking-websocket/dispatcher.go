package sinking_websocket

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultMinWriteDispatcherShards = 16
	defaultMaxWriteDispatcherShards = 256
	defaultHeartbeatSweepInterval   = time.Second
	minHeartbeatSweepInterval       = 100 * time.Millisecond
)

type writeDispatcherConfig struct {
	shardCount             int
	pingInterval           time.Duration
	heartbeatSweepInterval time.Duration
}

type writeDispatcher struct {
	shards []*writeDispatcherShard
	next   uint64
}

type writeDispatcherShard struct {
	queueMu   sync.Mutex
	queueCond *sync.Cond
	ready     []*Connection
	head      int

	heartbeatMu          sync.RWMutex
	heartbeatConnections map[*Connection]struct{}
	pingInterval         time.Duration
}

var connectionPool = sync.Pool{
	New: func() interface{} {
		return make([]*Connection, 0, 64)
	},
}

func newWriteDispatcher(config writeDispatcherConfig) *writeDispatcher {
	shardCount := resolvedWriteDispatcherShards(config.shardCount)

	dispatcher := &writeDispatcher{
		shards: make([]*writeDispatcherShard, shardCount),
	}

	for i := 0; i < shardCount; i++ {
		shard := &writeDispatcherShard{
			pingInterval: config.pingInterval,
		}
		shard.queueCond = sync.NewCond(&shard.queueMu)
		if config.pingInterval > 0 {
			shard.heartbeatConnections = make(map[*Connection]struct{})
		}
		dispatcher.shards[i] = shard
		go shard.run()
		if config.pingInterval > 0 {
			go shard.runHeartbeats(resolvedHeartbeatSweepInterval(config.heartbeatSweepInterval, config.pingInterval))
		}
	}

	return dispatcher
}

func (dispatcher *writeDispatcher) bind(connection *Connection) *writeDispatcherShard {
	if dispatcher == nil || len(dispatcher.shards) == 0 {
		return nil
	}

	index := int((atomic.AddUint64(&dispatcher.next, 1) - 1) % uint64(len(dispatcher.shards)))
	shard := dispatcher.shards[index]
	shard.register(connection)
	return shard
}

func (shard *writeDispatcherShard) register(connection *Connection) {
	if shard == nil || connection == nil || shard.heartbeatConnections == nil {
		return
	}

	shard.heartbeatMu.Lock()
	shard.heartbeatConnections[connection] = struct{}{}
	shard.heartbeatMu.Unlock()
}

func (shard *writeDispatcherShard) unregister(connection *Connection) {
	if shard == nil || connection == nil || shard.heartbeatConnections == nil {
		return
	}

	shard.heartbeatMu.Lock()
	delete(shard.heartbeatConnections, connection)
	shard.heartbeatMu.Unlock()
}

func (shard *writeDispatcherShard) schedule(connection *Connection) {
	if shard == nil || connection == nil {
		return
	}

	shard.queueMu.Lock()
	shard.ready = append(shard.ready, connection)
	shard.queueMu.Unlock()
	shard.queueCond.Signal()
}

func (shard *writeDispatcherShard) run() {
	for {
		connection := shard.waitForReadyConnection()
		if connection == nil {
			continue
		}
		connection.flushPending()
	}
}

func (shard *writeDispatcherShard) waitForReadyConnection() *Connection {
	shard.queueMu.Lock()
	defer shard.queueMu.Unlock()

	for shard.head >= len(shard.ready) {
		shard.ready = shard.ready[:0]
		shard.head = 0
		shard.queueCond.Wait()
	}

	connection := shard.ready[shard.head]
	shard.ready[shard.head] = nil
	shard.head++

	if shard.head == len(shard.ready) {
		shard.ready = shard.ready[:0]
		shard.head = 0
	} else if shard.head >= 1024 && shard.head*2 >= len(shard.ready) {
		copy(shard.ready, shard.ready[shard.head:])
		tail := len(shard.ready) - shard.head
		for i := tail; i < len(shard.ready); i++ {
			shard.ready[i] = nil
		}
		shard.ready = shard.ready[:tail]
		shard.head = 0
	}

	return connection
}

func (shard *writeDispatcherShard) runHeartbeats(sweepInterval time.Duration) {
	ticker := time.NewTicker(sweepInterval)
	defer ticker.Stop()

	for now := range ticker.C {
		connections := shard.snapshotHeartbeatConnections()
		for _, connection := range connections {
			if connection == nil || connection.Closed() {
				continue
			}

			nextPingAt := atomic.LoadInt64(&connection.nextPingAt)
			if nextPingAt == 0 || nextPingAt > now.UnixNano() {
				continue
			}

			if err := connection.writeControl(PingMessage, nil); err != nil {
				_ = connection.Close()
				continue
			}

			atomic.StoreInt64(&connection.nextPingAt, now.Add(shard.pingInterval).UnixNano())
		}
		releaseConnections(connections)
	}
}

func (shard *writeDispatcherShard) snapshotHeartbeatConnections() []*Connection {
	if shard == nil || shard.heartbeatConnections == nil {
		return nil
	}

	shard.heartbeatMu.RLock()
	connections := acquireConnections(len(shard.heartbeatConnections))
	for connection := range shard.heartbeatConnections {
		connections = append(connections, connection)
	}
	shard.heartbeatMu.RUnlock()
	return connections
}

func acquireConnections(sizeHint int) []*Connection {
	connections := connectionPool.Get().([]*Connection)
	if cap(connections) < sizeHint {
		return make([]*Connection, 0, sizeHint)
	}
	return connections[:0]
}

func releaseConnections(connections []*Connection) {
	if connections == nil || cap(connections) > 4096 {
		return
	}
	for i := range connections {
		connections[i] = nil
	}
	connectionPool.Put(connections[:0])
}

func resolvedWriteDispatcherShards(shardCount int) int {
	if shardCount > 0 {
		return shardCount
	}

	resolved := runtime.GOMAXPROCS(0) * 8
	if resolved < defaultMinWriteDispatcherShards {
		return defaultMinWriteDispatcherShards
	}
	if resolved > defaultMaxWriteDispatcherShards {
		return defaultMaxWriteDispatcherShards
	}
	return resolved
}

func resolvedHeartbeatSweepInterval(configured, pingInterval time.Duration) time.Duration {
	if configured > 0 {
		return configured
	}
	if pingInterval <= 0 {
		return defaultHeartbeatSweepInterval
	}

	interval := pingInterval / 4
	if interval < minHeartbeatSweepInterval {
		return minHeartbeatSweepInterval
	}
	if interval > defaultHeartbeatSweepInterval {
		return defaultHeartbeatSweepInterval
	}
	return interval
}

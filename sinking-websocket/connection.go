package sinking_websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type queuedWriteKind uint8

const (
	queuedMessageWrite queuedWriteKind = iota
	queuedPreparedWrite
)

type queuedWrite struct {
	kind        queuedWriteKind
	messageType int
	payload     []byte
	prepared    *PreparedMessage
}

type connectionConfig struct {
	writeTimeout   time.Duration
	pingInterval   time.Duration
	writeQueueSize int
	dispatcher     *writeDispatcher
}

type Connection struct {
	id             string
	request        *http.Request
	raw            *websocket.Conn
	writeTimeout   time.Duration
	writeQueueSize int
	dispatcher     *writeDispatcherShard
	nextPingAt     int64

	pendingMu sync.Mutex
	spaceCond *sync.Cond
	pending   []queuedWrite
	scheduled bool
	closed    bool

	done      chan struct{}
	closeOnce sync.Once
	closeErr  error
}

func newConnection(id string, request *http.Request, raw *websocket.Conn, config connectionConfig) *Connection {
	connection := &Connection{
		id:             id,
		request:        request,
		raw:            raw,
		writeTimeout:   config.writeTimeout,
		writeQueueSize: resolvedWriteQueueSize(config.writeQueueSize),
		done:           make(chan struct{}),
	}

	if config.dispatcher != nil {
		connection.dispatcher = config.dispatcher.bind(connection)
	}

	if config.pingInterval > 0 {
		atomic.StoreInt64(&connection.nextPingAt, time.Now().Add(config.pingInterval).UnixNano())
	}

	return connection
}

func (connection *Connection) ID() string {
	if connection == nil {
		return ""
	}
	return connection.id
}

func (connection *Connection) Request() *http.Request {
	if connection == nil {
		return nil
	}
	return connection.request
}

func (connection *Connection) Done() <-chan struct{} {
	if connection == nil {
		return nil
	}
	return connection.done
}

func (connection *Connection) Closed() bool {
	if connection == nil {
		return true
	}

	select {
	case <-connection.done:
		return true
	default:
		return false
	}
}

func (connection *Connection) Send(messageType int, payload []byte) error {
	return connection.enqueue(queuedWrite{
		kind:        queuedMessageWrite,
		messageType: messageType,
		payload:     cloneBytes(payload),
	}, true)
}

func (connection *Connection) TrySend(messageType int, payload []byte) error {
	return connection.enqueue(queuedWrite{
		kind:        queuedMessageWrite,
		messageType: messageType,
		payload:     cloneBytes(payload),
	}, false)
}

func (connection *Connection) SendJSON(value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return connection.enqueue(queuedWrite{
		kind:        queuedMessageWrite,
		messageType: TextMessage,
		payload:     payload,
	}, true)
}

func (connection *Connection) TrySendJSON(value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return connection.enqueue(queuedWrite{
		kind:        queuedMessageWrite,
		messageType: TextMessage,
		payload:     payload,
	}, false)
}

func (connection *Connection) SendPrepared(message *PreparedMessage) error {
	if message == nil {
		return errNilPreparedMessage
	}
	return connection.enqueue(queuedWrite{
		kind:     queuedPreparedWrite,
		prepared: message,
	}, true)
}

func (connection *Connection) TrySendPrepared(message *PreparedMessage) error {
	if message == nil {
		return errNilPreparedMessage
	}
	return connection.enqueue(queuedWrite{
		kind:     queuedPreparedWrite,
		prepared: message,
	}, false)
}

func (connection *Connection) Close() error {
	if connection == nil {
		return nil
	}

	connection.closeOnce.Do(func() {
		connection.pendingMu.Lock()
		connection.closed = true
		releaseQueuedWrites(connection.pending)
		connection.pending = nil
		if connection.spaceCond != nil {
			connection.spaceCond.Broadcast()
		}
		connection.pendingMu.Unlock()

		close(connection.done)

		if connection.dispatcher != nil {
			connection.dispatcher.unregister(connection)
		}

		if connection.raw != nil {
			connection.closeErr = connection.raw.Close()
		}
	})

	return connection.closeErr
}

func (connection *Connection) enqueue(write queuedWrite, blocking bool) error {
	if connection == nil || connection.raw == nil {
		return errNilConnection
	}

	shouldSchedule := false

	connection.pendingMu.Lock()
	for {
		if connection.closed {
			connection.pendingMu.Unlock()
			return ErrConnectionClosed
		}
		if len(connection.pending) < connection.writeQueueSize {
			connection.pending = appendQueuedWrite(connection.pending, write)
			if !connection.scheduled {
				connection.scheduled = true
				shouldSchedule = true
			}
			connection.pendingMu.Unlock()

			if shouldSchedule && connection.dispatcher != nil {
				connection.dispatcher.schedule(connection)
			}

			return nil
		}
		if !blocking {
			connection.pendingMu.Unlock()
			return ErrSendQueueFull
		}
		if connection.spaceCond == nil {
			connection.spaceCond = sync.NewCond(&connection.pendingMu)
		}
		connection.spaceCond.Wait()
	}
}

func (connection *Connection) flushPending() {
	for {
		batch := connection.takePendingBatch()
		if len(batch) == 0 {
			return
		}

		for _, write := range batch {
			if err := connection.write(write); err != nil {
				releaseQueuedWrites(batch)
				_ = connection.Close()
				return
			}
		}

		releaseQueuedWrites(batch)
	}
}

func (connection *Connection) takePendingBatch() []queuedWrite {
	connection.pendingMu.Lock()
	defer connection.pendingMu.Unlock()

	if len(connection.pending) == 0 {
		connection.scheduled = false
		return nil
	}

	batch := connection.pending
	connection.pending = nil
	if connection.spaceCond != nil {
		connection.spaceCond.Broadcast()
	}
	return batch
}

func (connection *Connection) write(write queuedWrite) error {
	if connection == nil || connection.raw == nil {
		return errNilConnection
	}

	if err := connection.raw.SetWriteDeadline(connection.writeDeadline()); err != nil {
		return err
	}

	switch write.kind {
	case queuedMessageWrite:
		return connection.raw.WriteMessage(write.messageType, write.payload)
	case queuedPreparedWrite:
		return connection.raw.WritePreparedMessage(write.prepared)
	default:
		return nil
	}
}

func (connection *Connection) writeControl(messageType int, payload []byte) error {
	if connection == nil || connection.raw == nil {
		return errNilConnection
	}
	return connection.raw.WriteControl(messageType, payload, connection.writeDeadline())
}

func (connection *Connection) writeDeadline() time.Time {
	if connection.writeTimeout <= 0 {
		return time.Time{}
	}
	return time.Now().Add(connection.writeTimeout)
}

var queuedWritePool = sync.Pool{
	New: func() interface{} {
		return make([]queuedWrite, 0, 8)
	},
}

func appendQueuedWrite(batch []queuedWrite, write queuedWrite) []queuedWrite {
	if batch == nil {
		batch = acquireQueuedWrites(1)
	}
	return append(batch, write)
}

func acquireQueuedWrites(sizeHint int) []queuedWrite {
	batch := queuedWritePool.Get().([]queuedWrite)
	if cap(batch) < sizeHint {
		return make([]queuedWrite, 0, sizeHint)
	}
	return batch[:0]
}

func releaseQueuedWrites(batch []queuedWrite) {
	if batch == nil {
		return
	}
	if cap(batch) > 1024 {
		return
	}
	for i := range batch {
		batch[i] = queuedWrite{}
	}
	queuedWritePool.Put(batch[:0])
}

func cloneBytes(payload []byte) []byte {
	if len(payload) == 0 {
		return nil
	}

	cloned := make([]byte, len(payload))
	copy(cloned, payload)
	return cloned
}

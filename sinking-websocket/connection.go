package sinking_websocket

import (
	"encoding/json"
	"net/http"
	"sync"
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
}

type Connection struct {
	id           string
	request      *http.Request
	raw          *websocket.Conn
	writeTimeout time.Duration
	outgoing     chan queuedWrite
	done         chan struct{}
	closeOnce    sync.Once
	closeErr     error
}

func newConnection(id string, request *http.Request, raw *websocket.Conn, config connectionConfig) *Connection {
	queueSize := config.writeQueueSize
	if queueSize <= 0 {
		queueSize = defaultWriteQueueSize
	}

	connection := &Connection{
		id:           id,
		request:      request,
		raw:          raw,
		writeTimeout: config.writeTimeout,
		outgoing:     make(chan queuedWrite, queueSize),
		done:         make(chan struct{}),
	}

	go connection.writeLoop(config.pingInterval)
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
	return connection.Send(TextMessage, payload)
}

func (connection *Connection) TrySendJSON(value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return connection.TrySend(TextMessage, payload)
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
		close(connection.done)
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

	if blocking {
		select {
		case <-connection.done:
			return ErrConnectionClosed
		case connection.outgoing <- write:
			return nil
		}
	}

	select {
	case <-connection.done:
		return ErrConnectionClosed
	case connection.outgoing <- write:
		return nil
	default:
		return ErrSendQueueFull
	}
}

func (connection *Connection) writeLoop(pingInterval time.Duration) {
	var ticker *time.Ticker

	if pingInterval > 0 {
		ticker = time.NewTicker(pingInterval)
		defer ticker.Stop()
	}

	for {
		select {
		case <-connection.done:
			return
		case write := <-connection.outgoing:
			if err := connection.write(write); err != nil {
				_ = connection.Close()
				return
			}
		case <-tickerChannel(ticker):
			if err := connection.writeControl(PingMessage, nil); err != nil {
				_ = connection.Close()
				return
			}
		}
	}
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

func cloneBytes(payload []byte) []byte {
	if len(payload) == 0 {
		return nil
	}

	cloned := make([]byte, len(payload))
	copy(cloned, payload)
	return cloned
}

func tickerChannel(ticker *time.Ticker) <-chan time.Time {
	if ticker == nil {
		return nil
	}
	return ticker.C
}

package sinking_websocket

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultConnectionPoolShards = 64
	defaultWriteQueueSize       = 256
)

var (
	errNilConnection    = errors.New("websocket connection is nil")
	ErrConnectionClosed = errors.New("websocket connection is closed")
	ErrSendQueueFull    = errors.New("websocket send queue is full")
)

type writeRequestKind uint8

const (
	writeMessageRequest writeRequestKind = iota
	writePreparedMessageRequest
)

type writeRequest struct {
	kind        writeRequestKind
	messageType int
	data        []byte
	prepared    *websocket.PreparedMessage
}

// Connection 对 gorilla websocket 连接做了一层并发、关闭和异步发送保护。
type Connection struct {
	*websocket.Conn
	stateOnce  sync.Once
	writerOnce sync.Once
	writeLock  sync.Mutex
	closeOnce  sync.Once
	closeErr   error
	done       chan struct{}
	sendQueue  chan writeRequest
	writeWait  time.Duration
	pingPeriod time.Duration
	queueSize  int
}

// NewConnection 包装原始 websocket 连接。
func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		Conn: conn,
	}
}

// initState 初始化连接内部状态。
func (connection *Connection) initState() {
	connection.stateOnce.Do(func() {
		connection.done = make(chan struct{})
	})
}

// configureAsyncWrite 配置写超时、写队列和心跳参数。
func (connection *Connection) configureAsyncWrite(writeWait, pingPeriod time.Duration, queueSize int) {
	connection.initState()
	connection.writeWait = writeWait
	connection.pingPeriod = pingPeriod
	if queueSize > 0 {
		connection.queueSize = queueSize
	}
}

// startWriter 启动单连接写协程，负责异步消息发送和 ping 心跳。
func (connection *Connection) startWriter() {
	connection.initState()
	connection.writerOnce.Do(func() {
		if connection.queueSize <= 0 {
			connection.queueSize = defaultWriteQueueSize
		}
		connection.sendQueue = make(chan writeRequest, connection.queueSize)
		go connection.writeLoop()
	})
}

// withWriteLock 串行化普通写入，避免 gorilla websocket 并发写 panic。
func (connection *Connection) withWriteLock(fun func() error) error {
	if connection == nil || connection.Conn == nil {
		return errNilConnection
	}

	connection.writeLock.Lock()
	defer connection.writeLock.Unlock()
	return fun()
}

// writeMessageNow 直接向底层连接写入消息。
func (connection *Connection) writeMessageNow(messageType int, data []byte) error {
	return connection.withWriteLock(func() error {
		if connection.writeWait > 0 {
			_ = connection.Conn.SetWriteDeadline(time.Now().Add(connection.writeWait))
		}
		return connection.Conn.WriteMessage(messageType, data)
	})
}

// writeJSONNow 直接向底层连接写入 JSON 消息。
func (connection *Connection) writeJSONNow(v interface{}) error {
	return connection.withWriteLock(func() error {
		if connection.writeWait > 0 {
			_ = connection.Conn.SetWriteDeadline(time.Now().Add(connection.writeWait))
		}
		return connection.Conn.WriteJSON(v)
	})
}

// writePreparedMessageNow 直接向底层连接写入 PreparedMessage。
func (connection *Connection) writePreparedMessageNow(pm *websocket.PreparedMessage) error {
	return connection.withWriteLock(func() error {
		if connection.writeWait > 0 {
			_ = connection.Conn.SetWriteDeadline(time.Now().Add(connection.writeWait))
		}
		return connection.Conn.WritePreparedMessage(pm)
	})
}

// writeLoop 串行处理异步发送队列和连接心跳。
func (connection *Connection) writeLoop() {
	var ticker *time.Ticker

	if connection.pingPeriod > 0 {
		ticker = time.NewTicker(connection.pingPeriod)
		defer ticker.Stop()
	}

	for {
		select {
		case <-connection.done:
			return
		case request := <-connection.sendQueue:
			if err := connection.handleWriteRequest(request); err != nil {
				_ = connection.Close()
				return
			}
		case <-tickerChannel(ticker):
			if err := connection.WriteControl(websocket.PingMessage, nil, time.Now().Add(connection.writeWait)); err != nil {
				_ = connection.Close()
				return
			}
		}
	}
}

// handleWriteRequest 执行异步写请求。
func (connection *Connection) handleWriteRequest(request writeRequest) error {
	switch request.kind {
	case writeMessageRequest:
		return connection.writeMessageNow(request.messageType, request.data)
	case writePreparedMessageRequest:
		return connection.writePreparedMessageNow(request.prepared)
	default:
		return nil
	}
}

// enqueueWrite 将写请求投递到异步发送队列。
func (connection *Connection) enqueueWrite(request writeRequest, blocking bool) error {
	if connection == nil || connection.Conn == nil {
		return errNilConnection
	}

	connection.startWriter()

	select {
	case <-connection.done:
		return ErrConnectionClosed
	default:
	}

	if blocking {
		select {
		case connection.sendQueue <- request:
			return nil
		case <-connection.done:
			return ErrConnectionClosed
		}
	}

	select {
	case connection.sendQueue <- request:
		return nil
	case <-connection.done:
		return ErrConnectionClosed
	default:
		return ErrSendQueueFull
	}
}

// cloneBytes 复制异步发送使用的消息体，避免调用方后续修改原始切片。
func cloneBytes(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}

	cloned := make([]byte, len(data))
	copy(cloned, data)
	return cloned
}

// tickerChannel 返回 ticker 的只读通道，便于在未启用心跳时关闭该分支。
func tickerChannel(ticker *time.Ticker) <-chan time.Time {
	if ticker == nil {
		return nil
	}
	return ticker.C
}

// WriteMessage 同步写入业务消息。
func (connection *Connection) WriteMessage(messageType int, data []byte) error {
	return connection.writeMessageNow(messageType, data)
}

// WriteJSON 同步写入 JSON。
func (connection *Connection) WriteJSON(v interface{}) error {
	return connection.writeJSONNow(v)
}

// WritePreparedMessage 同步写入 PreparedMessage。
func (connection *Connection) WritePreparedMessage(pm *websocket.PreparedMessage) error {
	return connection.writePreparedMessageNow(pm)
}

// SendMessage 将消息异步投递到发送队列，适合高并发推送场景。
func (connection *Connection) SendMessage(messageType int, data []byte) error {
	return connection.enqueueWrite(writeRequest{
		kind:        writeMessageRequest,
		messageType: messageType,
		data:        cloneBytes(data),
	}, true)
}

// TrySendMessage 非阻塞地将消息投递到发送队列，队列满时立即返回错误。
func (connection *Connection) TrySendMessage(messageType int, data []byte) error {
	return connection.enqueueWrite(writeRequest{
		kind:        writeMessageRequest,
		messageType: messageType,
		data:        cloneBytes(data),
	}, false)
}

// SendJSON 序列化 JSON 后异步投递到发送队列。
func (connection *Connection) SendJSON(v interface{}) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return connection.SendMessage(websocket.TextMessage, payload)
}

// TrySendJSON 非阻塞地序列化 JSON 并投递到发送队列。
func (connection *Connection) TrySendJSON(v interface{}) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return connection.TrySendMessage(websocket.TextMessage, payload)
}

// SendPreparedMessage 将 PreparedMessage 异步投递到发送队列。
func (connection *Connection) SendPreparedMessage(pm *websocket.PreparedMessage) error {
	return connection.enqueueWrite(writeRequest{
		kind:     writePreparedMessageRequest,
		prepared: pm,
	}, true)
}

// TrySendPreparedMessage 非阻塞地将 PreparedMessage 投递到发送队列。
func (connection *Connection) TrySendPreparedMessage(pm *websocket.PreparedMessage) error {
	return connection.enqueueWrite(writeRequest{
		kind:     writePreparedMessageRequest,
		prepared: pm,
	}, false)
}

// WriteControl 直接写入控制帧。
func (connection *Connection) WriteControl(messageType int, data []byte, deadline time.Time) error {
	if connection == nil || connection.Conn == nil {
		return errNilConnection
	}
	return connection.Conn.WriteControl(messageType, data, deadline)
}

// Close 保证连接只会被真正关闭一次。
func (connection *Connection) Close() error {
	if connection == nil {
		return nil
	}

	connection.initState()
	connection.closeOnce.Do(func() {
		close(connection.done)
		if connection.sendQueue != nil {
			close(connection.sendQueue)
		}
		if connection.Conn != nil {
			connection.closeErr = connection.Conn.Close()
		}
	})
	return connection.closeErr
}

// IsClosed 是否关闭
func (connection *Connection) IsClosed() bool {
	select {
	case <-connection.done:
		return true
	default:
		return false
	}
}

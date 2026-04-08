package sinking_websocket

import (
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var errNilConnection = errors.New("websocket connection is nil")

// NewConnections 初始化连接池
func NewConnections() *ConnectionPool {
	return &ConnectionPool{
		conn: make(map[string]*Connection),
	}
}

// ConnectionPool 管理当前在线的 websocket 连接。
type ConnectionPool struct {
	conn map[string]*Connection
	lock sync.RWMutex
}

// Get 获取长连接对象
func (connections *ConnectionPool) Get(key string) *Connection {
	connections.lock.RLock()
	defer connections.lock.RUnlock()
	return connections.conn[key]
}

// GetAll 获取所有长连接对象
func (connections *ConnectionPool) GetAll() map[string]*Connection {
	connections.lock.RLock()
	defer connections.lock.RUnlock()

	conn := make(map[string]*Connection)
	for k, v := range connections.conn {
		conn[k] = v
	}
	return conn
}

// Set 设置长连接对象
func (connections *ConnectionPool) Set(key string, conn *Connection) {
	var oldConn *Connection

	connections.lock.Lock()
	if connections.conn == nil {
		connections.conn = make(map[string]*Connection)
	}
	oldConn = connections.conn[key]
	connections.conn[key] = conn
	connections.lock.Unlock()

	if oldConn != nil && oldConn != conn {
		_ = oldConn.Close()
	}
}

// Delete 删除长连接对象。
// 传入连接实例时，仅在 key 当前绑定的连接与实例一致时才删除，
// 可避免新连接被旧连接的关闭事件误删。
func (connections *ConnectionPool) Delete(key string, expected ...*Connection) bool {
	var current *Connection
	var expectedConn *Connection
	if len(expected) > 0 {
		expectedConn = expected[0]
	}
	connections.lock.Lock()
	if connections.conn == nil {
		connections.lock.Unlock()
		return true
	}
	current = connections.conn[key]
	if current == nil {
		delete(connections.conn, key)
		connections.lock.Unlock()
		return true
	}
	if expectedConn != nil && current != expectedConn {
		connections.lock.Unlock()
		return false
	}
	delete(connections.conn, key)
	connections.lock.Unlock()

	if err := current.Close(); err != nil {
		return false
	}
	return true
}

// Connection 对 gorilla websocket 连接做了一层并发和关闭保护。
type Connection struct {
	*websocket.Conn
	writeLock sync.Mutex
	closeOnce sync.Once
	closeErr  error
	writeWait time.Duration
}

// NewConnection 包装原始 websocket 连接
func NewConnection(conn *websocket.Conn) *Connection {
	return &Connection{
		Conn: conn,
	}
}

// withWriteLock 串行化写入，避免 gorilla websocket 并发写 panic。
func (connection *Connection) withWriteLock(fun func() error) error {
	if connection == nil || connection.Conn == nil {
		return errNilConnection
	}

	connection.writeLock.Lock()
	defer connection.writeLock.Unlock()
	return fun()
}

// WriteMessage 在写入业务消息前加写锁和写超时控制。
func (connection *Connection) WriteMessage(messageType int, data []byte) error {
	return connection.withWriteLock(func() error {
		if connection.writeWait > 0 {
			_ = connection.Conn.SetWriteDeadline(time.Now().Add(connection.writeWait))
		}
		return connection.Conn.WriteMessage(messageType, data)
	})
}

// WriteJSON 在写入 JSON 前加写锁和写超时控制。
func (connection *Connection) WriteJSON(v interface{}) error {
	return connection.withWriteLock(func() error {
		if connection.writeWait > 0 {
			_ = connection.Conn.SetWriteDeadline(time.Now().Add(connection.writeWait))
		}
		return connection.Conn.WriteJSON(v)
	})
}

// WritePreparedMessage 在写入 PreparedMessage 前加写锁和写超时控制。
func (connection *Connection) WritePreparedMessage(pm *websocket.PreparedMessage) error {
	return connection.withWriteLock(func() error {
		if connection.writeWait > 0 {
			_ = connection.Conn.SetWriteDeadline(time.Now().Add(connection.writeWait))
		}
		return connection.Conn.WritePreparedMessage(pm)
	})
}

// WriteControl 串行化控制帧写入。
func (connection *Connection) WriteControl(messageType int, data []byte, deadline time.Time) error {
	return connection.withWriteLock(func() error {
		return connection.Conn.WriteControl(messageType, data, deadline)
	})
}

// Close 保证连接只会被真正关闭一次。
func (connection *Connection) Close() error {
	if connection == nil || connection.Conn == nil {
		return nil
	}

	connection.closeOnce.Do(func() {
		connection.closeErr = connection.Conn.Close()
	})
	return connection.closeErr
}

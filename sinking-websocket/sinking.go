package sinking_websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewWebSocketConnections 单例
func NewWebSocketConnections() *WebSocketConnections {
	wsConn := &WebSocketConnections{}
	if wsConn.conn == nil {
		wsConn.conn = make(map[string]*Conn)
	}
	return wsConn
}

// WebSocketConnections ws连接用户
type WebSocketConnections struct {
	conn map[string]*Conn
	lock sync.Mutex
}

// Get 获取长连接对象
func (connections *WebSocketConnections) Get(key string) *Conn {
	connections.lock.Lock()
	defer connections.lock.Unlock()
	return connections.conn[key]
}

// GetAll 获取所有长连接对象
func (connections *WebSocketConnections) GetAll() map[string]*Conn {
	connections.lock.Lock()
	defer connections.lock.Unlock()
	conn := make(map[string]*Conn)
	for k, v := range connections.conn {
		conn[k] = v
	}
	return conn
}

// Set 设置长连接对象
func (connections *WebSocketConnections) Set(key string, conn *Conn) {
	connections.lock.Lock()
	defer connections.lock.Unlock()
	connections.conn[key] = conn
}

// Delete 删除长连接对象
func (connections *WebSocketConnections) Delete(key string) bool {
	connections.lock.Lock()
	defer connections.lock.Unlock()
	if connections.conn[key] != nil {
		err := connections.conn[key].Close()
		if err != nil {
			return false
		}
	}
	connections.conn[key] = nil
	return true
}

// Conn conn 包装
type Conn struct {
	*websocket.Conn
}

// NewWebSocket 单例
func NewWebSocket() *WebSocket {
	wsConn := &WebSocket{}
	return wsConn
}

// WebSocket 执行
type WebSocket struct {
	OnError   func(err error)
	OnConnect func(ws *Conn)
	OnClose   func(err error)
	OnMessage func(ws *Conn, messageType int, data []byte)
}

func (handle *WebSocket) SetErrorHandle(fun func(err error)) *WebSocket {
	handle.OnError = fun
	return handle
}

func (handle *WebSocket) SetConnectHandle(fun func(ws *Conn)) *WebSocket {
	handle.OnConnect = fun
	return handle
}

func (handle *WebSocket) SetCloseHandle(fun func(err error)) *WebSocket {
	handle.OnClose = fun
	return handle
}

func (handle *WebSocket) SetOnMessageHandle(fun func(ws *Conn, messageType int, data []byte)) *WebSocket {
	handle.OnMessage = fun
	return handle
}

type Error struct {
	ErrCode int
	ErrMsg  string
}

func (err *Error) Error() string {
	return err.ErrMsg
}

func (handle *WebSocket) Listen(writer http.ResponseWriter, request *http.Request) {
	defer func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Connection") != "Upgrade" {
			if handle.OnError != nil {
				handle.OnError(&Error{ErrCode: 500, ErrMsg: "Connection Close"})
			}
			return
		}
		ws, err := upGrader.Upgrade(writer, request, nil)
		if err != nil {
			if handle.OnError != nil {
				handle.OnError(err)
			}
			return
		}
		conn := &Conn{ws}
		defer func(con *Conn) {
			_ = con.Close()
		}(conn) //返回前关闭
		if handle.OnConnect != nil {
			handle.OnConnect(conn)
		}
		for {
			//读取ws中的数据
			mt, message, err := conn.ReadMessage()
			if err != nil {
				if handle.OnClose != nil {
					handle.OnClose(err)
				}
				return
			} else {
				if handle.OnMessage != nil {
					handle.OnMessage(conn, mt, message)
				}
			}
		}
	}(writer, request)
}

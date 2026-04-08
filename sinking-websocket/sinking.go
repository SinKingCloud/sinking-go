package sinking_websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultHandshakeTimeout = 10 * time.Second
	defaultWriteWait        = 10 * time.Second
	defaultPongWait         = 60 * time.Second
)

var (
	defaultCheckOrigin = func(r *http.Request) bool {
		return true
	}
)

// NewWebSocket 创建 websocket 处理器
func NewWebSocket() *WebSocket {
	return &WebSocket{}
}

// WebSocket 封装 websocket 连接的升级、心跳和回调处理。
type WebSocket struct {
	Id                string
	HandshakeTimeout  time.Duration
	ReadBufferSize    int
	WriteBufferSize   int
	ReadLimit         int64
	WriteWait         time.Duration
	PongWait          time.Duration
	PingPeriod        time.Duration
	EnableCompression bool
	DisableHeartbeat  bool
	CheckOrigin       func(r *http.Request) bool
	OnError           func(id string, err error)
	OnConnect         func(id string, ws *Connection)
	OnCloseFrame      func(id string, ws *Connection, code int, text string) error
	OnClose           func(id string, ws *Connection, err error)
	OnMessage         func(id string, ws *Connection, messageType int, data []byte)
	OnPing            func(id string, ws *Connection, appData string)
	OnPong            func(id string, ws *Connection, appData string)
}

// SetId 设置当前 websocket 处理器关联的业务标识。
func (handle *WebSocket) SetId(id string) *WebSocket {
	handle.Id = id
	return handle
}

// SetErrorHandle 设置握手失败等错误回调。
func (handle *WebSocket) SetErrorHandle(fun func(id string, err error)) *WebSocket {
	handle.OnError = fun
	return handle
}

// SetConnectHandle 设置连接建立成功后的回调。
func (handle *WebSocket) SetConnectHandle(fun func(id string, ws *Connection)) *WebSocket {
	handle.OnConnect = fun
	return handle
}

// SetCloseHandler 设置收到 close frame 时的协议层回调。
func (handle *WebSocket) SetCloseHandler(fun func(id string, ws *Connection, code int, text string) error) *WebSocket {
	handle.OnCloseFrame = fun
	return handle
}

// SetCloseFrameHandle 是 SetCloseHandler 的语义化别名。
func (handle *WebSocket) SetCloseFrameHandle(fun func(id string, ws *Connection, code int, text string) error) *WebSocket {
	handle.OnCloseFrame = fun
	return handle
}

// SetCloseHandle 设置连接读循环退出时的业务层关闭回调。
func (handle *WebSocket) SetCloseHandle(fun func(id string, ws *Connection, err error)) *WebSocket {
	handle.OnClose = fun
	return handle
}

// SetOnMessageHandle 设置收到业务消息后的回调。
func (handle *WebSocket) SetOnMessageHandle(fun func(id string, ws *Connection, messageType int, data []byte)) *WebSocket {
	handle.OnMessage = fun
	return handle
}

// SetPingHandle 设置收到 ping 控制帧后的回调。
func (handle *WebSocket) SetPingHandle(fun func(id string, ws *Connection, appData string)) *WebSocket {
	handle.OnPing = fun
	return handle
}

// SetPongHandle 设置收到 pong 控制帧后的回调。
func (handle *WebSocket) SetPongHandle(fun func(id string, ws *Connection, appData string)) *WebSocket {
	handle.OnPong = fun
	return handle
}

// Error 表示自定义业务错误。
type Error struct {
	ErrCode int
	ErrMsg  string
}

// Error 返回错误信息文本。
func (err *Error) Error() string {
	return err.ErrMsg
}

// Listen 完成协议升级并持续处理连接生命周期。
func (handle *WebSocket) Listen(writer http.ResponseWriter, request *http.Request, responseHeader http.Header) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout:  handle.handshakeTimeout(),
		ReadBufferSize:    handle.ReadBufferSize,
		WriteBufferSize:   handle.WriteBufferSize,
		CheckOrigin:       handle.checkOrigin(),
		EnableCompression: handle.EnableCompression,
	}
	ws, err := upgrader.Upgrade(writer, request, responseHeader)
	if err != nil {
		if handle.OnError != nil {
			handle.OnError(handle.Id, err)
		}
		return
	}
	conn := NewConnection(ws)
	done := make(chan struct{})
	defer func() {
		close(done)
		_ = conn.Close()
	}()
	handle.configureConn(conn)
	if handle.OnConnect != nil {
		handle.OnConnect(handle.Id, conn)
	}
	if pingPeriod := handle.pingPeriod(); pingPeriod > 0 {
		go handle.keepAlive(conn, done, pingPeriod)
	}
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if handle.OnClose != nil {
				handle.OnClose(handle.Id, conn, err)
			}
			return
		}
		if handle.OnMessage != nil {
			handle.OnMessage(handle.Id, conn, messageType, data)
		}
	}
}

// configureConn 配置读写限制、心跳处理和关闭回调。
func (handle *WebSocket) configureConn(conn *Connection) {
	defaultCloseHandler := conn.CloseHandler()
	conn.writeWait = handle.writeWait()
	if handle.ReadLimit > 0 {
		conn.SetReadLimit(handle.ReadLimit)
	}
	if !handle.DisableHeartbeat {
		_ = conn.SetReadDeadline(time.Now().Add(handle.pongWait()))
	}
	conn.SetCloseHandler(func(code int, text string) error {
		if handle.OnCloseFrame != nil {
			if err := handle.OnCloseFrame(handle.Id, conn, code, text); err != nil {
				return err
			}
		}
		if defaultCloseHandler != nil {
			return defaultCloseHandler(code, text)
		}
		return nil
	})
	conn.SetPingHandler(func(appData string) error {
		if !handle.DisableHeartbeat {
			_ = conn.SetReadDeadline(time.Now().Add(handle.pongWait()))
		}
		if handle.OnPing != nil {
			handle.OnPing(handle.Id, conn, appData)
		}
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(handle.writeWait()))
	})
	conn.SetPongHandler(func(appData string) error {
		if !handle.DisableHeartbeat {
			_ = conn.SetReadDeadline(time.Now().Add(handle.pongWait()))
		}
		if handle.OnPong != nil {
			handle.OnPong(handle.Id, conn, appData)
		}
		return nil
	})
}

// keepAlive 按固定周期主动发送 ping 保持连接活性。
func (handle *WebSocket) keepAlive(conn *Connection, done <-chan struct{}, pingPeriod time.Duration) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(handle.writeWait())); err != nil {
				_ = conn.Close()
				return
			}
		}
	}
}

// checkOrigin 返回本次升级请求使用的 Origin 校验函数。
func (handle *WebSocket) checkOrigin() func(r *http.Request) bool {
	if handle.CheckOrigin != nil {
		return handle.CheckOrigin
	}
	return defaultCheckOrigin
}

// handshakeTimeout 返回 websocket 握手超时时间。
func (handle *WebSocket) handshakeTimeout() time.Duration {
	if handle.HandshakeTimeout > 0 {
		return handle.HandshakeTimeout
	}
	return defaultHandshakeTimeout
}

// writeWait 返回单次写操作允许的最长时间。
func (handle *WebSocket) writeWait() time.Duration {
	if handle.WriteWait > 0 {
		return handle.WriteWait
	}
	return defaultWriteWait
}

// pongWait 返回等待对端 pong 的最长时间。
func (handle *WebSocket) pongWait() time.Duration {
	if handle.PongWait > 0 {
		return handle.PongWait
	}
	return defaultPongWait
}

// pingPeriod 返回主动发送 ping 的时间间隔。
func (handle *WebSocket) pingPeriod() time.Duration {
	if handle.DisableHeartbeat {
		return 0
	}
	if handle.PingPeriod > 0 {
		pongWait := handle.pongWait()
		if pongWait > 0 && handle.PingPeriod >= pongWait {
			return pongWait * 9 / 10
		}
		return handle.PingPeriod
	}
	return handle.pongWait() * 9 / 10
}

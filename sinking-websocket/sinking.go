package sinking_websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Websocket struct {
	OnError   func(err error)
	OnConnect func(ws *websocket.Conn)
	OnClose   func(err error)
	OnMessage func(ws *websocket.Conn, messageType int, data []byte)
}

func (handle *Websocket) SetErrorHandle(fun func(err error)) *Websocket {
	handle.OnError = fun
	return handle
}

func (handle *Websocket) SetConnectHandle(fun func(ws *websocket.Conn)) *Websocket {
	handle.OnConnect = fun
	return handle
}

func (handle *Websocket) SetCloseHandle(fun func(err error)) *Websocket {
	handle.OnClose = fun
	return handle
}

func (handle *Websocket) SetOnMessageHandle(fun func(ws *websocket.Conn, messageType int, data []byte)) *Websocket {
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

func (handle *Websocket) Listen(writer http.ResponseWriter, request *http.Request) {
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
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				if handle.OnClose != nil {
					handle.OnClose(err)
				}
				return
			}
		}(ws) //返回前关闭
		if handle.OnConnect != nil {
			handle.OnConnect(ws)
		}
		for {
			//读取ws中的数据
			mt, message, err := ws.ReadMessage()
			if err != nil {
				if handle.OnClose != nil {
					handle.OnClose(err)
				}
				return
			} else {
				if handle.OnMessage != nil {
					handle.OnMessage(ws, mt, message)
				}
			}
		}
	}(writer, request)
}

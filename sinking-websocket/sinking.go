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

type WebsocketHandle struct {
	ErrorHandle     func(err error)
	OnMessageHandle func(ws *websocket.Conn, messageType int, data []byte)
}

func (handle *WebsocketHandle) SetErrorHandle(fun func(err error)) *WebsocketHandle {
	handle.ErrorHandle = fun
	return handle
}

func (handle *WebsocketHandle) SetOnMessageHandle(fun func(ws *websocket.Conn, messageType int, data []byte)) *WebsocketHandle {
	handle.OnMessageHandle = fun
	return handle
}

type Error struct {
	ErrCode int
	ErrMsg  string
}

func (err *Error) Error() string {
	return err.ErrMsg
}

func (handle *WebsocketHandle) Listen(writer http.ResponseWriter, request *http.Request) {
	defer func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Connection") != "Upgrade" {
			if handle.ErrorHandle != nil {
				handle.ErrorHandle(&Error{ErrCode: 500, ErrMsg: "Connection Close"})
			}
			return
		}
		ws, err := upGrader.Upgrade(writer, request, nil)
		if err != nil {
			if handle.ErrorHandle != nil {
				handle.ErrorHandle(err)
			}
			return
		}
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				if handle.ErrorHandle != nil {
					handle.ErrorHandle(err)
				}
				return
			}
		}(ws) //返回前关闭
		for {
			//读取ws中的数据
			mt, message, err := ws.ReadMessage()
			if err != nil {
				if handle.ErrorHandle != nil {
					handle.ErrorHandle(err)
				}
				return
			} else {
				if handle.OnMessageHandle != nil {
					handle.OnMessageHandle(ws, mt, message)
				}
			}
		}
	}(writer, request)
}

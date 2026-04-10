package sinking_websocket

import "github.com/gorilla/websocket"

const (
	TextMessage   = websocket.TextMessage
	BinaryMessage = websocket.BinaryMessage
	CloseMessage  = websocket.CloseMessage
	PingMessage   = websocket.PingMessage
	PongMessage   = websocket.PongMessage
)

type PreparedMessage = websocket.PreparedMessage

type Message struct {
	Type    int
	Payload []byte
}

type BroadcastResult struct {
	Queued  int
	Dropped int
	Closed  int
}

func PrepareMessage(messageType int, payload []byte) (*PreparedMessage, error) {
	return websocket.NewPreparedMessage(messageType, payload)
}

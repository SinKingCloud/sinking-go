package sinking_websocket

import (
	"errors"
	"fmt"
)

var (
	errNilConnection      = errors.New("websocket connection is nil")
	errNilPreparedMessage = errors.New("websocket prepared message is nil")

	ErrConnectionClosed = errors.New("websocket connection is closed")
	ErrSendQueueFull    = errors.New("websocket send queue is full")
)

type panicError struct {
	phase string
	value interface{}
	stack []byte
}

func (err *panicError) Error() string {
	return fmt.Sprintf("websocket %s panic: %v", err.phase, err.value)
}

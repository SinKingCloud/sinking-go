package sinking_websocket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultHandshakeTimeout = 10 * time.Second
	defaultWriteBufferSize  = 4096
	defaultWriteTimeout     = 10 * time.Second
	defaultPongTimeout      = 60 * time.Second
	defaultWriteQueueSize   = 64
)

var defaultOriginValidator = func(r *http.Request) bool {
	return true
}

type ConnectionIDResolver func(request *http.Request) (string, error)

type ConnectHandler func(connection *Connection) error

type MessageHandler func(connection *Connection, message Message) error

type DisconnectHandler func(connection *Connection, err error)

type UpgradeErrorHandler func(request *http.Request, err error)

type CloseFrameHandler func(connection *Connection, code int, text string) error

type ControlHandler func(connection *Connection, payload string)

type ServerOption func(config *serverConfig)

type serverConfig struct {
	handshakeTimeout       time.Duration
	readBufferSize         int
	writeBufferSize        int
	writeQueueSize         int
	readLimit              int64
	writeTimeout           time.Duration
	pongTimeout            time.Duration
	pingInterval           time.Duration
	heartbeatSweepInterval time.Duration
	enableCompression      bool
	disableHeartbeat       bool
	originValidator        func(r *http.Request) bool
	writeBufferPool        websocket.BufferPool
	writeDispatcherShards  int
	idResolver             ConnectionIDResolver
	upgradeError           UpgradeErrorHandler
	connectHandler         ConnectHandler
	closeFrameHandler      CloseFrameHandler
	disconnectHandler      DisconnectHandler
	messageHandler         MessageHandler
	pingHandler            ControlHandler
	pongHandler            ControlHandler
}

func defaultServerConfig() serverConfig {
	return serverConfig{
		handshakeTimeout: defaultHandshakeTimeout,
		writeBufferSize:  defaultWriteBufferSize,
		writeQueueSize:   defaultWriteQueueSize,
		writeTimeout:     defaultWriteTimeout,
		pongTimeout:      defaultPongTimeout,
		originValidator:  defaultOriginValidator,
	}
}

func (config serverConfig) resolvedPingInterval() time.Duration {
	if config.disableHeartbeat {
		return 0
	}
	if config.pingInterval > 0 {
		if config.pongTimeout > 0 && config.pingInterval >= config.pongTimeout {
			return config.pongTimeout * 9 / 10
		}
		return config.pingInterval
	}
	if config.pongTimeout <= 0 {
		return 0
	}
	return config.pongTimeout * 9 / 10
}

func WithConnectionID(id string) ServerOption {
	return WithConnectionIDResolver(func(request *http.Request) (string, error) {
		return id, nil
	})
}

func WithConnectionIDResolver(resolver ConnectionIDResolver) ServerOption {
	return func(config *serverConfig) {
		config.idResolver = resolver
	}
}

func WithHandshakeTimeout(timeout time.Duration) ServerOption {
	return func(config *serverConfig) {
		if timeout > 0 {
			config.handshakeTimeout = timeout
		}
	}
}

func WithReadBufferSize(size int) ServerOption {
	return func(config *serverConfig) {
		if size > 0 {
			config.readBufferSize = size
		}
	}
}

func WithWriteBufferSize(size int) ServerOption {
	return func(config *serverConfig) {
		if size > 0 {
			config.writeBufferSize = size
		}
	}
}

func WithWriteBufferPool(pool websocket.BufferPool) ServerOption {
	return func(config *serverConfig) {
		config.writeBufferPool = pool
	}
}

func WithWriteQueueSize(size int) ServerOption {
	return func(config *serverConfig) {
		if size > 0 {
			config.writeQueueSize = size
		}
	}
}

func WithWriteDispatcherShards(shards int) ServerOption {
	return func(config *serverConfig) {
		if shards > 0 {
			config.writeDispatcherShards = shards
		}
	}
}

func WithReadLimit(limit int64) ServerOption {
	return func(config *serverConfig) {
		if limit > 0 {
			config.readLimit = limit
		}
	}
}

func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(config *serverConfig) {
		if timeout > 0 {
			config.writeTimeout = timeout
		}
	}
}

func WithPongTimeout(timeout time.Duration) ServerOption {
	return func(config *serverConfig) {
		if timeout > 0 {
			config.pongTimeout = timeout
		}
	}
}

func WithPingInterval(interval time.Duration) ServerOption {
	return func(config *serverConfig) {
		if interval > 0 {
			config.pingInterval = interval
		}
	}
}

func WithHeartbeatSweepInterval(interval time.Duration) ServerOption {
	return func(config *serverConfig) {
		if interval > 0 {
			config.heartbeatSweepInterval = interval
		}
	}
}

func resolvedWriteQueueSize(size int) int {
	if size > 0 {
		return size
	}
	return defaultWriteQueueSize
}

func WithCompression(enabled bool) ServerOption {
	return func(config *serverConfig) {
		config.enableCompression = enabled
	}
}

func WithoutHeartbeat() ServerOption {
	return func(config *serverConfig) {
		config.disableHeartbeat = true
	}
}

func WithOriginValidator(validator func(r *http.Request) bool) ServerOption {
	return func(config *serverConfig) {
		if validator != nil {
			config.originValidator = validator
		}
	}
}

func WithUpgradeErrorHandler(handler UpgradeErrorHandler) ServerOption {
	return func(config *serverConfig) {
		config.upgradeError = handler
	}
}

func WithConnectHandler(handler ConnectHandler) ServerOption {
	return func(config *serverConfig) {
		config.connectHandler = handler
	}
}

func WithCloseFrameHandler(handler CloseFrameHandler) ServerOption {
	return func(config *serverConfig) {
		config.closeFrameHandler = handler
	}
}

func WithDisconnectHandler(handler DisconnectHandler) ServerOption {
	return func(config *serverConfig) {
		config.disconnectHandler = handler
	}
}

func WithMessageHandler(handler MessageHandler) ServerOption {
	return func(config *serverConfig) {
		config.messageHandler = handler
	}
}

func WithPingHandler(handler ControlHandler) ServerOption {
	return func(config *serverConfig) {
		config.pingHandler = handler
	}
}

func WithPongHandler(handler ControlHandler) ServerOption {
	return func(config *serverConfig) {
		config.pongHandler = handler
	}
}

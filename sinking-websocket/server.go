package sinking_websocket

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	config     serverConfig
	upgrader   websocket.Upgrader
	dispatcher *writeDispatcher
}

func NewServer(options ...ServerOption) *Server {
	config := defaultServerConfig()
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}

	return &Server{
		config: config,
		upgrader: websocket.Upgrader{
			HandshakeTimeout:  config.handshakeTimeout,
			ReadBufferSize:    config.readBufferSize,
			WriteBufferSize:   resolvedWriteBufferSize(config.writeBufferSize),
			WriteBufferPool:   resolvedWriteBufferPool(config),
			CheckOrigin:       config.originValidator,
			EnableCompression: config.enableCompression,
		},
		dispatcher: newWriteDispatcher(writeDispatcherConfig{
			shardCount:             config.writeDispatcherShards,
			pingInterval:           config.resolvedPingInterval(),
			heartbeatSweepInterval: config.heartbeatSweepInterval,
		}),
	}
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_ = server.Handle(writer, request, nil)
}

func (server *Server) Handle(writer http.ResponseWriter, request *http.Request, responseHeader http.Header) error {
	connectionID, err := server.resolveConnectionID(request)
	if err != nil {
		server.reportUpgradeError(request, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return err
	}

	raw, err := server.upgrader.Upgrade(writer, request, responseHeader)
	if err != nil {
		server.reportUpgradeError(request, err)
		return err
	}

	connection := newConnection(connectionID, request, raw, connectionConfig{
		writeTimeout:   server.config.writeTimeout,
		pingInterval:   server.config.resolvedPingInterval(),
		writeQueueSize: server.config.writeQueueSize,
		dispatcher:     server.dispatcher,
	})

	server.configureConnection(connection)
	return server.serveConnection(connection)
}

func (server *Server) configureConnection(connection *Connection) {
	defaultCloseHandler := connection.raw.CloseHandler()

	if server.config.readLimit > 0 {
		connection.raw.SetReadLimit(server.config.readLimit)
	}

	if !server.config.disableHeartbeat && server.config.pongTimeout > 0 {
		_ = connection.raw.SetReadDeadline(time.Now().Add(server.config.pongTimeout))
	}

	connection.raw.SetCloseHandler(func(code int, text string) error {
		if err := server.callCloseFrameHandler(connection, code, text); err != nil {
			return err
		}
		if defaultCloseHandler != nil {
			return defaultCloseHandler(code, text)
		}
		return nil
	})

	connection.raw.SetPingHandler(func(payload string) error {
		if !server.config.disableHeartbeat && server.config.pongTimeout > 0 {
			_ = connection.raw.SetReadDeadline(time.Now().Add(server.config.pongTimeout))
		}
		if err := server.callPingHandler(connection, payload); err != nil {
			return err
		}
		return connection.writeControl(PongMessage, []byte(payload))
	})

	connection.raw.SetPongHandler(func(payload string) error {
		if !server.config.disableHeartbeat && server.config.pongTimeout > 0 {
			_ = connection.raw.SetReadDeadline(time.Now().Add(server.config.pongTimeout))
		}
		return server.callPongHandler(connection, payload)
	})
}

func (server *Server) serveConnection(connection *Connection) (runErr error) {
	defer func() {
		_ = connection.Close()
		server.callDisconnectHandler(connection, normalizeDisconnectError(runErr))
	}()

	if err := server.callConnectHandler(connection); err != nil {
		runErr = err
		return runErr
	}

	for {
		messageType, payload, err := connection.raw.ReadMessage()
		if err != nil {
			runErr = err
			return runErr
		}

		if err := server.callMessageHandler(connection, Message{
			Type:    messageType,
			Payload: payload,
		}); err != nil {
			runErr = err
			return runErr
		}
	}
}

func (server *Server) resolveConnectionID(request *http.Request) (string, error) {
	if server.config.idResolver == nil {
		return "", nil
	}
	return server.config.idResolver(request)
}

func (server *Server) reportUpgradeError(request *http.Request, err error) {
	if server.config.upgradeError != nil {
		server.config.upgradeError(request, err)
	}
}

func (server *Server) callConnectHandler(connection *Connection) error {
	if server.config.connectHandler == nil {
		return nil
	}
	return recoverCall("connect handler", func() error {
		return server.config.connectHandler(connection)
	})
}

func (server *Server) callMessageHandler(connection *Connection, message Message) error {
	if server.config.messageHandler == nil {
		return nil
	}
	return recoverCall("message handler", func() error {
		return server.config.messageHandler(connection, message)
	})
}

func (server *Server) callDisconnectHandler(connection *Connection, err error) {
	if server.config.disconnectHandler == nil {
		return
	}

	_ = recoverCall("disconnect handler", func() error {
		server.config.disconnectHandler(connection, err)
		return nil
	})
}

func (server *Server) callCloseFrameHandler(connection *Connection, code int, text string) error {
	if server.config.closeFrameHandler == nil {
		return nil
	}
	return recoverCall("close frame handler", func() error {
		return server.config.closeFrameHandler(connection, code, text)
	})
}

func (server *Server) callPingHandler(connection *Connection, payload string) error {
	if server.config.pingHandler == nil {
		return nil
	}
	return recoverCall("ping handler", func() error {
		server.config.pingHandler(connection, payload)
		return nil
	})
}

func (server *Server) callPongHandler(connection *Connection, payload string) error {
	if server.config.pongHandler == nil {
		return nil
	}
	return recoverCall("pong handler", func() error {
		server.config.pongHandler(connection, payload)
		return nil
	})
}

func recoverCall(phase string, call func() error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = &panicError{
				phase: phase,
				value: recovered,
				stack: debug.Stack(),
			}
		}
	}()

	return call()
}

func normalizeDisconnectError(err error) error {
	if err == nil {
		return nil
	}
	if err == websocket.ErrCloseSent {
		return nil
	}
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		return nil
	}
	return err
}

func resolvedWriteBufferPool(config serverConfig) websocket.BufferPool {
	if config.writeBufferPool != nil {
		return config.writeBufferPool
	}
	return sharedWriteBufferPool(config.writeBufferSize)
}

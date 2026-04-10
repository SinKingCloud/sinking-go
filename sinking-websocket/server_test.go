package sinking_websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestRegistryDeleteIfMatch(t *testing.T) {
	registry := NewRegistry()
	first := newConnection("user-1", nil, nil, connectionConfig{})
	second := newConnection("user-1", nil, nil, connectionConfig{})
	defer func() {
		_ = first.Close()
		_ = second.Close()
	}()

	registry.Store("user-1", first)
	if registry.Len() != 1 {
		t.Fatalf("expected registry size 1, got %d", registry.Len())
	}

	if registry.DeleteIfMatch("user-1", second) {
		t.Fatalf("expected delete with stale connection to fail")
	}

	current, ok := registry.Load("user-1")
	if !ok || current != first {
		t.Fatalf("expected first connection to stay registered")
	}

	registry.Store("user-1", second)
	if !first.Closed() {
		t.Fatalf("expected replaced connection to be closed")
	}

	if registry.DeleteIfMatch("user-1", first) {
		t.Fatalf("expected delete with replaced connection to fail")
	}

	if !registry.DeleteIfMatch("user-1", second) {
		t.Fatalf("expected delete with current connection to succeed")
	}

	if registry.Len() != 0 {
		t.Fatalf("expected registry size 0, got %d", registry.Len())
	}
}

func TestRegistryRangeAllowsDeletion(t *testing.T) {
	registry := NewRegistry()
	first := newConnection("user-1", nil, nil, connectionConfig{})
	second := newConnection("user-2", nil, nil, connectionConfig{})
	defer func() {
		_ = first.Close()
		_ = second.Close()
	}()

	registry.Store("user-1", first)
	registry.Store("user-2", second)

	registry.Range(func(id string, connection *Connection) bool {
		registry.DeleteIfMatch(id, connection)
		return true
	})

	if registry.Len() != 0 {
		t.Fatalf("expected registry to be empty after range delete, got %d", registry.Len())
	}
}

func TestServerEchoLifecycle(t *testing.T) {
	registry := NewRegistry()
	disconnects := make(chan error, 1)

	server := httptest.NewServer(NewServer(
		WithConnectionIDResolver(func(request *http.Request) (string, error) {
			return request.URL.Query().Get("id"), nil
		}),
		WithConnectHandler(func(connection *Connection) error {
			registry.Store(connection.ID(), connection)
			return nil
		}),
		WithMessageHandler(func(connection *Connection, message Message) error {
			return connection.Send(message.Type, message.Payload)
		}),
		WithDisconnectHandler(func(connection *Connection, err error) {
			registry.DeleteIfMatch(connection.ID(), connection)
			disconnects <- err
		}),
	))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(websocketURL(server.URL+"/?id=echo-user"), nil)
	if err != nil {
		t.Fatalf("failed to dial websocket server: %v", err)
	}
	defer client.Close()

	if err := client.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
		t.Fatalf("failed to send websocket message: %v", err)
	}

	messageType, payload, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read websocket message: %v", err)
	}

	if messageType != websocket.TextMessage {
		t.Fatalf("expected text message, got %d", messageType)
	}

	if string(payload) != "hello" {
		t.Fatalf("expected echoed payload hello, got %s", string(payload))
	}

	connection, ok := registry.Load("echo-user")
	if !ok || connection == nil {
		t.Fatalf("expected websocket connection to be registered")
	}

	if err := client.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second),
	); err != nil {
		t.Fatalf("failed to send close frame: %v", err)
	}

	select {
	case disconnectErr := <-disconnects:
		if disconnectErr != nil {
			t.Fatalf("expected clean disconnect, got %v", disconnectErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for disconnect callback")
	}

	waitFor(t, func() bool {
		return registry.Len() == 0
	}, "registry should be empty after disconnect")
}

func TestServerRecoversMessageHandlerPanic(t *testing.T) {
	disconnects := make(chan error, 1)

	server := httptest.NewServer(NewServer(
		WithMessageHandler(func(connection *Connection, message Message) error {
			panic("boom")
		}),
		WithDisconnectHandler(func(connection *Connection, err error) {
			disconnects <- err
		}),
	))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(websocketURL(server.URL), nil)
	if err != nil {
		t.Fatalf("failed to dial websocket server: %v", err)
	}
	defer client.Close()

	if err := client.WriteMessage(websocket.TextMessage, []byte("trigger")); err != nil {
		t.Fatalf("failed to send websocket message: %v", err)
	}

	select {
	case disconnectErr := <-disconnects:
		if disconnectErr == nil {
			t.Fatalf("expected panic error, got nil")
		}
		if !strings.Contains(disconnectErr.Error(), "message handler panic") {
			t.Fatalf("expected panic error, got %v", disconnectErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for disconnect callback")
	}
}

func websocketURL(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}

func waitFor(t *testing.T, condition func() bool, message string) {
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal(message)
}

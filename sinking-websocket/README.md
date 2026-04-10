# Sinking-WebSocket

Sinking-WebSocket 是一个围绕 `gorilla/websocket` 封装的轻量级服务端模块，目标是把连接生命周期、单写协程、连接注册表和回调边界收敛成一套更稳妥也更简洁的 API。

## 设计目标

- 单连接固定一个写协程，避免并发写 panic
- 关闭路径只依赖 `done` 信号，不关闭发送队列，避免 `send on closed channel`
- 连接注册表支持分片存储、条件删除和安全遍历
- 对外只暴露高频需要的能力：建连、收消息、发消息、注册连接、广播连接

## 快速开始

```go
package main

import (
	"log"
	"net/http"

	"github.com/SinKingCloud/sinking-go/sinking-websocket"
)

func main() {
	registry := sinking_websocket.NewRegistry()

	handler := sinking_websocket.NewServer(
		sinking_websocket.WithConnectionIDResolver(func(request *http.Request) (string, error) {
			return request.URL.Query().Get("id"), nil
		}),
		sinking_websocket.WithConnectHandler(func(connection *sinking_websocket.Connection) error {
			registry.Store(connection.ID(), connection)
			return connection.SendJSON(map[string]string{
				"type": "welcome",
				"id":   connection.ID(),
			})
		}),
		sinking_websocket.WithMessageHandler(func(connection *sinking_websocket.Connection, message sinking_websocket.Message) error {
			return connection.Send(message.Type, message.Payload)
		}),
		sinking_websocket.WithDisconnectHandler(func(connection *sinking_websocket.Connection, err error) {
			registry.DeleteIfMatch(connection.ID(), connection)
			log.Printf("disconnect: id=%s err=%v", connection.ID(), err)
		}),
	)

	http.Handle("/ws", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## 广播消息

```go
func broadcast(registry *sinking_websocket.Registry, payload []byte) {
	result, err := registry.Broadcast(sinking_websocket.TextMessage, payload)
	if err != nil {
		return
	}
	_ = result
}
```

## 核心 API

- `NewServer(options ...ServerOption)`：创建 WebSocket 服务端处理器
- `(*Server).ServeHTTP` / `(*Server).Handle`：升级并托管连接生命周期
- `(*Connection).Send` / `TrySend` / `SendJSON`：异步串行发送消息
- `NewRegistry()`：创建连接注册表
- `(*Registry).Load` / `Store` / `DeleteIfMatch` / `Range`：管理在线连接
- `(*Registry).Broadcast` / `BroadcastPrepared` / `BroadcastJSON`：高频广播快路径
- `PrepareMessage(...)`：预构建消息，适合高频广播

更详细的示例见 [doc.md](./doc.md)。

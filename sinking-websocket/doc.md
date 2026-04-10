# Sinking-WebSocket 使用文档

## 整体结构

重构后的模块只保留四个核心概念：

- `Server`：负责升级请求、运行读循环、托管心跳和回调
- `Connection`：代表一条已建立的 WebSocket 连接，所有业务消息都通过内部单写协程串行发送
- `Registry`：分片化连接注册表，负责在线连接的存取、遍历和条件删除
- `Message`：读循环交给业务层的消息载体

## 推荐用法

```go
registry := sinking_websocket.NewRegistry()

handler := sinking_websocket.NewServer(
	sinking_websocket.WithConnectionIDResolver(func(request *http.Request) (string, error) {
		return request.URL.Query().Get("user_id"), nil
	}),
	sinking_websocket.WithConnectHandler(func(connection *sinking_websocket.Connection) error {
		registry.Store(connection.ID(), connection)
		return nil
	}),
	sinking_websocket.WithMessageHandler(func(connection *sinking_websocket.Connection, message sinking_websocket.Message) error {
		return connection.Send(message.Type, message.Payload)
	}),
	sinking_websocket.WithDisconnectHandler(func(connection *sinking_websocket.Connection, err error) {
		registry.DeleteIfMatch(connection.ID(), connection)
	}),
)

http.Handle("/ws", handler)
```

## Server 选项

### 生命周期

- `WithConnectionID(...)`：为当前 handler 指定固定连接 ID
- `WithConnectionIDResolver(...)`：根据请求动态生成连接 ID
- `WithConnectHandler(...)`：建连成功后执行
- `WithMessageHandler(...)`：收到业务消息后执行
- `WithDisconnectHandler(...)`：连接结束后执行
- `WithUpgradeErrorHandler(...)`：升级失败或连接建立前置校验失败时执行

### 网络参数

- `WithHandshakeTimeout(...)`
- `WithReadBufferSize(...)`
- `WithWriteBufferSize(...)`
- `WithWriteBufferPool(...)`
- `WithReadLimit(...)`
- `WithWriteTimeout(...)`
- `WithWriteQueueSize(...)`
- `WithPongTimeout(...)`
- `WithPingInterval(...)`
- `WithoutHeartbeat()`
- `WithCompression(true)`
- `WithOriginValidator(...)`

### 控制帧回调

- `WithCloseFrameHandler(...)`
- `WithPingHandler(...)`
- `WithPongHandler(...)`

## Connection API

### 身份与生命周期

- `ID() string`
- `Request() *http.Request`
- `Done() <-chan struct{}`
- `Closed() bool`
- `Close() error`

### 发送消息

- `Send(messageType, payload)`：阻塞直到消息成功进入发送队列
- `TrySend(messageType, payload)`：队列满时立即返回 `ErrSendQueueFull`
- `SendJSON(value)` / `TrySendJSON(value)`：发送 JSON 文本消息
- `SendPrepared(message)` / `TrySendPrepared(message)`：发送预构建消息

`Connection` 不再直接暴露底层 `*websocket.Conn`，这样可以从 API 层面禁止并发写入和随意绕过生命周期管理。

## Registry API

```go
registry := sinking_websocket.NewRegistry()

registry.Store("user-1", connection)

connection, ok := registry.Load("user-1")

registry.Range(func(id string, connection *sinking_websocket.Connection) bool {
	return true
})

registry.Delete("user-1")
registry.DeleteIfMatch("user-1", connection)
registry.Close()
```

`DeleteIfMatch` 的意义是避免旧连接断开时误删同一个用户的新连接。

## 广播示例

```go
func broadcast(registry *sinking_websocket.Registry, payload []byte) {
	result, err := registry.Broadcast(sinking_websocket.TextMessage, payload)
	if err != nil {
		return
	}
	_ = result
}
```

## 高频广播优化

如果同一份消息会发给很多连接，优先使用 `PrepareMessage`：

```go
prepared, err := sinking_websocket.PrepareMessage(
	sinking_websocket.TextMessage,
	[]byte(`{"type":"broadcast"}`),
)
if err != nil {
	return err
}

registry.Range(func(id string, connection *sinking_websocket.Connection) bool {
	if err := connection.TrySendPrepared(prepared); err != nil {
		registry.DeleteIfMatch(id, connection)
	}
	return true
})
```

如果是单条消息扇出到很多连接，更推荐直接走注册表内置广播：

```go
result := registry.BroadcastPrepared(prepared)
_ = result
```

## 行为保证

- 普通业务消息只会由单写协程写入底层连接
- 连接关闭时不会关闭发送队列，因此不会出现 `send on closed channel`
- `DisconnectHandler` 总会在连接生命周期结束时执行一次
- 回调 panic 会被收敛成错误并结束当前连接，不会把整个服务拖崩
- 高频广播默认走 `PreparedMessage`，避免每个连接重复编码同一条消息

# Sinking-WebSocket

🌐 **轻量级 WebSocket 连接管理框架**

Sinking-WebSocket 是一个基于 Go 语言开发的轻量级 WebSocket 连接管理框架，提供简洁易用的 API 来处理 WebSocket 连接的生命周期管理、消息处理和连接池管理。

## ✨ 主要特性

- 🔗 **连接管理**: 自动管理 WebSocket 连接的生命周期
- 📦 **连接池**: 内置连接池，支持连接的存储、获取和删除
- 🎯 **事件处理**: 支持连接、断开、消息和错误事件的自定义处理
- 🛡️ **并发安全**: 使用读写锁保证并发操作的安全性
- 🚀 **轻量级**: 基于 gorilla/websocket，代码简洁高效
- 📝 **简单易用**: 提供链式调用 API，使用简单直观

## 🛠️ 技术规格

- **Go 版本**: 1.11+
- **依赖**: gorilla/websocket v1.4.2
- **架构**: 轻量级，核心代码不到 200 行
- **并发**: 支持高并发连接管理

## 🚀 快速开始

### 安装

```bash
go get github.com/SinKingCloud/sinking-go/sinking-websocket
```

### 基本用法

```go
package main

import (
    "log"
    "net/http"
    "github.com/SinKingCloud/sinking-go/sinking-websocket"
)

func main() {
    // 创建连接池
    connections := sinking_websocket.NewWebSocketConnections()
    
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        // 获取客户端ID（可以从请求中获取）
        clientID := r.URL.Query().Get("id")
        if clientID == "" {
            clientID = "default"
        }
        
        // 创建 WebSocket 处理器
        ws := sinking_websocket.NewWebSocket().
            SetId(clientID).
            SetConnectHandle(func(id string, conn *sinking_websocket.Conn) {
                log.Printf("客户端 %s 连接成功", id)
                connections.Set(id, conn)
            }).
            SetOnMessageHandle(func(id string, conn *sinking_websocket.Conn, messageType int, data []byte) {
                log.Printf("收到客户端 %s 消息: %s", id, string(data))
                // 回复消息
                conn.WriteMessage(messageType, []byte("收到消息: "+string(data)))
            }).
            SetCloseHandle(func(id string, err error) {
                log.Printf("客户端 %s 断开连接: %v", id, err)
                connections.Delete(id)
            }).
            SetErrorHandle(func(id string, err error) {
                log.Printf("客户端 %s 发生错误: %v", id, err)
            })
        
        // 开始监听
        ws.Listen(w, r, nil)
    })
    
    log.Println("WebSocket 服务启动在 :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 广播消息

```go
// 向所有连接的客户端广播消息
func broadcast(connections *sinking_websocket.WebSocketConnections, message []byte) {
    allConns := connections.GetAll()
    for id, conn := range allConns {
        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Printf("向客户端 %s 发送消息失败: %v", id, err)
            connections.Delete(id)
        }
    }
}
```

## 📋 API 说明

### WebSocket 处理器

- `NewWebSocket()` - 创建 WebSocket 处理器
- `SetId(id)` - 设置连接 ID
- `SetConnectHandle(func)` - 设置连接成功回调
- `SetOnMessageHandle(func)` - 设置消息接收回调
- `SetCloseHandle(func)` - 设置连接关闭回调
- `SetErrorHandle(func)` - 设置错误处理回调
- `Listen(w, r, header)` - 开始监听连接

### 连接池管理

- `NewWebSocketConnections()` - 创建连接池
- `Get(key)` - 获取指定连接
- `GetAll()` - 获取所有连接
- `Set(key, conn)` - 存储连接
- `Delete(key)` - 删除连接

## 📖 文档

详细的使用文档和示例请参考 [doc.md](./doc.md)。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进项目！

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](../LICENSE) 文件。

## 📞 联系方式

- 作者: SinKingCloud
- 博客: www.clwl.online
- 项目地址: https://github.com/SinKingCloud/sinking-go

---

⭐ 如果这个项目对您有帮助，请给我们一个 Star！

# Sinking-WebSocket 使用文档

## 快速开始

### 基本服务器

```go
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/SinKingCloud/sinking-go/sinking-websocket"
)

func main() {
    connections := sinking_websocket.NewWebSocketConnections()
    
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        clientID := r.URL.Query().Get("id")
        
        ws := sinking_websocket.NewWebSocket().
            SetId(clientID).
            SetConnectHandle(func(id string, conn *sinking_websocket.Conn) {
                log.Printf("客户端连接: %s", id)
                connections.Set(id, conn)
            }).
            SetOnMessageHandle(func(id string, conn *sinking_websocket.Conn, messageType int, data []byte) {
                log.Printf("收到消息 [%s]: %s", id, string(data))
                conn.WriteMessage(messageType, data)
            }).
            SetCloseHandle(func(id string, err error) {
                log.Printf("客户端断开: %s, 错误: %v", id, err)
                connections.Delete(id)
            }).
            SetErrorHandle(func(id string, err error) {
                log.Printf("连接错误 [%s]: %v", id, err)
            })
        
        ws.Listen(w, r, nil)
    })
    
    log.Println("WebSocket 服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}
```

## 连接池管理

### 创建连接池

```go
// 创建连接池实例
connections := sinking_websocket.NewWebSocketConnections()
```

### 连接操作

```go
// 存储连接
connections.Set("user123", conn)

// 获取单个连接
conn := connections.Get("user123")
if conn != nil {
    conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
}

// 获取所有连接
allConns := connections.GetAll()
for id, conn := range allConns {
    log.Printf("连接ID: %s", id)
}

// 删除连接
success := connections.Delete("user123")
```

## WebSocket 处理器

### 创建处理器

```go
ws := sinking_websocket.NewWebSocket()
```

### 设置连接ID

```go
// 静态ID
ws.SetId("client-001")

// 动态ID（从请求中获取）
clientID := r.URL.Query().Get("client_id")
ws.SetId(clientID)
```

### 事件处理

#### 连接成功事件

```go
ws.SetConnectHandle(func(id string, conn *sinking_websocket.Conn) {
    log.Printf("新连接: %s", id)
    
    // 保存到连接池
    connections.Set(id, conn)
    
    // 发送欢迎消息
    welcome := map[string]string{
        "type":    "welcome",
        "message": "连接成功",
        "id":      id,
    }
    data, _ := json.Marshal(welcome)
    conn.WriteMessage(websocket.TextMessage, data)
})
```

#### 消息处理事件

```go
ws.SetOnMessageHandle(func(id string, conn *sinking_websocket.Conn, messageType int, data []byte) {
    log.Printf("收到消息 [%s]: %s", id, string(data))
    
    // 解析消息
    var msg map[string]interface{}
    if err := json.Unmarshal(data, &msg); err != nil {
        log.Printf("消息解析失败: %v", err)
        return
    }
    
    // 根据消息类型处理
    switch msg["type"] {
    case "ping":
        // 回复 pong
        response := map[string]string{"type": "pong"}
        responseData, _ := json.Marshal(response)
        conn.WriteMessage(websocket.TextMessage, responseData)
        
    case "broadcast":
        // 广播消息给所有客户端
        broadcastMessage(connections, msg["message"].(string))
        
    case "private":
        // 私信
        targetID := msg["target"].(string)
        message := msg["message"].(string)
        sendPrivateMessage(connections, targetID, message)
    }
})
```

#### 连接关闭事件

```go
ws.SetCloseHandle(func(id string, err error) {
    log.Printf("连接关闭 [%s]: %v", id, err)
    
    // 从连接池中移除
    connections.Delete(id)
    
    // 通知其他客户端用户离开
    notifyUserLeft(connections, id)
})
```

#### 错误处理事件

```go
ws.SetErrorHandle(func(id string, err error) {
    log.Printf("连接错误 [%s]: %v", id, err)
    
    // 记录错误日志
    // 可以根据错误类型进行不同处理
})
```

## 消息处理示例

### 广播消息

```go
func broadcastMessage(connections *sinking_websocket.WebSocketConnections, message string) {
    broadcast := map[string]string{
        "type":    "broadcast",
        "message": message,
        "time":    time.Now().Format("15:04:05"),
    }
    
    data, _ := json.Marshal(broadcast)
    allConns := connections.GetAll()
    
    for id, conn := range allConns {
        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            log.Printf("向 %s 发送广播失败: %v", id, err)
            connections.Delete(id)
        }
    }
}
```

### 私信功能

```go
func sendPrivateMessage(connections *sinking_websocket.WebSocketConnections, targetID, message string) {
    conn := connections.Get(targetID)
    if conn == nil {
        log.Printf("目标用户 %s 不在线", targetID)
        return
    }
    
    privateMsg := map[string]string{
        "type":    "private",
        "message": message,
        "time":    time.Now().Format("15:04:05"),
    }
    
    data, _ := json.Marshal(privateMsg)
    if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
        log.Printf("发送私信失败: %v", err)
        connections.Delete(targetID)
    }
}
```

### 用户状态通知

```go
func notifyUserLeft(connections *sinking_websocket.WebSocketConnections, leftUserID string) {
    notification := map[string]string{
        "type":    "user_left",
        "user_id": leftUserID,
        "time":    time.Now().Format("15:04:05"),
    }
    
    data, _ := json.Marshal(notification)
    allConns := connections.GetAll()
    
    for id, conn := range allConns {
        if id != leftUserID { // 不发送给已离开的用户
            conn.WriteMessage(websocket.TextMessage, data)
        }
    }
}
```

## 完整聊天室示例

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/websocket"
    "github.com/SinKingCloud/sinking-go/sinking-websocket"
)

type Message struct {
    Type    string `json:"type"`
    Content string `json:"content"`
    From    string `json:"from"`
    To      string `json:"to,omitempty"`
    Time    string `json:"time"`
}

func main() {
    connections := sinking_websocket.NewWebSocketConnections()
    
    // WebSocket 处理
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        username := r.URL.Query().Get("username")
        if username == "" {
            http.Error(w, "需要用户名", http.StatusBadRequest)
            return
        }
        
        ws := sinking_websocket.NewWebSocket().
            SetId(username).
            SetConnectHandle(func(id string, conn *sinking_websocket.Conn) {
                log.Printf("用户 %s 加入聊天室", id)
                connections.Set(id, conn)
                
                // 通知其他用户
                notifyUserJoined(connections, id)
                
                // 发送在线用户列表
                sendOnlineUsers(connections, id)
            }).
            SetOnMessageHandle(func(id string, conn *sinking_websocket.Conn, messageType int, data []byte) {
                var msg Message
                if err := json.Unmarshal(data, &msg); err != nil {
                    log.Printf("消息解析失败: %v", err)
                    return
                }
                
                msg.From = id
                msg.Time = time.Now().Format("15:04:05")
                
                switch msg.Type {
                case "public":
                    broadcastPublicMessage(connections, msg)
                case "private":
                    sendPrivateMessage(connections, msg)
                }
            }).
            SetCloseHandle(func(id string, err error) {
                log.Printf("用户 %s 离开聊天室", id)
                connections.Delete(id)
                notifyUserLeft(connections, id)
            }).
            SetErrorHandle(func(id string, err error) {
                log.Printf("用户 %s 连接错误: %v", id, err)
            })
        
        ws.Listen(w, r, nil)
    })
    
    // 静态文件服务
    http.Handle("/", http.FileServer(http.Dir("./static/")))
    
    log.Println("聊天室服务器启动在 :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func notifyUserJoined(connections *sinking_websocket.WebSocketConnections, username string) {
    msg := Message{
        Type:    "system",
        Content: username + " 加入了聊天室",
        Time:    time.Now().Format("15:04:05"),
    }
    
    data, _ := json.Marshal(msg)
    broadcastToAll(connections, data, username)
}

func notifyUserLeft(connections *sinking_websocket.WebSocketConnections, username string) {
    msg := Message{
        Type:    "system",
        Content: username + " 离开了聊天室",
        Time:    time.Now().Format("15:04:05"),
    }
    
    data, _ := json.Marshal(msg)
    broadcastToAll(connections, data, username)
}

func sendOnlineUsers(connections *sinking_websocket.WebSocketConnections, username string) {
    allConns := connections.GetAll()
    users := make([]string, 0, len(allConns))
    
    for id := range allConns {
        users = append(users, id)
    }
    
    msg := Message{
        Type:    "online_users",
        Content: "",
        Time:    time.Now().Format("15:04:05"),
    }
    
    // 这里可以扩展 Message 结构体来包含用户列表
    data, _ := json.Marshal(map[string]interface{}{
        "type":  "online_users",
        "users": users,
        "time":  msg.Time,
    })
    
    if conn := connections.Get(username); conn != nil {
        conn.WriteMessage(websocket.TextMessage, data)
    }
}

func broadcastPublicMessage(connections *sinking_websocket.WebSocketConnections, msg Message) {
    data, _ := json.Marshal(msg)
    broadcastToAll(connections, data, "")
}

func sendPrivateMessage(connections *sinking_websocket.WebSocketConnections, msg Message) {
    if msg.To == "" {
        return
    }
    
    data, _ := json.Marshal(msg)
    
    // 发送给目标用户
    if conn := connections.Get(msg.To); conn != nil {
        conn.WriteMessage(websocket.TextMessage, data)
    }
    
    // 也发送给发送者确认
    if conn := connections.Get(msg.From); conn != nil {
        conn.WriteMessage(websocket.TextMessage, data)
    }
}

func broadcastToAll(connections *sinking_websocket.WebSocketConnections, data []byte, excludeUser string) {
    allConns := connections.GetAll()
    
    for id, conn := range allConns {
        if id != excludeUser {
            if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
                log.Printf("向用户 %s 发送消息失败: %v", id, err)
                connections.Delete(id)
            }
        }
    }
}
```

## 客户端示例 (JavaScript)

```javascript
// 连接 WebSocket
const ws = new WebSocket('ws://localhost:8080/ws?username=用户名');

// 连接成功
ws.onopen = function(event) {
    console.log('连接成功');
};

// 接收消息
ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('收到消息:', message);
    
    switch(message.type) {
        case 'public':
            displayPublicMessage(message);
            break;
        case 'private':
            displayPrivateMessage(message);
            break;
        case 'system':
            displaySystemMessage(message);
            break;
        case 'online_users':
            updateOnlineUsers(message.users);
            break;
    }
};

// 发送公共消息
function sendPublicMessage(content) {
    const message = {
        type: 'public',
        content: content
    };
    ws.send(JSON.stringify(message));
}

// 发送私信
function sendPrivateMessage(to, content) {
    const message = {
        type: 'private',
        content: content,
        to: to
    };
    ws.send(JSON.stringify(message));
}

// 连接关闭
ws.onclose = function(event) {
    console.log('连接关闭');
};

// 连接错误
ws.onerror = function(error) {
    console.log('连接错误:', error);
};
```

---

更多高级用法和示例请参考源码。

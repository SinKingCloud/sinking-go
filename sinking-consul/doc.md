# Sinking Consul 使用文档

## 部署说明

### 环境要求
- Go 1.24+
- 内存 512MB+
- 开放端口 5678（或配置的其他端口）

### 快速部署

```bash
# 1. 进入服务端目录
cd sinking-consul/server

# 2. 编译项目
go build -o server main.go

# 3. 运行服务
./server run      # 前台运行
./server start    # 后台运行（Linux/macOS）
./server stop     # 停止服务
./server restart  # 重启服务
```

### 配置文件

编辑 `config/application.yml`：

```yaml
# 服务配置
server:
  mode: dev          # dev/prod
  host: 0.0.0.0     # 监听地址
  port: 5678        # 监听端口

# 认证配置  
auth:
  account: admin                                    # 管理员账户
  password: 123456                                  # 管理员密码
  expire: 86400                                     # JWT过期时间(秒)
  token: E10ADC3949BA59ABBE56E057F20F883E          # API访问令牌

# 集群配置
cluster:
  local: 127.0.0.1:5678    # 本地节点地址
  nodes:                    # 集群节点列表
    - 127.0.0.1:5678
    - 127.0.0.1:5679
```

### 集群部署

多节点部署时，每个节点的配置：

```yaml
# 节点1 (192.168.1.10)
cluster:
  local: 192.168.1.10:5678
  nodes:
    - 192.168.1.10:5678
    - 192.168.1.11:5678

# 节点2 (192.168.1.11)  
cluster:
  local: 192.168.1.11:5678
  nodes:
    - 192.168.1.10:5678
    - 192.168.1.11:5678
```

## 客户端使用

### 基本用法

```go
package main

import (
    "log"
    "github.com/SinKingCloud/sinking-go/sinking-consul/client"
)

func main() {
    // 创建客户端
    cli := client.NewClient(
        []string{"127.0.0.1:5678", "127.0.0.1:5679"}, // 集群地址
        "default",                                     // 服务组
        "user-service",                               // 服务名
        "192.168.1.100:8080",                        // 本机应用地址
        "E10ADC3949BA59ABBE56E057F20F883E",          // 集群密钥
    )

    // 连接集群
    _ = cli.Connect()
    defer cli.Close()

    // 获取所有服务
    _, _ = cli.GetAllService()

    // 获取所有配置
    _, _ = cli.GetAllConfigs()

    // 获取指定配置
    _, _ = cli.GetConfig("配置名")

    // 获取服务节点，支持轮询和随机
    _, _ = cli.GetService("服务名", client.Poll) // 轮询
    _, _ = cli.GetService("服务名", client.Rand) // 随机
}
```

### 配置解析

```go
// 获取配置解析器
configParser, err := cli.GetConfig("database.yml")
if err != nil {
    log.Fatal("获取配置失败:", err)
}

// 获取配置值
host := configParser.GetString("host")           // 字符串
port := configParser.GetInt("port")              // 整数
debug := configParser.GetBool("debug")           // 布尔值
timeout := configParser.GetDuration("timeout")  // 时间间隔

// 获取数组
ips := configParser.GetStringSlice("allowed_ips")
ports := configParser.GetIntSlice("ports")

// 检查配置是否存在
if configParser.IsSet("database.host") {
    // 配置存在
}
```

## 访问服务

部署完成后访问：
- **管理界面**: http://localhost:5678
- **默认账户**: admin / 123456

---

更多问题请查看日志文件或提交 Issue。

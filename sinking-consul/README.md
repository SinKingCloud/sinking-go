# Sinking Consul

🚀 **基于 Sinking-Web 框架的轻量级微服务注册与配置中心**

Sinking Consul 是一个使用自研 Sinking-Web 框架开发的轻量级微服务注册与配置中心，提供完整的服务发现、配置管理、集群监控功能。支持多集群管理，具备直观的 Web 管理界面，专为简化分布式系统的服务治理而设计。

## ✨ 核心功能

### 🌐 服务注册与发现
- 支持多节点服务注册，自动健康检查
- 实时同步服务状态，支持集群间数据同步
- 提供 RESTful API 接口，便于服务集成

### ⚙️ 配置中心
- 集中化配置存储与分发
- 支持配置分组管理，按环境隔离
- 配置变更实时推送，支持动态更新
- 多种配置格式支持（JSON、YAML、Properties等）

### 📊 集群管理
- 多集群统一管理，支持集群间数据同步
- 实时监控集群状态和节点健康度
- 分布式锁支持，保证数据一致性
- 集群节点自动发现和故障转移

### 🔐 安全与认证
- 基于 JWT 的身份验证系统
- 滑块验证码防护，增强安全性
- API 访问令牌控制，支持细粒度权限管理

## 🛠️ 技术架构

- **后端语言**: Go 1.24+
- **Web框架**: Sinking-Web (自研高性能 HTTP 框架)
- **数据库**: SQLite (轻量级嵌入式数据库)
- **缓存**: Go-Cache (内存缓存)
- **认证**: JWT + 滑块验证
- **前端**: React + TypeScript + Ant Design
- **部署**: 支持守护进程模式，跨平台运行

## 📁 项目结构

```
sinking-consul/
├── server/                    # 后端服务
│   ├── app/                  # 应用核心代码
│   │   ├── command/          # 命令处理（任务调度、队列管理）
│   │   ├── constant/         # 常量定义
│   │   ├── http/            # HTTP 层
│   │   │   ├── controller/   # 控制器（API、管理后台、认证）
│   │   │   ├── middleware/   # 中间件（CORS、权限检查）
│   │   │   └── route/       # 路由配置
│   │   ├── model/           # 数据模型（集群、配置、节点、日志）
│   │   ├── service/         # 业务服务层
│   │   │   ├── auth/        # 认证服务
│   │   │   ├── cluster/     # 集群管理服务
│   │   │   ├── config/      # 配置管理服务
│   │   │   ├── node/        # 节点管理服务
│   │   │   └── log/         # 日志服务
│   │   └── util/            # 工具类（缓存、验证码、数据库等）
│   ├── bootstrap/           # 系统启动引导
│   ├── config/             # 配置文件和数据库
│   ├── public/             # 静态资源和 SQL 脚本
│   └── main.go             # 程序入口
├── web/                     # 前端管理界面
├── client/                  # Go 客户端 SDK
└── README.md
```

## 🚀 快速开始

### 环境要求
- Go 1.24+ 
- Node.js 16+ (仅前端开发需要)

### 编译运行

```bash
# 进入服务端目录
cd sinking-consul/server

# 编译项目
go build -o server main.go

# Linux/macOS 运行
./server start    # 守护进程模式
./server run      # 前台运行模式
./server stop     # 停止服务
./server restart  # 重启服务

# Windows 或调试模式直接运行
go run main.go
```

### 访问服务
- **管理界面**: http://localhost:5678
- **默认账户**: admin / 123456

## ⚙️ 配置说明

配置文件位于 `server/config/application.yml`：

```yaml
# 服务配置
server:
  mode: dev          # 运行模式 (dev/prod)
  host: 0.0.0.0     # 监听地址
  port: 5678        # 监听端口

# 认证配置  
auth:
  account: admin                                    # 管理员账户
  password: 123456                                  # 管理员密码
  expire: 86400                                     # JWT 过期时间(秒)
  token: E10ADC3949BA59ABBE56E057F20F883E          # API 访问令牌

# 集群配置
cluster:
  local: 127.0.0.1:5678    # 本地节点地址(可选，不填自动获取)
  nodes:                    # 集群节点列表
    - 127.0.0.1:5678
    - 127.0.0.1:5679
```

## 📋 API 接口

### 服务注册 API
```bash
# 注册服务节点
POST /api/node/register
# 同步服务状态  
POST /api/node/sync
```

### 配置管理 API
```bash
# 同步配置
POST /api/config/sync
```

### 集群管理 API
```bash
# 集群注册
POST /api/cluster/register
# 节点测试
POST /api/cluster/testing
# 获取服务列表
GET /api/cluster/node
# 获取配置列表
GET /api/cluster/config
```

详细的 API 文档请参考 [doc.md](./doc.md)。

## 🤝 贡献指南

我们欢迎所有形式的贡献！

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](../LICENSE) 文件。

## 📞 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 [Issue](https://github.com/SinKingCloud/sinking-go/issues)
- 发起 [Discussion](https://github.com/SinKingCloud/sinking-go/discussions)

---

⭐ 如果这个项目对您有帮助，请给我们一个 Star！
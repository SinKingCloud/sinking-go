# Sinking-Go

🚀 **轻量级高性能微服务解决方案**

Sinking-Go 是一个基于 Go 语言开发的完整微服务解决方案，提供从 Web 框架、WebSocket 支持到服务注册发现的全套组件。专为构建高性能、低成本的分布式系统而设计，让开发者能够快速构建现代化的微服务应用。

## ✨ 核心特性

- 🌐 **完整的微服务生态**: 包含 Web 框架、WebSocket、服务注册发现等完整组件
- 🚀 **高性能设计**: 基于 Go 语言，提供卓越的并发性能
- 🛠️ **开箱即用**: 提供完整的示例和文档，快速上手
- 💰 **低成本方案**: 轻量级设计，资源占用少，降低运维成本
- 🔧 **灵活扩展**: 模块化设计，可根据需求选择使用
- 📱 **跨平台支持**: 支持 Linux、macOS、Windows 等多平台部署

## 📦 项目架构

### 🌐 [Sinking-Web](./sinking-web/) - HTTP 框架
**轻量级高性能 Web 框架**
- ✅ 高性能路由匹配（基于前缀树）
- ✅ 灵活的中间件机制
- ✅ 完整的参数绑定支持
- ✅ 静态文件服务
- ✅ 模板渲染支持
- ✅ HTTP/WebSocket 反向代理
- ✅ 请求限流控制

### 🔗 [Sinking-WebSocket](./sinking-websocket/) - WebSocket 框架
**WebSocket 连接管理框架**
- ✅ 连接池管理
- ✅ 事件驱动处理
- ✅ 并发安全设计
- ✅ 简洁易用 API

### 🏢 [Sinking-Consul](./sinking-consul/) - 服务治理中心
**微服务注册与配置中心**
- ✅ 服务注册与发现
- ✅ 配置中心管理
- ✅ 集群管理与同步
- ✅ Web 管理界面
- ✅ Go 客户端 SDK

### 📚 [Sinking-Demo](./sinking-demo/) - 使用示例
**完整的框架使用示例**
- ✅ Web 框架使用示例
- ✅ WebSocket 实现示例
- ✅ 中间件使用示例
- ✅ 代理功能示例

### 🐳 [Sinking-Env](./sinking-env/) - 开发环境
**Docker 开发环境整合**
- ✅ MySQL、Redis、Nginx 等常用服务
- ✅ 消息队列、搜索引擎等中间件
- ✅ 一键启动开发环境


## 🚀 快速开始

### 1. 环境准备

```bash
# 确保安装了 Go 1.21+ 和 Node.js 16+
go version
node --version
```

### 2. 克隆项目

```bash
git clone https://github.com/SinKingCloud/sinking-go.git
cd sinking-go
```

### 3. 运行示例

```bash
# 运行 Web 框架示例
cd sinking-demo
go mod tidy
go run main.go

# 访问 http://localhost:8888
```

### 4. 启动服务注册中心

```bash
# 启动 Consul 服务
cd sinking-consul/server
go build -o server main.go
./server run

# 访问管理界面 http://localhost:5678
# 默认账户: admin / 123456
```

### 5. 启动开发环境

```bash
# 启动完整开发环境
cd sinking-env
docker-compose up -d
```

## 📖 使用文档

### 核心框架文档
- [Sinking-Web 使用指南](./sinking-web/doc.md) - Web 框架完整使用文档
- [Sinking-WebSocket 使用指南](./sinking-websocket/doc.md) - WebSocket 框架使用文档
- [Sinking-Consul 使用指南](./sinking-consul/doc.md) - 服务治理中心使用文档

### 快速参考
- [Web 框架快速开始](./sinking-web/README.md)
- [WebSocket 框架快速开始](./sinking-websocket/README.md)
- [服务治理中心快速开始](./sinking-consul/README.md)

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    Sinking-Go 生态系统                      │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   前端应用   │ │   移动应用   │ │  第三方应用  │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
├─────────────────────────────────────────────────────────────┤
│                      API 网关层                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              Sinking-Web (HTTP + WebSocket)            │ │
│  │  • 路由管理  • 中间件  • 代理转发  • 限流控制          │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                      微服务层                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │  用户服务   │ │  订单服务   │ │  支付服务   │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
├─────────────────────────────────────────────────────────────┤
│                    服务治理层                               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                 Sinking-Consul                         │ │
│  │  • 服务注册  • 配置管理  • 健康检查  • 负载均衡        │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                     基础设施层                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐            │
│  │   数据库    │ │    缓存     │ │   消息队列   │            │
│  └─────────────┘ └─────────────┘ └─────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

## 🎯 适用场景

### 🌐 Web 应用开发
- RESTful API 服务
- 前后端分离应用
- 单页面应用 (SPA)
- 移动应用后端

### 🏢 微服务架构
- 服务注册与发现
- 配置中心管理
- API 网关
- 负载均衡

### 💬 实时通信
- 在线聊天系统
- 实时通知
- 数据推送
- 协作应用

### 🔄 系统集成
- 反向代理
- 服务编排
- 数据同步
- 第三方集成

## 🔄 开发路线

### 已完成 ✅
- [x] Sinking-Web HTTP 框架核心功能
- [x] Sinking-WebSocket 连接管理
- [x] Sinking-Consul 服务注册发现
- [x] 完整使用示例和文档
- [x] Docker 开发环境

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何参与
1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 贡献方式
- 🐛 报告 Bug
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码
- 🧪 编写测试用例

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](./LICENSE) 文件。

## 📞 联系我们

### 作者信息
- **作者**: SinKingCloud
- **博客**: [https://www.clwl.online](https://www.clwl.online)

### 社区交流
- **QQ**: 1178710004
- **微信**: sinking-cloud
- **GitHub**: [https://github.com/SinKingCloud/sinking-go](https://github.com/SinKingCloud/sinking-go)

### 问题反馈
- [提交 Issue](https://github.com/SinKingCloud/sinking-go/issues)
- [参与讨论](https://github.com/SinKingCloud/sinking-go/discussions)

---

⭐ **如果这个项目对您有帮助，请给我们一个 Star！**

让我们一起构建更好的微服务生态系统！ 🚀


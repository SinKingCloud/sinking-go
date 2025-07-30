# Sinking Consul

一个轻量级的微服务注册和配置中心。该项目提供了直观的 Web 界面，让用户可以轻松地管理多个集群,监控集群状态和健康度,配置服务发现和配置管理

## 技术栈

- **后端**: Go 1.21+
- **Web框架**: Gin
- **数据库**: SQLite
- **缓存**: Mem
- **认证**: JWT
- **前端**: Web界面（位于 web 目录）

## 项目结构

```
sinking-consul/
├── server/                 # 后端服务
│   ├── app/               # 应用核心代码
│   │   ├── http/          # HTTP 控制器和路由
│   │   ├── service/       # 业务逻辑层
│   │   ├── model/         # 数据模型
│   │   └── util/          # 工具类
│   ├── bootstrap/         # 启动配置
│   ├── config/           # 配置文件
│   └── public/           # 静态资源和SQL脚本
└── web/                  # 前端界面
```

## 配置文件说明 (application.yml)

### 服务器配置

```yaml
server:
  mode: dev #模式
  host: 0.0.0.0 #监听地址
  port: 5678 #监听端口
```

### 鉴权配置

```yaml
auth:
  account: admin #登录账户
  password: 123456 #登录密码
  expire: 86400 #登录过期时间（秒）
  token: E10ADC3949BA59ABBE56E057F20F883E #API令牌
```

### 集群配置

```yaml
cluster:
  local: 1.1.1.1:5678 #本地集群地址(不填则自动获取地址)
  nodes:
    - 111.111.111.111:8500 #服务器demo地址1
    - 222.222.222.222:8500 #服务器demo地址2
```

## 许可证

本项目采用 MIT 许可证，详见 LICENSE 文件。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进项目。

## 联系方式

如有问题或建议，请通过 Issue 联系我们。
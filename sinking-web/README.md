# Sinking-Web

🚀 **轻量级高性能 Go Web 框架**

Sinking-Web 是一个基于 Go 语言开发的轻量级、高性能 Web 框架，提供简洁易用的 API 和丰富的功能特性。框架设计简洁，性能优异，适合构建各种规模的 Web 应用和 API 服务。

## ✨ 主要特性

- 🚀 **高性能路由**: 基于前缀树的高效路由匹配算法
- 🛠️ **中间件支持**: 灵活的中间件机制，支持请求拦截和处理
- 📝 **参数绑定**: 支持 JSON、Form、Query、路径参数的自动绑定
- 🔧 **错误处理**: 统一的错误处理机制和自定义错误处理器
- 📁 **静态文件**: 内置静态文件服务支持
- 🌐 **代理支持**: 支持 HTTP 和 WebSocket 代理
- 📊 **限流控制**: 内置请求限流功能
- 🔍 **日志中间件**: 详细的请求日志记录
- 🛡️ **异常恢复**: 自动捕获和处理 panic 异常
- 🎯 **路由分组**: 支持路由分组和嵌套路由

## 🛠️ 技术规格

- **Go 版本**: 1.11+
- **依赖**: 无第三方依赖，仅使用 Go 标准库
- **架构**: 轻量级，核心代码不到 2000 行
- **性能**: 高性能路由匹配，支持高并发请求

## 🚀 快速开始

### 安装

```bash
go get github.com/SinKingCloud/sinking-go/sinking-web
```

### 基本用法

```go
package main

import "github.com/SinKingCloud/sinking-go/sinking-web"

func main() {
    // 创建引擎（包含日志和恢复中间件）
    r := sinking_web.Default()
    
    // 简单路由
    r.GET("/", func(c *sinking_web.Context) {
        c.String(200, "Hello, Sinking-Web!")
    })
    
    // JSON 响应
    r.GET("/json", func(c *sinking_web.Context) {
        c.JSON(200, sinking_web.H{
            "message": "Hello JSON",
            "status":  "success",
        })
    })
    
    // 启动服务器
    r.Run(":8080")
}
```

### 路由参数

```go
// 路径参数
r.GET("/user/:id", func(c *sinking_web.Context) {
    id := c.Param("id")
    c.String(200, "User ID: %s", id)
})

// 通配符路由
r.GET("/assets/*filepath", func(c *sinking_web.Context) {
    filepath := c.Param("filepath")
    c.String(200, "File path: %s", filepath)
})
```

### 参数绑定

```go
type User struct {
    Name  string `json:"name" form:"name"`
    Email string `json:"email" form:"email"`
}

r.POST("/user", func(c *sinking_web.Context) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(400, sinking_web.H{"error": err.Error()})
        return
    }
    c.JSON(200, user)
})
```

## 📋 支持的 HTTP 方法

- `GET` - 获取资源
- `POST` - 创建资源  
- `PUT` - 更新资源
- `DELETE` - 删除资源
- `PATCH` - 部分更新
- `HEAD` - 获取头信息
- `OPTIONS` - 获取选项
- `ANY` - 匹配所有方法

## 📖 文档

详细的使用文档和 API 说明请参考 [doc.md](./doc.md)。

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

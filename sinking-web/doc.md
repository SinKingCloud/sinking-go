# Sinking-Web 使用文档

## 快速开始

### 创建应用

```go
package main

import (
    "github.com/SinKingCloud/sinking-go/sinking-web"
    "time"
)

func main() {
    // 设置是否以debug模式运行
    sinking_web.SetDebugMode(false)
    
    // 设置读写超时时间
    sinking_web.SetTimeOut(10*time.Second, 10*time.Second)
    
    // 创建引擎（推荐使用 Default，包含日志和异常恢复中间件）
    r := sinking_web.Default()
    
    r.GET("/", func(c *sinking_web.Context) {
        c.String(200, "Hello World")
    })
    
    // 启动服务器
    r.Run("0.0.0.0:8080")
}
```

## 路由

### 基本路由

```go
r.GET("/get", getHandler)
r.POST("/post", postHandler)
r.PUT("/put", putHandler)
r.DELETE("/delete", deleteHandler)
r.PATCH("/patch", patchHandler)
r.HEAD("/head", headHandler)
r.OPTIONS("/options", optionsHandler)

// 匹配所有 HTTP 方法
r.ANY("/any", anyHandler)
```

### 路径参数

```go
// 单个参数
r.GET("/user/:id", func(c *sinking_web.Context) {
    id := c.Param("id")
    c.String(200, "User ID: %s", id)
})

// 多个参数
r.GET("/user/:id/post/:pid", func(c *sinking_web.Context) {
    id := c.Param("id")
    pid := c.Param("pid")
    c.String(200, "User: %s, Post: %s", id, pid)
})

// 通配符参数
r.GET("/files/*filepath", func(c *sinking_web.Context) {
    filepath := c.Param("filepath")
    c.String(200, "File: %s", filepath)
})
```

### 路由分组

```go
// 基本分组
v1 := r.Group("/api/v1")
{
    v1.GET("/users", getUsersHandler)
    v1.POST("/users", createUserHandler)
}

// 嵌套分组
api := r.Group("/api")
{
    v1 := api.Group("/v1")
    {
        v1.GET("/users", getUsersHandler)
    }
    
    v2 := api.Group("/v2")
    {
        v2.GET("/users", getUsersV2Handler)
    }
}
```

## 中间件

### 使用内置中间件

```go
r := sinking_web.New()

// 添加日志中间件
r.Use(sinking_web.Logger())

// 添加异常恢复中间件
r.Use(sinking_web.Recovery())
```

### 自定义中间件

```go
// 测试中间件
func TestMiddle() sinking_web.HandlerFunc {
    return func(c *sinking_web.Context) {
        log.Println("开始执行请求")
        c.Set("user", "admin") // 中间件传值
        c.Next()
        log.Println("请求执行完毕")
    }
}

// 限流中间件
func LimitRateMiddle() sinking_web.HandlerFunc {
    return func(c *sinking_web.Context) {
        // 令牌桶算法限流
        limitRate := sinking_web.GetLimitRateIns(c.ClientIP(false), 1) // 每秒颁发令牌总数
        
        // 方式1：等待式限流
        limitRate.Wait(1) // 消耗令牌数
        c.Next()          // 继续执行
        
        // 方式2：快速失败式限流
        // if limitRate.Check(1) {
        //     log.Println("触发限流")
        //     c.JSON(200, sinking_web.H{"code": 503, "message": "触发限流"})
        //     c.Abort() // 终止后面的方法执行
        // } else {
        //     c.Next() // 继续执行
        // }
    }
}

// 路由组中间件使用
method := r.Group("/method")
method.Use(LimitRateMiddle(), TestMiddle()) // 多个中间件
```

## 请求处理

### 获取请求参数

```go
// 获取所有类型参数的示例
method.ANY("/all", func(s *sinking_web.Context) {
    s.JSON(200, sinking_web.H{
        "code": "200", 
        "message": "获取所有请求参数及内容", 
        "data": sinking_web.H{
            "get":   s.AllQuery(), // 所有get参数
            "post":  s.AllForm(),  // 所有post参数
            "param": s.AllParam(), // 所有路径参数
        }
    })
})

// 获取单个参数
r.GET("/query", func(c *sinking_web.Context) {
    // Query 参数
    name := c.Query("name")                    // 获取参数
    age := c.DefaultQuery("age", "18")         // 带默认值
    all := c.AllQuery()                        // 获取所有 Query 参数
    
    c.JSON(200, sinking_web.H{
        "name": name,
        "age":  age,
        "all":  all,
    })
})

// 路径参数
r.POST("/user/:id", func(s *sinking_web.Context) {
    s.JSON(200, sinking_web.H{"code": "200", "message": "用户id为" + s.Param("id")})
})
```

### 获取请求头

```go
r.GET("/headers", func(c *sinking_web.Context) {
    userAgent := c.GetHeader("User-Agent")
    contentType := c.GetHeader("Content-Type")
    
    c.JSON(200, sinking_web.H{
        "user_agent":   userAgent,
        "content_type": contentType,
    })
})
```

### 获取客户端 IP

```go
r.GET("/ip", func(c *sinking_web.Context) {
    ip := c.ClientIP(true)  // true: 获取真实IP，false: 获取直连IP
    c.String(200, "Your IP: %s", ip)
})
```

### 获取请求体

```go
r.POST("/body", func(c *sinking_web.Context) {
    body := c.Body()
    c.String(200, "Body: %s", body)
})
```

## 参数绑定

### 参数绑定示例

```go
// 综合绑定示例（推荐）
r.ANY("/bind/:code", func(s *sinking_web.Context) {
    type Login struct {
        User string `default:"admin" json:"user"`   // default:默认值 json:json输出格式
        Pwd  string `default:"123456" json:"pwd"`   // default:默认值 json:json输出格式
        Code string `default:"000000" json:"code"`
    }
    login := &Login{}
    // BindAll 绑定所有参数（路径、Query、Form、JSON）
    err := s.BindAll(login)
    if err != nil {
        s.JSON(200, sinking_web.H{"code": "500", "message": "绑定参数失败", "data": err})
    } else {
        s.JSON(200, sinking_web.H{"code": "200", "message": "绑定参数成功", "data": login})
    }
})

// Query 参数绑定
r.ANY("/page", func(c *sinking_web.Context) {
    type ValidatePage struct {
        Page     uint `json:"page" default:"1"`
        PageSize int  `json:"page_size" default:"20"`
        Test     bool `json:"test" default:"true"`
    }
    pageInfo := &ValidatePage{}
    if c.BindQuery(pageInfo) != nil {
        fmt.Println("get参数绑定失败")
    } else {
        c.JSON(200, pageInfo)
    }
})

// JSON 绑定
r.POST("/json", func(s *sinking_web.Context) {
    s.JSON(200, sinking_web.H{"code": "200", "message": "请求为json", "data": s.Body()})
})
```

### 绑定方法说明

```go
// 可用的绑定方法：
// - BindAll(obj)    // 绑定所有参数（路径、Query、Form、JSON）
// - BindQuery(obj)  // 绑定 GET 参数
// - BindForm(obj)   // 绑定 POST 参数  
// - BindJSON(obj)   // 绑定 JSON 参数
// - BindParam(obj)  // 绑定路径参数

// 结构体标签说明：
// - `default:"value"` // 默认值
// - `json:"field"`    // 绑定字段名
```

## 响应处理

### 字符串响应

```go
r.GET("/string", func(c *sinking_web.Context) {
    c.String(200, "Hello %s", "World")
})
```

### JSON 响应

```go
r.GET("/json", func(c *sinking_web.Context) {
    c.JSON(200, sinking_web.H{
        "message": "success",
        "data":    []string{"item1", "item2"},
    })
})
```

### HTML 响应

```go
r.GET("/html", func(c *sinking_web.Context) {
    c.HTML(200, "<h1>Hello HTML</h1>")
})
```

### 文件响应

```go
r.GET("/file", func(c *sinking_web.Context) {
    c.File("./static/file.pdf")
})

r.GET("/download", func(c *sinking_web.Context) {
    c.FileAttachment("./static/file.pdf", "download.pdf")
})
```

### 重定向

```go
r.GET("/redirect", func(c *sinking_web.Context) {
    c.Redirect(302, "https://www.example.com")
})
```

### 设置响应头

```go
r.GET("/headers", func(c *sinking_web.Context) {
    c.SetHeader("X-Custom-Header", "CustomValue")
    c.SetHeader("Cache-Control", "no-cache")
    c.JSON(200, sinking_web.H{"message": "success"})
})
```

## 文件上传

### 单文件上传

```go
r.POST("/upload", func(c *sinking_web.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, sinking_web.H{"error": err.Error()})
        return
    }
    
    // 保存文件
    if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
        c.JSON(500, sinking_web.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, sinking_web.H{
        "message":  "上传成功",
        "filename": file.Filename,
        "size":     file.Size,
    })
})
```

### 多文件上传

```go
r.POST("/uploads", func(c *sinking_web.Context) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, sinking_web.H{"error": err.Error()})
        return
    }
    
    files := form.File["files"]
    for _, file := range files {
        if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
            c.JSON(500, sinking_web.H{"error": err.Error()})
            return
        }
    }
    
    c.JSON(200, sinking_web.H{
        "message": "上传成功",
        "count":   len(files),
    })
})
```

## 静态文件服务

```go
// 静态资源目录
r.Static("/static", "./static/") // 访问地址：ip:port/static/test.txt
```

## 模板渲染

```go
import (
    "html/template"
    "time"
)

// 自定义模板函数
r.SetFuncMap(template.FuncMap{
    // 格式转换器
    "DateTimeNow": func(t time.Time) string {
        year, month, day := t.Date()
        return fmt.Sprintf("%d-%02d-%02d", year, month, day)
    },
})

// 加载模板目录
r.LoadHtmlGlob("template/*")

r.GET("/template", func(s *sinking_web.Context) {
    // 使用模板
    s.HTML(200, "index.tmpl", sinking_web.H{
        "name": "test name",
        "now":  time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
    })
})
```

## 代理功能

### 自定义 HTTP 反向代理

```go
r.GET("/proxyHttp/*", func(s *sinking_web.Context) {
    _ = s.HttpProxy("http://127.0.0.1:1004", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
        // 过滤器 可以执行自定义过滤或修改内容
    }, nil)
})
```

### 自定义 WebSocket 反向代理

```go
r.GET("/proxyWs/*", func(s *sinking_web.Context) {
    _ = s.WebSocketProxy("ws://127.0.0.1:1004/test/1", nil, func(r *http.Request, w http.ResponseWriter) {
        // 过滤器 可以执行自定义过滤或修改内容
    }, nil)
})
```

### 自定义通用反向代理

```go
g := r.Group("/proxyAll")
g.GET("/*", func(s *sinking_web.Context) {
    // 支持 ws 和 http
    s.Proxy("ws://127.0.0.1:1004/test/1", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
        // 过滤器 可以执行自定义过滤或修改内容
        r.URL.Path = strings.Replace(s.Request.URL.Path, "/proxyAll", "", 1)
        r.URL.RawPath, r.RequestURI = r.URL.Path, r.URL.Path
    }, nil)
})
```

### 通用反向代理

```go
// HTTP 代理
r.PROXY("/proxy/http", "http://127.0.0.1:1004", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
    // 过滤器 可以执行自定义过滤或修改内容
}, nil)

// WebSocket 代理
r.PROXY("/proxy/ws", "ws://127.0.0.1:1004/test/1", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
    // 过滤器 可以执行自定义过滤或修改内容
}, nil)
```

## 错误处理

### 自定义错误处理

```go
// 设置错误回调
r.SetErrorHandle(&sinking_web.ErrorHandel{
    // 资源不存在错误
    NotFound: func(c *sinking_web.Context) {
        c.JSON(404, sinking_web.H{"code": 404, "message": "资源不存在"})
    },
    // 系统错误
    Fail: func(c *sinking_web.Context, code int, message string) {
        c.JSON(500, sinking_web.H{"code": code, "message": message})
    },
})
```

### 手动触发错误

```go
r.GET("/error", func(c *sinking_web.Context) {
    c.Fail(500, "服务器内部错误")
})
```

## 限流控制

```go
// 令牌桶算法限流
func LimitRateMiddle() sinking_web.HandlerFunc {
    return func(c *sinking_web.Context) {
        // 获取限流实例（基于客户端IP）
        limitRate := sinking_web.GetLimitRateIns(c.ClientIP(false), 1) // 每秒颁发令牌总数
        
        // 方式1：等待式限流
        limitRate.Wait(1) // 消耗令牌数
        c.Next()          // 继续执行
        
        // 方式2：快速失败式限流
        // if limitRate.Check(1) {
        //     log.Println("触发限流")
        //     c.JSON(200, sinking_web.H{"code": 503, "message": "触发限流"})
        //     c.Abort() // 终止后面的方法执行
        // } else {
        //     c.Next() // 继续执行
        // }
    }
}

// 使用限流中间件
r.Use(LimitRateMiddle())
```

## 上下文操作

### 存储和获取数据

```go
// 中间件中设置数据
func TestMiddle() sinking_web.HandlerFunc {
    return func(c *sinking_web.Context) {
        c.Set("user", "admin") // 中间件传值
        c.Next()
    }
}

// 处理器中获取数据
r.GET("/user", func(c *sinking_web.Context) {
    user := c.Get("user") // 中间件取值
    c.JSON(200, sinking_web.H{"user": user})
})
```

### 判断请求是否被中止

```go
r.Use(func(c *sinking_web.Context) {
    if someCondition {
        c.Abort() // 中止后续处理
        return
    }
    c.Next()
})

r.GET("/test", func(c *sinking_web.Context) {
    if c.IsAborted() {
        return
    }
    c.String(200, "OK")
})
```

## 启动服务

### 基本启动

```go
// 启动 HTTP 服务
r.Run(":8080")
```

### HTTPS 启动

```go
// 启动 HTTPS 服务
r.RunTLS(":8443", "cert.pem", "key.pem")
```

### 自定义服务器

```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      r,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}

log.Fatal(server.ListenAndServe())
```

---

更多高级用法和示例请参考源码和测试文件。

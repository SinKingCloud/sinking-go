package main

import (
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"github.com/SinKingCloud/sinking-go/sinking-websocket"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// TestMiddle 测试中间件
func TestMiddle() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		log.Println("开始执行请求")
		c.Set("user", "admin") // 中间件传值
		//c.Get("user")// 中间件取值 后面的中间件可通过get set 方法传值
		c.Next()
		log.Println("请求执行完毕")
	}
}

// LimitRateMiddle 限流
func LimitRateMiddle() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		//令牌桶算法限流
		limitRate := sinking_web.GetLimitRateIns(c.ClientIP(false), 1) //每秒颁发令牌总数
		mode := 0
		switch mode {
		case 0:
			//1.等待式限流
			limitRate.Wait(1) //消耗令牌数
			c.Next()          //继续执行
		case 1:
			//2.快速失败式限流
			if limitRate.Check(1) {
				log.Println("触发限流")
				c.JSON(200, sinking_web.H{"code": 503, "message": "触发限流"})
				c.Abort() //终止后面的方法执行
			} else {
				log.Println("请求成功")
				c.Next() //继续执行
			}
		}
	}
}

//示例使用说明
//1.项目目录启动go mod集成
//2.项目目录执行go get -u github.com/SinKingCloud/sinking-go/sinking-web
//3.项目目录执行go get -u github.com/SinKingCloud/sinking-go/sinking-websocket
//4.复制本示例代码执行main函数

func main() {
	//设置是否以debug模式运行
	sinking_web.SetDebugMode(true)

	//设置读写超时时间
	sinking_web.SetTimeOut(10*time.Second, 10*time.Second)

	//实例化一个http server
	r := sinking_web.Default()

	//设置错误回调
	r.SetErrorHandle(&sinking_web.ErrorHandel{
		//资源不存在错误
		NotFound: func(c *sinking_web.Context) {
			c.JSON(404, sinking_web.H{"code": 404, "message": "资源不存在"})
		},
		//系统错误
		Fail: func(c *sinking_web.Context, code int, message string) {
			c.JSON(500, sinking_web.H{"code": code, "message": message})
		},
	})

	//静态资源
	r.Static("/static", "./static/") //测试访问地址 ip:port/static/test.txt

	//模板
	r.SetFuncMap(template.FuncMap{
		//格式转换器
		"DateTimeNow": func(t time.Time) string {
			year, month, day := t.Date()
			return fmt.Sprintf("%d-%02d-%02d", year, month, day)
		},
	})
	r.LoadHtmlGlob("template/*") //模板目录
	r.GET("/template", func(s *sinking_web.Context) {
		//使用模板
		s.HTML(200, "index.tmpl", sinking_web.H{
			"name": "test name",
			"now":  time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	//路由请求及中间件示例(仅写出常用示例)
	method := r.Group("/method")                // 路由组
	method.Use(LimitRateMiddle(), TestMiddle()) // 中间件
	{
		method.ANY("/any", func(s *sinking_web.Context) {
			s.JSON(200, sinking_web.H{"code": "200", "message": "所有请求都可以捕获"})
		}) // 访问地址 ip:port/method/any
		method.GET("/get", func(s *sinking_web.Context) {
			s.JSON(200, sinking_web.H{"code": "200", "message": "请求为" + s.Method})
		}) // 访问地址 ip:port/method/get
		method.POST("/post", func(s *sinking_web.Context) {
			s.JSON(200, sinking_web.H{"code": "200", "message": "请求为" + s.Method})
		}) // 访问地址 ip:port/method/post
		method.POST("/json", func(s *sinking_web.Context) {
			s.JSON(200, sinking_web.H{"code": "200", "message": "请求为json", "data": s.Body()})
		}) // 访问地址 ip:port/method/json
		method.POST("/user/:id", func(s *sinking_web.Context) {
			s.JSON(200, sinking_web.H{"code": "200", "message": "用户id为" + s.Param("id")})
		}) // 访问地址 ip:port/method/user/[参数]
		method.ANY("/all", func(s *sinking_web.Context) {
			s.JSON(200, sinking_web.H{"code": "200", "message": "获取所有请求参数及内容", "data": sinking_web.H{
				"get":   s.AllQuery(), //所有get参数
				"post":  s.AllForm(),  //所有post参数
				"param": s.AllParam(), //所有路径参数
			}})
		}) // 访问地址 ip:port/method/all
	}

	//参数绑定 访问地址 ip:port/bind?user=admin
	r.GET("/bind", func(s *sinking_web.Context) {
		type Login struct {
			User string `form:"user" default:"admin" json:"user"` //form:接受的参数名 default:默认值 json:json输出格式
			Pwd  string `form:"pwd" default:"123456" json:"pwd"`  //form:接受的参数名 default:默认值 json:json输出格式
		}
		login := &Login{}
		err := s.BindQuery(&login) //BindQuery:绑定get参数 BindForm:绑定post参数 BindJson:绑定json BindParam:绑定路由参数
		if err != nil {
			s.JSON(200, sinking_web.H{"code": "500", "message": "绑定参数失败"})
		} else {
			s.JSON(200, sinking_web.H{"code": "200", "message": "绑定参数成功", "data": login})
		}
	})
	//websocket功能
	wsConn := sinking_websocket.NewWebSocketConnections() //ws连接池
	ws := r.Group("/ws")
	//ws本质是get长连接,可使用get建立短连接在升级为长连接最后使用协程监听消息
	ws.GET("/message/listen/:id", func(s *sinking_web.Context) {
		//生成uid
		uid := "user-" + s.Param("id")
		wsServer := sinking_websocket.WebSocket{
			Id: uid,
			OnError: func(id string, err error) {
				wsConn.Delete(id)
				log.Println("websocket错误", err)
			},
			OnConnect: func(id string, ws *sinking_websocket.Conn) {
				wsConn.Set(id, ws)
				log.Println("websocket连接", uid)
			},
			OnClose: func(id string, err error) {
				wsConn.Delete(id)
				log.Println("websocket关闭", err)
			},
			OnMessage: func(id string, ws *sinking_websocket.Conn, messageType int, data []byte) {
				log.Println("websocket消息", string(data), messageType)
				conn := wsConn.Get(id)
				if conn != nil {
					_ = conn.WriteMessage(1, data)
				}
			},
		}
		wsServer.Listen(s.Writer, s.Request, nil)
		//ws地址 ws://ip:port/ws/message/listen/[示例ID]
		//在线测试ws工具 http://coolaf.com/zh/tool/chattest
	})
	//单播消息
	ws.GET("/message/send/:id/:message", func(s *sinking_web.Context) {
		uid := "user-" + s.Param("id")
		conn := wsConn.Get(uid)
		if conn != nil {
			err := conn.WriteMessage(1, []byte(s.Param("message")))
			if err != nil {
				s.JSON(200, sinking_web.H{"code": "500", "message": "发送消息失败"})
			} else {
				s.JSON(200, sinking_web.H{"code": "500", "message": "发送消息成功"})
			}
		} else {
			s.JSON(200, sinking_web.H{"code": "500", "message": "发送消息失败"})
		}
	})
	//广播消息
	ws.GET("/message/send/:message", func(s *sinking_web.Context) {
		for _, v := range wsConn.GetAll() {
			if v != nil {
				_ = v.WriteMessage(1, []byte(s.Param("message")))
			}
		}
		s.JSON(200, sinking_web.H{"code": "500", "message": "发送消息成功"})
	})
	//反向代理功能
	//(1.)自定义http反向代理
	r.GET("/proxyHttp/*", func(s *sinking_web.Context) {
		s.HttpProxy("http://127.0.0.1:1004", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
			//过滤器 可以执行自定义过滤或修改内容
		}, nil)
	})
	//(2.)自定义websocket反向代理
	r.GET("/proxyWs/*", func(s *sinking_web.Context) {
		s.WebSocketProxy("ws://127.0.0.1:1004/test/1", nil, func(r *http.Request, w http.ResponseWriter) {
			//过滤器 可以执行自定义过滤或修改内容
		}, nil)
	})
	//(2.)自定义通用反向代理
	g := r.Group("/proxyAll")
	g.GET("/*", func(s *sinking_web.Context) {
		//支持ws 和 http
		s.Proxy("ws://127.0.0.1:1004/test/1", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
			//过滤器 可以执行自定义过滤或修改内容
			r.URL.Path = strings.Replace(s.Request.URL.Path, "/proxyAll", "", 1)
			r.URL.RawPath, r.RequestURI = r.URL.Path, r.URL.Path
		}, nil)
	})
	//(4.)通用反向代理
	r.PROXY("/proxy/http", "http://127.0.0.1:1004", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
		//过滤器 可以执行自定义过滤或修改内容
	}, nil)
	r.PROXY("/proxy/ws", "ws://127.0.0.1:1004/test/1", nil, func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
		//过滤器 可以执行自定义过滤或修改内容
	}, nil)
	//启动http server
	err := r.Run("0.0.0.0:8888")
	if err != nil {
		log.Println(err.Error())
		return
	}
}

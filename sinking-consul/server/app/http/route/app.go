package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/http/controller/api"
	"server/app/http/controller/auth"
	"server/app/http/controller/system"
	"server/app/http/middleware"
	"server/app/util"
	"server/app/util/server"
	"server/public"
)

func loadApp(s *sinking_web.Engine) {
	if util.IsDebug() {
		s.Use(server.HandleFunc(middleware.Cors))
	}
	loadApiRoute(s)
	loadAuthRoute(s)
	loadSystemRoute(s)
	loadStaticRoute(s)
}

func loadStaticRoute(s *sinking_web.Engine) {
	s.ANY("/", server.HandleFunc(func(c *server.Context) {
		c.SetHeader("content-type", "text/html;charset=utf-8;")
		c.Data(200, public.ReadDistFile("index.html"))
	}))
	s.ANY("/*", server.HandleFunc(func(c *server.Context) {
		c.Request.URL.Path = public.Path() + c.Request.URL.Path
		public.FileServer.ServeHTTP(c.Writer, c.Request)
	}))
}

func loadApiRoute(s *sinking_web.Engine) {
	g := s.Group("/api")
	g.Use(server.HandleFunc(middleware.CheckToken))

	//集群中心相关路由
	cluster := g.Group("/cluster")
	cluster.ANY("/register", server.HandleFunc(api.Cluster.Register)) //注册集群
	cluster.ANY("/node", server.HandleFunc(api.Cluster.Node))         //服务列表
	cluster.ANY("/config", server.HandleFunc(api.Cluster.Config))     //配置列表

	//注册中心相关路由
	node := g.Group("/node")
	node.ANY("/register", server.HandleFunc(api.Node.Register)) //注册服务

	//配置中心相关路由
	config := g.Group("/config")
	config.ANY("/list", nil) //可用配置
}

func loadAuthRoute(s *sinking_web.Engine) {
	s.ANY("/login", server.HandleFunc(auth.Login))     //账号登录
	s.ANY("/logout", server.HandleFunc(auth.Logout))   //注销登录
	s.ANY("/captcha", server.HandleFunc(auth.Captcha)) //验证码
}

func loadSystemRoute(s *sinking_web.Engine) {
	g := s.Group("/system")
	g.Use(server.HandleFunc(middleware.CheckLogin))
	g.ANY("/enum", server.HandleFunc(system.Enum)) //枚举类型
}

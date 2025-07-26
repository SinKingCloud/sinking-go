package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/http/controller/auth"
	"server/app/http/controller/config"
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
	loadAuthRoute(s)
	loadConfigRoute(s)
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

func loadAuthRoute(s *sinking_web.Engine) {
	s.ANY("/login", server.HandleFunc(auth.Login))     //账号登录
	s.ANY("/logout", server.HandleFunc(auth.Logout))   //注销登录
	s.ANY("/captcha", server.HandleFunc(auth.Captcha)) //验证码
}

func loadConfigRoute(s *sinking_web.Engine) {
	g := s.Group("/config")
	g.Use(server.HandleFunc(middleware.CheckLogin))
	g.ANY("/get", server.HandleFunc(config.Get)) //获取配置
	g.ANY("/set", server.HandleFunc(config.Set)) //修改配置
}

func loadSystemRoute(s *sinking_web.Engine) {
	g := s.Group("/system")
	g.Use(server.HandleFunc(middleware.CheckLogin))
	g.ANY("/log", server.HandleFunc(system.Log))   //系统日志
	g.ANY("/enum", server.HandleFunc(system.Enum)) //枚举类型
}

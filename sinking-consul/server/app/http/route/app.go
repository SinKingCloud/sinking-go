package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/http/controller/admin"
	"server/app/http/controller/api"
	"server/app/http/controller/auth"
	"server/app/http/middleware"
	"server/app/util/context"
	"server/public"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadMiddleware(s)
	loadApiRoute(s)
	loadAuthRoute(s)
	loadAdminRoute(s)
	loadStaticRoute(s)
}

// loadMiddleware 加载中间件
func loadMiddleware(s *sinking_web.Engine) {
	s.Use(context.HandleFunc(middleware.Cors))
}

// loadStaticRoute 加载静态资源路由
func loadStaticRoute(s *sinking_web.Engine) {
	s.ANY("/", context.HandleFunc(func(c *context.Context) {
		c.SetHeader("content-type", "text/html;charset=utf-8;")
		c.Data(200, public.ReadDistFile("index.html"))
	}))
	s.ANY("/*", context.HandleFunc(func(c *context.Context) {
		c.Request.URL.Path = public.Path() + c.Request.URL.Path
		public.FileServer.ServeHTTP(c.Writer, c.Request)
	}))
}

// loadApiRoute 加载api路由
func loadApiRoute(s *sinking_web.Engine) {
	g := s.Group("/api")
	g.Use(context.HandleFunc(middleware.CheckToken))

	//集群中心相关路由
	cluster := g.Group("/cluster")
	{
		cluster.ANY("/testing", context.HandleFunc(api.Cluster.Testing))   //节点测试
		cluster.ANY("/register", context.HandleFunc(api.Cluster.Register)) //注册集群
		cluster.ANY("/node", context.HandleFunc(api.Cluster.Node))         //服务列表
		cluster.ANY("/config", context.HandleFunc(api.Cluster.Config))     //配置列表
		cluster.ANY("/lock", context.HandleFunc(api.Cluster.Lock))         //分布式锁
		cluster.ANY("/delete", context.HandleFunc(api.Cluster.Delete))     //删除数据
		cluster.ANY("/create", context.HandleFunc(api.Cluster.Create))     //创建数据
		cluster.ANY("/update", context.HandleFunc(api.Cluster.Update))     //更新数据
	}

	//注册中心相关路由
	node := g.Group("/node")
	{
		node.ANY("/register", context.HandleFunc(api.Node.Register)) //注册服务
		node.ANY("/sync", context.HandleFunc(api.Node.Sync))         //同步服务
	}

	//配置中心相关路由
	config := g.Group("/config")
	{
		config.ANY("/sync", context.HandleFunc(api.Config.Sync)) //同步配置
	}
}

// loadAuthRoute 加载认证路由
func loadAuthRoute(s *sinking_web.Engine) {
	g := s.Group("/auth")
	{
		g.ANY("/login", context.HandleFunc(auth.Login))     //账号登录
		g.ANY("/logout", context.HandleFunc(auth.Logout))   //注销登录
		g.ANY("/captcha", context.HandleFunc(auth.Captcha)) //滑块验证
		g.ANY("/info", context.HandleFunc(auth.Info))       //站点信息
	}
}

// loadAdminRoute 加载管理后台路由
func loadAdminRoute(s *sinking_web.Engine) {
	g := s.Group("/admin")
	g.Use(context.HandleFunc(middleware.CheckLogin))

	person := g.Group("/person")
	{
		person.ANY("/info", context.HandleFunc(admin.Person.Info))         //个人信息
		person.ANY("/password", context.HandleFunc(admin.Person.Password)) //修改密码
		person.ANY("/log", context.HandleFunc(admin.Person.Log))           //操作日志
	}

	system := g.Group("/system")
	{
		system.ANY("/overview", context.HandleFunc(admin.System.Overview)) //统计数据
		system.ANY("/enum", context.HandleFunc(admin.System.Enum))         //枚举类型
		system.ANY("/config", context.HandleFunc(admin.System.Config))     //系统配置
	}

	cluster := g.Group("/cluster")
	{
		cluster.ANY("/list", context.HandleFunc(admin.Cluster.List)) //集群列表
	}

	node := g.Group("/node")
	{
		node.ANY("/list", context.HandleFunc(admin.Node.List))     //节点列表
		node.ANY("/update", context.HandleFunc(admin.Node.Update)) //编辑节点
		node.ANY("/delete", context.HandleFunc(admin.Node.Delete)) //删除节点
	}

	config := g.Group("/config")
	{
		config.ANY("/list", context.HandleFunc(admin.Config.List))     //配置列表
		config.ANY("/info", context.HandleFunc(admin.Config.Info))     //配置信息
		config.ANY("/update", context.HandleFunc(admin.Config.Update)) //编辑配置
		config.ANY("/create", context.HandleFunc(admin.Config.Create)) //创建配置
		config.ANY("/delete", context.HandleFunc(admin.Config.Delete)) //删除配置
	}
}

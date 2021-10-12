package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"strconv"
	"time"
)

// loadCommonMiddleware 加载全局中间件
func loadCommonMiddleware(s *sinking_web.Engine) {

}

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	InitApiRouter(s)
	InitAdminRouter(s)
}

// loadErrorHandle 设置错误回调
func loadErrorHandle(s *sinking_web.Engine) {
	//设置错误回调
	s.SetErrorHandle(&sinking_web.ErrorHandel{
		//资源不存在错误
		NotFound: func(c *sinking_web.Context) {
			c.JSON(404, sinking_web.H{"code": 404, "message": "请求资源不存在"})
		},
		//系统错误
		Fail: func(c *sinking_web.Context, code int, message string) {
			c.JSON(500, sinking_web.H{"code": code, "message": message})
		},
	})
}

func Init() {
	//设置是否以debug模式运行
	sinking_web.SetDebugMode(setting.GetSystemConfig().App.Debug)
	//设置读写超时时间
	sinking_web.SetTimeOut(600*time.Second, 600*time.Second)
	//实例化一个http server
	r := sinking_web.Default()
	//加载错误handle
	loadErrorHandle(r)
	//加载全局中间件
	loadCommonMiddleware(r)
	//加载路由
	loadRoute(r)
	//启动http server
	err := r.Run(setting.GetSystemConfig().App.Ip + ":" + strconv.Itoa(setting.GetSystemConfig().App.Port))
	if err != nil {
		logs.Println(err.Error())
		return
	}
}

package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"os"
	"os/signal"
	"server/app/util"
	"server/app/util/http"
	"strconv"
	"syscall"
)

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

// Init 初始化server
func Init() {
	host, port := util.ServerAddr()
	addr := host + ":" + strconv.Itoa(port)
	r := http.NewServer(addr, util.IsDebug())
	r.ErrorHandle(&sinking_web.ErrorHandel{
		NotFound: func(c *sinking_web.Context) {
			c.JSON(404, sinking_web.H{"code": 404, "message": "请求资源不存在"})
		},
		Fail: func(c *sinking_web.Context, code int, message string) {
			c.JSON(500, sinking_web.H{"code": code, "message": message})
		},
	})
	r.Handle(func(engine *sinking_web.Engine) {
		loadErrorHandle(engine)
		loadApp(engine)
	})
	_ = r.Start()
	//等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	_ = <-sigChan
	_ = r.Stop()
}

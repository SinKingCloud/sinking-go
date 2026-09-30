package route

import (
	"fmt"
	"net"
	"strconv"

	"server/app/util/http"
	"server/global"

	"github.com/SinKingCloud/sinking-go/sinking-web"
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
func Init(stop <-chan struct{}) {
	host, port := global.App.ServerAddr()
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	r := http.NewServer(addr, global.App.IsDebug())
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
	if err := r.Start(); err != nil {
		panic(fmt.Errorf("启动HTTP服务失败: %w", err))
	}
	select {
	case err := <-r.Done():
		if err != nil {
			panic(fmt.Errorf("HTTP服务运行失败: %w", err))
		}
		return
	case <-stop:
	}
	if err := r.Stop(); err != nil {
		panic(fmt.Errorf("停止HTTP服务失败: %w", err))
	}
}

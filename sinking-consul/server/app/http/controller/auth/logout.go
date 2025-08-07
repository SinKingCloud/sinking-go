package auth

import (
	"server/app/http/middleware"
	"server/app/service"
	"server/app/util/server"
)

// Logout 退出登录
func Logout(c *server.Context) {
	middleware.CheckLogin(c)
	if c.IsAborted() {
		return
	}
	_ = service.Auth.ClearLoginToken()
	c.Success("注销成功")
}

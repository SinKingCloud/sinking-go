package auth

import (
	"server/app/http/middleware"
	"server/app/service"
	"server/app/service/log"
	"server/app/util/server"
)

// Logout 退出登录
func Logout(c *server.Context) {
	middleware.CheckLogin(c)
	if c.IsAborted() {
		return
	}
	_ = service.Auth.ClearLoginToken()
	service.Log.Create(c.GetRequestIp(), log.EventLogin, "注销账户登录", "注销登录成功")
	c.Success("注销成功")
}

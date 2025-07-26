package middleware

import (
	"server/app/constant"
	"server/app/service"
	"server/app/util/jwt"
	"server/app/util/server"
)

// CheckLogin 判断登录
func CheckLogin(c *server.Context) {
	token := c.Request.Header.Get(constant.JwtTokenName)
	types := c.Request.Header.Get(constant.JwtDeviceName)
	if token == "" || types == "" {
		c.NotLogin("您还未登陆,请先登陆账户", nil)
		c.Abort()
		return
	}
	key := jwt.CheckToken(token)
	if key == nil || key.User == nil {
		c.TokenError("登陆超时,请重新登陆", nil)
		c.Abort()
		return
	}
	loginToken := service.Config.Get(constant.LoginGroup, constant.LoginToken+"."+types)
	if loginToken == "" {
		c.TokenError("您的账户已注销登陆,请重新登陆", nil)
		c.Abort()
		return
	}
	if key.User.LoginToken != loginToken {
		c.TokenError("您的账户已在其他设备登陆,请重新登陆", nil)
		c.Abort()
		return
	}
	c.SetUserInfo(key.User)
	c.Next()
}

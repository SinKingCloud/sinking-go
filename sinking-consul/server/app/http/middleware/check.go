package middleware

import (
	"server/app/constant"
	"server/app/service"
	"server/app/util/context"
	"server/app/util/jwt"
)

// CheckLogin 判断登录
func CheckLogin(c *context.Context) {
	token := c.Request.Header.Get(constant.JwtTokenName)
	if token == "" {
		c.NotLogin("您还未登陆,请先登陆账户")
		c.Abort()
		return
	}
	key := jwt.CheckToken(token)
	if key == nil || key.User == nil {
		c.TokenError("登陆超时,请重新登陆")
		c.Abort()
		return
	}
	err := service.Auth.CheckLoginToken(key.User.LoginToken)
	if err != nil {
		c.TokenError(err.Error())
		c.Abort()
		return
	}
	c.SetUserInfo(key.User)
	c.Next()
}

// CheckToken 判断token
func CheckToken(c *context.Context) {
	err := service.Auth.CheckApiToken(c.Request.Header.Get(constant.JwtTokenName))
	if err != nil {
		c.TokenError(err.Error())
		c.Abort()
		return
	}
	c.Next()
}

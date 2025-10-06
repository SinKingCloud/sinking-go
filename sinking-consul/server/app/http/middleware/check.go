package middleware

import (
	"server/app/constant"
	"server/app/util"
	"server/app/util/context"
	"server/app/util/jwt"
)

// CheckLogin 判断登录
func CheckLogin(c *context.Context) {
	token := c.Request.Header.Get(constant.JwtTokenName)
	if token == "" {
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
	loginToken := util.Conf.GetString(constant.AuthLoginToken)
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

// CheckToken 判断token
func CheckToken(c *context.Context) {
	token := c.Request.Header.Get(constant.JwtTokenName)
	if token == "" {
		c.NotLogin("鉴权token缺失,请检查请求", nil)
		c.Abort()
		return
	}
	if token != util.Conf.GetString(constant.AuthApiToken) {
		c.TokenError("鉴权token无效,请确认token是否正确", nil)
		c.Abort()
		return
	}
	c.Next()
}

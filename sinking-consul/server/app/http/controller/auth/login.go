package auth

import (
	"server/app/service"
	"server/app/service/log"
	"server/app/util/server"
)

// Login 账号登录
func Login(c *server.Context) {
	type Form struct {
		Account  string `json:"account" default:"" validate:"required" label:"账户"`
		Password string `json:"password" default:"" validate:"required" label:"密码"`
		Token    string `json:"token" default:"" validate:"required" label:"验证码标识"`
		CaptchaX int    `json:"captcha_x" default:"" validate:"required,numeric" label:"验证码X坐标"`
		CaptchaY int    `json:"captcha_y" default:"" validate:"required,numeric" label:"验证码Y坐标"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if !service.Auth.CheckCaptcha(form.Token, form.CaptchaX, form.CaptchaY) {
		c.Error("验证码验证失败请重试")
		return
	}
	if e := service.Auth.CheckAccount(form.Account, form.Password); e != nil {
		c.Error(e.Error())
		return
	}
	token, e := service.Auth.GenLoginToken(c.GetRequestIp())
	if e != nil {
		c.Error(e.Error())
		return
	}
	service.Log.Create(c.GetRequestIp(), log.EventLogin, "系统账户登录", "系统登录成功")
	c.SuccessWithData("登录成功", token)
}

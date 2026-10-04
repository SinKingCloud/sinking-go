package auth

import (
	"server/app/service"
	"server/app/util/context"
)

// Captcha 获取验证吗
func Captcha(c *context.Context) {
	var form struct {
		Token string `json:"token" default:"" validate:"required,len=16" label:"验证码标识"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	captcha, err := service.Auth.GetCaptcha(form.Token)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取验证码成功", captcha)
}

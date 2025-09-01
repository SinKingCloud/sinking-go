package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/constant"
	"server/app/util"
	"server/app/util/server"
	"server/app/util/str"
)

type ControllerPerson struct {
}

func (ControllerPerson) Info(c *server.Context) {
	user := c.GetUserInfo()
	c.SuccessWithData("获取成功", sinking_web.H{
		"account":    util.Conf.GetString(constant.AuthAccount),
		"login_ip":   user.LoginIp,
		"login_time": user.LoginTime,
	})
}

func (ControllerPerson) Password(c *server.Context) {
	type Form struct {
		Password string `json:"password" default:"" validate:"required" label:"密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	sPwd, _ := str.NewStringTool().BcryptHash(form.Password)
	util.Conf.Set(constant.AuthPassword, sPwd)
	if util.Conf.WriteConfig() != nil {
		c.Error("修改失败")
	} else {
		c.Success("修改成功")
	}
}

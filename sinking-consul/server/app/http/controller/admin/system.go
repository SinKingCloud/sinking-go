package admin

import (
	"server/app/constant"
	"server/app/util"
	"server/app/util/server"
	"server/app/util/str"
)

type ControllerSystem struct{}

func (ControllerSystem) Password(c *server.Context) {
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
		c.Error("修改密码失败")
	} else {
		c.Success("修改密码成功")
	}
}

func (ControllerSystem) Overview(c *server.Context) {
	c.Success("获取成功")
}

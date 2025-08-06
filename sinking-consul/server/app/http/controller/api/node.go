package api

import (
	"server/app/service"
	"server/app/util/server"
)

type ControllerNode struct {
}

// Register 注册服务
func (ControllerNode) Register(c *server.Context) {
	type Form struct {
		Group   string `json:"group" default:"" validate:"required" label:"服务组"`
		Name    string `json:"name" default:"" validate:"required" label:"服务名"`
		Address string `json:"address" default:"" validate:"required" label:"注册地址"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	service.Node.Register(form.Group, form.Name, form.Address)
	c.Success("注册成功")
}

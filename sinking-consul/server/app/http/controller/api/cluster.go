package api

import (
	"server/app/service"
	"server/app/util/server"
)

type ControllerCluster struct {
}

// Register 注册服务
func (ControllerCluster) Register(c *server.Context) {
	type Form struct {
		Address string `json:"address" default:"" validate:"required" label:"注册地址"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	service.Cluster.Register(form.Address)
	c.Success("注册成功")
}

// Sync 同步数据
func (ControllerCluster) Sync(c *server.Context) {
	c.Success("注册成功")
}

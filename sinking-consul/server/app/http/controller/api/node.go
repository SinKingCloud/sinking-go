package api

import (
	"server/app/model"
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

// Sync 同步服务
func (ControllerNode) Sync(c *server.Context) {
	type Form struct {
		Group        string `json:"group" default:"" validate:"required" label:"服务组"`
		LastSyncTime int64  `json:"last_sync_time" default:"" validate:"required" label:"上次同步时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	var list []*model.Node
	lastSyncTime := service.Node.GetOperateTime(form.Group)
	if form.LastSyncTime == 0 || lastSyncTime != form.LastSyncTime {
		list = service.Node.GetAllOnlineNodes(form.Group)
	}
	c.SuccessWithData("获取成功", list)
}

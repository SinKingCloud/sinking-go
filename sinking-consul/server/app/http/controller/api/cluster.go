package api

import (
	"server/app/model"
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

// Node 服务列表
func (ControllerCluster) Node(c *server.Context) {
	c.SuccessWithData("获取成功", service.Node.GetLocalNodes())
}

// Config 配置列表
func (ControllerCluster) Config(c *server.Context) {
	type Form struct {
		ShowContent bool `json:"show_content" default:"" validate:"omitempty" label:"是否返回内容"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	c.SuccessWithData("获取成功", service.Config.GetAllConfigs("*", form.ShowContent, false))
}

// Lock 分布式锁
func (ControllerCluster) Lock(c *server.Context) {
	type Form struct {
		Status int `json:"status" default:"" validate:"oneof=0 1" label:"锁状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if form.Status == 0 {
		err := service.Cluster.SyncDataLock()
		if err != nil {
			c.Error(err.Error())
			return
		}
		c.Success("上锁成功")
	} else {
		err := service.Cluster.SyncDataUnLock()
		if err != nil {
			c.Error(err.Error())
			return
		}
		c.Success("解锁成功")
	}
}

// Delete 删除数据
func (ControllerCluster) Delete(c *server.Context) {
	type Form struct {
		Configs []*model.Config `json:"configs" default:"" validate:"omitempty" label:"配置列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if form.Configs != nil && len(form.Configs) > 0 {
		_ = service.Config.DeleteByGroupAndName(form.Configs)
	}
	c.Success("执行成功")
}

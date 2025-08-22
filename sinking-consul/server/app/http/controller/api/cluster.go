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

// List 集群列表
func (ControllerCluster) List(c *server.Context) {
	c.SuccessWithData("获取成功", service.Cluster.GetAllClusters())
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
	c.SuccessWithData("获取成功", service.Config.GetAllConfigs("*", form.ShowContent))
}

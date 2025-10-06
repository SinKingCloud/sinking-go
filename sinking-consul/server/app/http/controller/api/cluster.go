package api

import (
	"server/app/model"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/util/context"
	"server/app/util/str"
	"time"
)

type ControllerCluster struct {
}

// Testing 节点测试
func (ControllerCluster) Testing(c *context.Context) {
	c.SuccessWithData("获取成功", time.Now().Unix())
}

// Register 注册服务
func (ControllerCluster) Register(c *context.Context) {
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
func (ControllerCluster) Node(c *context.Context) {
	c.SuccessWithData("获取成功", service.Node.GetLocalNodes())
}

// Config 配置列表
func (ControllerCluster) Config(c *context.Context) {
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
func (ControllerCluster) Lock(c *context.Context) {
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
func (ControllerCluster) Delete(c *context.Context) {
	type Form struct {
		Configs []*model.Config `json:"configs" default:"" validate:"omitempty" label:"配置列表"`
		Nodes   []*model.Node   `json:"nodes" default:"" validate:"omitempty" label:"节点列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if form.Configs != nil && len(form.Configs) > 0 {
		_ = service.Config.DeleteByGroupAndName(form.Configs)
	}
	if form.Nodes != nil && len(form.Nodes) > 0 {
		addresses := make([]string, 100)
		for _, node := range form.Nodes {
			if node != nil && node.Address != "" {
				addresses = append(addresses, node.Address)
			}
		}
		if addresses != nil && len(addresses) > 0 {
			_ = service.Node.DeleteByAddress(addresses)
		}
	}
	c.Success("执行成功")
}

// Create 创建数据
func (ControllerCluster) Create(c *context.Context) {
	type Form struct {
		Config *model.Config `json:"config" default:"" validate:"omitempty" label:"配置信息"`
		Node   *model.Node   `json:"node" default:"" validate:"omitempty" label:"节点信息"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if form.Config != nil {
		_ = service.Config.Create(form.Config)
	}
	if form.Node != nil {
		_ = service.Node.Create(form.Node)
	}
	c.Success("执行成功")
}

// Update 更新数据
func (ControllerCluster) Update(c *context.Context) {
	type Form struct {
		Configs *cluster.ConfigUpdateValidate `json:"configs" default:"" validate:"omitempty" label:"配置列表"`
		Nodes   *cluster.NodeUpdateValidate   `json:"nodes" default:"" validate:"omitempty" label:"节点列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if form.Configs != nil && len(form.Configs.Keys) > 0 {
		data := make(map[string]interface{})
		if form.Configs.Type != "" {
			data["type"] = form.Configs.Type
		}
		if form.Configs.Content != "" {
			data["content"] = form.Configs.Content
			data["hash"] = str.NewStringTool().Md5(form.Configs.Content)
		}
		if form.Configs.Status != "" {
			data["status"] = form.Configs.Status
		}
		_ = service.Config.UpdateByGroupAndName(form.Configs.Keys, data)
	}
	if form.Nodes != nil && len(form.Nodes.Addresses) > 0 {
		data := make(map[string]interface{})
		if form.Nodes.Status != "" {
			data["status"] = form.Nodes.Status
		}
		_ = service.Node.UpdateByAddresses(form.Nodes.Addresses, data)
	}
	c.Success("执行成功")
}

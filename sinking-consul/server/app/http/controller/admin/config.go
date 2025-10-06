package admin

import (
	"server/app/model"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/log"
	"server/app/util/context"
	"server/app/util/page"
	"server/app/util/str"
	"strconv"
)

type ControllerConfig struct {
}

func (ControllerConfig) List(c *context.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `json:"order_by_field" default:"create_time" validate:"oneof=group name update_time create_time" label:"排序字段"`
		OrderByType     string `json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Group           string `json:"group" default:"" validate:"omitempty" label:"配置分组"`
		Name            string `json:"name" default:"" validate:"omitempty" label:"配置名称"`
		Type            string `json:"type" default:"" validate:"omitempty" label:"配置类型"`
		Hash            string `json:"hash" default:"" validate:"omitempty" label:"配置hash"`
		Content         string `json:"content" default:"" validate:"omitempty" label:"配置内容"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Group != "" {
		where["group"] = form.Group
	}
	if form.Name != "" {
		where["name"] = form.Name
	}
	if form.Type != "" {
		where["type"] = form.Type
	}
	if form.Hash != "" {
		where["hash"] = form.Hash
	}
	if form.Content != "" {
		where["content"] = form.Content
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	data, total, err := service.Config.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取失败")
	} else {
		service.Log.Create(c.GetRequestIp(), log.EventShow, "查看服务配置", "查看服务配置列表")
		c.SuccessWithData("获取成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

func (ControllerConfig) Info(c *context.Context) {
	type Form struct {
		Group string `json:"group" default:"" validate:"required" label:"配置分组"`
		Name  string `json:"name" default:"" validate:"required" label:"配置名称"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Config.FindByGroupAndName(form.Group, form.Name)
	if err != nil {
		c.Error("获取失败")
	} else {
		service.Log.Create(c.GetRequestIp(), log.EventShow, "查看配置详情", "查看配置["+form.Group+":"+form.Name+"]详情")
		c.SuccessWithData("获取成功", info)
	}
}

func (ControllerConfig) Update(c *context.Context) {
	form := &cluster.ConfigUpdateValidate{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Type != "" {
		data["type"] = form.Type
	}
	if form.Content != "" {
		data["content"] = form.Content
		data["hash"] = str.NewStringTool().Md5(form.Content)
	}
	if form.Status != "" {
		n, _ := strconv.Atoi(form.Status)
		if _, ok := service.Config.Status()[config.Status(n)]; !ok {
			c.Error("状态值不合法")
			return
		}
		data["status"] = form.Status
	}
	err := service.Cluster.ChangeAllClusterLockStatus(0)
	if err != nil {
		c.Error("获取分布式锁失败")
		return
	}
	defer func() {
		_ = service.Cluster.ChangeAllClusterLockStatus(1)
	}()
	err = service.Config.UpdateByGroupAndName(form.Keys, data)
	if err != nil {
		c.Success("修改失败")
		return
	}
	service.Cluster.UpdateAllClusterData(form, nil)
	service.Log.Create(c.GetRequestIp(), log.EventUpdate, "修改服务配置", "修改服务配置数据")
	c.Success("修改成功")
}

func (ControllerConfig) Create(c *context.Context) {
	type Form struct {
		Group   string `json:"group" default:"" validate:"required" label:"配置分组"`
		Name    string `json:"name" default:"" validate:"required" label:"配置名称"`
		Type    string `json:"type" default:"" validate:"required" label:"配置类型"`
		Content string `json:"content" default:"" validate:"omitempty" label:"配置内容"`
		Status  int    `json:"status" default:"" validate:"numeric" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if _, ok := service.Config.Status()[config.Status(form.Status)]; !ok {
		c.Error("状态值不合法")
		return
	}
	if _, ok := service.Config.Types()[config.Type(form.Type)]; !ok {
		c.Error("配置类型不合法")
		return
	}
	d := &model.Config{
		Group:  form.Group,
		Name:   form.Name,
		Type:   form.Type,
		Status: form.Status,
	}
	if form.Content != "" {
		d.Hash = str.NewStringTool().Md5(form.Content)
		d.Content = form.Content
	}
	err := service.Cluster.ChangeAllClusterLockStatus(0)
	if err != nil {
		c.Error("获取分布式锁失败")
		return
	}
	defer func() {
		_ = service.Cluster.ChangeAllClusterLockStatus(1)
	}()
	err = service.Config.Create(d)
	if err != nil {
		c.Error("创建失败")
		return
	}
	service.Cluster.CreateAllClusterData(d, nil)
	service.Log.Create(c.GetRequestIp(), log.EventCreate, "创建服务配置", "创建服务配置["+form.Group+":"+form.Name+"]")
	c.Success("创建成功")
}

func (ControllerConfig) Delete(c *context.Context) {
	type Form struct {
		Keys []*model.Config `json:"keys" default:"" validate:"required,min=1,max=1000" label:"配置列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	err := service.Cluster.ChangeAllClusterLockStatus(0)
	if err != nil {
		c.Error("获取分布式锁失败")
		return
	}
	defer func() {
		_ = service.Cluster.ChangeAllClusterLockStatus(1)
	}()
	err = service.Config.DeleteByGroupAndName(form.Keys)
	if err != nil {
		c.Error("删除失败")
		return
	}
	service.Cluster.DeleteAllClusterData(form.Keys, nil)
	service.Log.Create(c.GetRequestIp(), log.EventDelete, "删除服务配置", "删除服务配置数据")
	c.Success("删除成功")
}

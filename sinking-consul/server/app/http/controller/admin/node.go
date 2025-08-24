package admin

import (
	"server/app/service"
	"server/app/service/node"
	"server/app/util/page"
	"server/app/util/server"
	"strconv"
)

type ControllerNode struct {
}

func (ControllerNode) List(c *server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `json:"order_by_field" default:"id" validate:"oneof=group name update_time create_time" label:"排序字段"`
		OrderByType     string `json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Group           string `json:"group" default:"" validate:"omitempty" label:"服务分组"`
		Name            string `json:"name" default:"" validate:"omitempty" label:"服务名称"`
		OnlineStatus    string `json:"online_status" default:"" validate:"omitempty,numeric" label:"在线状态"`
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
	if form.OnlineStatus != "" {
		where["online_status"] = form.OnlineStatus
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
	data, total, err := service.Node.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取失败")
	} else {
		c.SuccessWithData("获取成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

func (ControllerNode) Update(c *server.Context) {
	type Form struct {
		Addresses []string `json:"addresses" default:"" validate:"required,min=1,max=1000,unique" label:"节点列表"`
		Status    string   `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Status != "" {
		n, _ := strconv.Atoi(form.Status)
		if _, ok := service.Node.Status()[node.Status(n)]; !ok {
			c.Error("状态值不合法")
			return
		}
		data["status"] = form.Status
	}
	err := service.Node.UpdateByAddresses(form.Addresses, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("修改成功")
}

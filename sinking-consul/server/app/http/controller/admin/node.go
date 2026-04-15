package admin

import (
	"server/app/enum/log_type"
	"server/app/enum/node_status"
	"server/app/model"
	"server/app/repository/node"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/util/context"
	"server/app/util/page"
	"strconv"
)

type ControllerNode struct {
}

func (ControllerNode) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("create_time", "desc", "group,name,update_time,create_time")
	type Form struct {
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
	where := &node.SelectNode{}
	if form.Group != "" {
		where.Group = form.Group
	}
	if form.Name != "" {
		where.Name = form.Name
	}
	if form.OnlineStatus != "" {
		where.OnlineStatus = form.OnlineStatus
	}
	if form.Status != "" {
		where.Status = form.Status
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = form.UpdateTimeEnd
	}
	data, total, err := service.Node.Select(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error("获取失败")
	} else {
		service.Log.Create(c.GetRequestIp(), log_type.EventShow, "查看服务节点", "查看服务节点列表")
		c.SuccessWithData("获取成功", page.NewPage(total, pageNum, pageSize, data))
	}
}

func (ControllerNode) Update(c *context.Context) {
	form := &cluster.NodeUpdateValidate{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	data := &node.UpdateNode{}
	if form.Status != "" {
		n, _ := strconv.Atoi(form.Status)
		if _, ok := node_status.Map()[n]; !ok {
			c.Error("状态值不合法")
			return
		}
		data.Status = form.Status
	}
	err := service.Cluster.ChangeAllClusterLockStatus(0)
	if err != nil {
		c.Error("获取分布式锁失败")
		return
	}
	defer func() {
		_ = service.Cluster.ChangeAllClusterLockStatus(1)
	}()
	err = service.Node.UpdateByAddresses(form.Addresses, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	service.Cluster.UpdateAllClusterData(nil, form)
	service.Log.Create(c.GetRequestIp(), log_type.EventUpdate, "修改服务节点", "修改服务节点数据")
	c.Success("修改成功")
}

func (ControllerNode) Delete(c *context.Context) {
	type Form struct {
		Addresses []string `json:"addresses" default:"" validate:"required,min=1,max=1000,unique" label:"节点列表"`
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
	list, err2 := service.Node.SelectInAddress(form.Addresses)
	err = service.Node.DeleteByAddress(form.Addresses)
	if err != nil {
		c.Error("删除失败")
		return
	}
	if err2 == nil && list != nil {
		var list2 []*model.Node
		for _, v := range list {
			list2 = append(list2, &model.Node{
				Group:   v.Group,
				Name:    v.Name,
				Address: v.Address,
			})
		}
		service.Cluster.DeleteAllClusterData(nil, list2)
	}
	service.Log.Create(c.GetRequestIp(), log_type.EventDelete, "删除服务节点", "删除服务节点数据")
	c.Success("删除成功")
}

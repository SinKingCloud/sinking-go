package admin

import (
	"strconv"

	"server/app/enum/log_type"
	"server/app/repository/cluster"
	"server/app/service"
	"server/app/util/context"
)

type ControllerCluster struct {
}

func (ControllerCluster) List(c *context.Context) {
	query := c.ValidatePage("create_time", "desc", "address,create_time,update_time", "create_time,update_time")
	var form struct {
		Keyword         string `json:"keyword" default:"" validate:"omitempty,max=200" label:"关键词"`
		Address         string `json:"address" default:"" validate:"omitempty,max=200" label:"集群地址"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"在线状态"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &cluster.SelectCluster{}
	if form.Keyword != "" {
		where.Keyword = &form.Keyword
	}
	if form.Address != "" {
		where.Address = &form.Address
	}
	if form.Status != "" {
		status, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("在线状态参数错误")
			return
		}
		where.Status = &status
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = &form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = &form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = &form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = &form.UpdateTimeEnd
	}
	data, err := service.Cluster.Select(where, query)
	if err != nil {
		c.Error("获取失败")
	} else {
		service.Log.Create(c.GetRequestIp(), log_type.EventShow, "查看系统集群", "查看系统集群列表")
		c.SuccessWithData("获取成功", data)
	}
}

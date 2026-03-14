package api

import (
	"server/app/service"
	"server/app/service/config"
	"server/app/util/context"

	sinking_web "github.com/SinKingCloud/sinking-go/sinking-web"
)

type ControllerConfig struct {
}

// Sync 同步配置
func (ControllerConfig) Sync(c *context.Context) {
	type Form struct {
		Group        string `json:"group" default:"" validate:"required" label:"服务组"`
		LastSyncTime int64  `json:"last_sync_time" default:"" validate:"" label:"上次同步时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	var list []*config.Config
	lastOperateTime := service.Config.GetOperateTime(form.Group)
	if form.LastSyncTime <= 0 || form.LastSyncTime < lastOperateTime {
		list = service.Config.GetAllConfigs(form.Group, true, true)
	}
	c.SuccessWithData("获取成功", sinking_web.H{
		"last_operate_time": lastOperateTime,
		"list":              list,
	})
}

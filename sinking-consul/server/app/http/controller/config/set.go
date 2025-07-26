package config

import (
	"server/app/service"
	"server/app/util/server"
)

// Set 修改配置
func Set(c *server.Context) {
	type Config struct {
		Key   string `json:"key" default:"" validate:"required,max=100" label:"配置标识"`
		Value string `json:"value" default:"" validate:"omitempty" label:"配置内容"`
	}
	type Form struct {
		Configs []*Config `json:"configs" default:"" validate:"gte=1" label:"配置标识"`
		Group   string    `json:"group" default:"" validate:"required,max=100" label:"配置组"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	for _, v := range form.Configs {
		if v.Value != "" {
			_ = service.Config.Set(form.Group, v.Key, v.Value)
		}
	}
	c.Success("修改数据成功")
}

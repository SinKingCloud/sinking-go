package system

import (
	"server/app/service"
	"server/app/util/server"
)

func Enum(c *server.Context) {
	type Form struct {
		Name string `json:"name" default:"" validate:"required,name" label:"枚举名称"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if value, ok := service.Enum[form.Name]; ok {
		c.SuccessWithData("获取数据成功", value)
	} else {
		c.Error("枚举类型不存在")
	}
}

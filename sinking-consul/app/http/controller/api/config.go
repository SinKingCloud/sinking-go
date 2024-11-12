package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ConfigList 获取配置
func ConfigList(s *sinking_web.Context) {
	type register struct {
		AppName string `form:"app_name" json:"app_name"` //所属应用
		EnvName string `form:"env_name" json:"env_name"` //环境标识
	}
	form := &register{}
	err := s.BindAll(form)
	if err != nil || form.AppName == "" || form.EnvName == "" {
		response.Error(s, "参数不足", nil)
		return
	}
	list := (&model.Config{AppName: form.AppName, EnvName: form.EnvName}).SelectByNameCache()
	response.Success(s, "获取配置列表成功", list)
}

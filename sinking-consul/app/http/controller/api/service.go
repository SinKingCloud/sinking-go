package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ServiceRegister 注册服务
func ServiceRegister(s *sinking_web.Context) {
	type register struct {
		Name      string `form:"name" json:"name"`             //服务名称
		AppName   string `form:"app_name" json:"app_name"`     //所属应用
		EnvName   string `form:"env_name" json:"env_name"`     //环境标识
		GroupName string `form:"group_name" json:"group_name"` //分组名称
		Addr      string `form:"addr" json:"addr"`             //服务地址(规则ip:port)
	}
	form := &register{}
	err := s.BindJson(&form)
	if err != nil || form.Name == "" || form.AppName == "" || form.EnvName == "" || form.GroupName == "" || form.Addr == "" {
		response.Error(s, "参数不足", nil)
		return
	}
	response.Success(s, "注册服务成功", nil)
}

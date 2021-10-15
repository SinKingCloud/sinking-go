package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"time"
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
	service.RegisterService(form.Name, form.AppName, form.EnvName, form.GroupName, form.Addr, time.Now().Unix(), 0)
	service.RegisterLocalService(form.Name, form.AppName, form.EnvName, form.GroupName, form.Addr, time.Now().Unix(), 0)
	response.Success(s, "注册服务成功", nil)
}

// ServiceStatus 更改服务状态
func ServiceStatus(s *sinking_web.Context) {
	type register struct {
		Name        string `form:"name" json:"name"`                 //服务名称
		AppName     string `form:"app_name" json:"app_name"`         //所属应用
		EnvName     string `form:"env_name" json:"env_name"`         //环境标识
		GroupName   string `form:"group_name" json:"group_name"`     //分组名称
		ServiceHash string `form:"service_hash" json:"service_hash"` //服务hash
		Status      int    `form:"addr" json:"addr"`                 //服务状态
	}
	form := &register{}
	err := s.BindJson(&form)
	if err != nil || form.ServiceHash == "" {
		response.Error(s, "参数不足", nil)
		return
	}
	if service.ChangeServiceStatus(form.Name, form.AppName, form.EnvName, form.GroupName, form.ServiceHash, form.Status) &&
		service.ChangeLocalServiceStatus(form.Name, form.AppName, form.EnvName, form.GroupName, form.ServiceHash, form.Status) {
		response.Success(s, "服务状态更改成功", nil)
	} else {
		response.Error(s, "服务状态更改失败", nil)
	}
}

// ServiceList 获取服务列表
func ServiceList(s *sinking_web.Context) {
	type register struct {
		Name      string `form:"name" json:"name"`             //服务名称
		AppName   string `form:"app_name" json:"app_name"`     //所属应用
		EnvName   string `form:"env_name" json:"env_name"`     //环境标识
		GroupName string `form:"group_name" json:"group_name"` //分组名称
	}
	form := &register{}
	err := s.BindJson(&form)
	if err != nil || form.Name == "" || form.AppName == "" || form.EnvName == "" || form.GroupName == "" {
		response.Error(s, "参数不足", nil)
		return
	}
	response.Success(s, "获取服务列表成功", service.GetServiceList(form.AppName, form.EnvName, form.GroupName, form.Name))
}

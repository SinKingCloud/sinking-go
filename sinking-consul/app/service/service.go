package service

import "github.com/SinKingCloud/sinking-go/sinking-consul/app/model"

// Services 服务列表
var Services = make(map[string]map[string]*Service)

// Service 服务列表
type Service struct {
	Name          string         `json:"name"`            //服务名称
	AppName       string         `json:"app_name"`        //所属应用
	EnvName       string         `json:"env_name"`        //环境标识
	GroupName     string         `json:"group_name"`      //分组名称
	AppHash       string         `json:"app_hash"`        //标识hash(规则md5(AppName-EnvName-GroupName))
	Ip            string         `json:"ip"`              //集群ip
	Port          string         `json:"port"`            //集群端口
	ServiceHash   string         `json:"service_hash"`    //标识hash(规则md5(AppName-EnvName-GroupName-Ip:Port))
	LastHeartTime model.DateTime `json:"last_heart_time"` //上次心跳时间
	Status        int            `json:"status"`          //服务状态(0:正常/1:异常)
}

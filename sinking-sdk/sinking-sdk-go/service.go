package sinking_sdk_go

import (
	"strings"
	"time"
)

// services 服务列表 AppName.EnvName.GroupName.Name.ServiceHash
var services = make(map[string]map[string]map[string]map[string]map[string]*Service)

// Service 服务列表
type Service struct {
	Name          string `json:"name"`
	AppName       string `json:"app_name"`
	EnvName       string `json:"env_name"`
	GroupName     string `json:"group_name"`
	Addr          string `json:"addr"`
	ServiceHash   string `json:"service_hash"`
	LastHeartTime int64  `json:"last_heart_time"`
	Status        int    `json:"status"`
}

// getServices 获取节点
func (r *Register) getServices() {
	//设置注册节点
	go func() {
		for {
			servers := strings.Split(r.Servers, ",")
			for _, v := range servers {
				test := &RequestServer{
					Server:    v,
					TokenName: r.TokenName,
					Token:     r.Token,
				}
				list := test.getServerList()
				for _, v2 := range list.Data {
					if v2.Status == 1 {
						continue
					}
					if services[v2.AppName] == nil {
						services[v2.AppName] = map[string]map[string]map[string]map[string]*Service{}
					}
					if services[v2.AppName][v2.EnvName] == nil {
						services[v2.AppName][v2.EnvName] = map[string]map[string]map[string]*Service{}
					}
					if services[v2.AppName][v2.EnvName][v2.GroupName] == nil {
						services[v2.AppName][v2.EnvName][v2.GroupName] = map[string]map[string]*Service{}
					}
					if services[v2.AppName][v2.EnvName][v2.GroupName][v2.Name] == nil {
						services[v2.AppName][v2.EnvName][v2.GroupName][v2.Name] = map[string]*Service{}
					}
					services[v2.AppName][v2.EnvName][v2.GroupName][v2.Name][v2.ServiceHash] = v2
				}
			}
			time.Sleep(time.Duration(checkTime) * time.Second)
		}
	}()
}
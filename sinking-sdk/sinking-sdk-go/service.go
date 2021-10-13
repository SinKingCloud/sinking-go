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
			servicesTemp := make(map[string]map[string]map[string]map[string]map[string]*Service)
			for _, v := range servers {
				test := &RequestServer{
					Server:    v,
					TokenName: r.TokenName,
					Token:     r.Token,
				}
				list := test.getServerList()
				if list.Code != 200 {
					continue
				}
				for _, v2 := range list.Data {
					if v2.Status == 1 {
						continue
					}
					if servicesTemp[v2.AppName] == nil {
						servicesTemp[v2.AppName] = map[string]map[string]map[string]map[string]*Service{}
					}
					if servicesTemp[v2.AppName][v2.EnvName] == nil {
						servicesTemp[v2.AppName][v2.EnvName] = map[string]map[string]map[string]*Service{}
					}
					if servicesTemp[v2.AppName][v2.EnvName][v2.GroupName] == nil {
						servicesTemp[v2.AppName][v2.EnvName][v2.GroupName] = map[string]map[string]*Service{}
					}
					if servicesTemp[v2.AppName][v2.EnvName][v2.GroupName][v2.Name] == nil {
						servicesTemp[v2.AppName][v2.EnvName][v2.GroupName][v2.Name] = map[string]*Service{}
					}
					servicesTemp[v2.AppName][v2.EnvName][v2.GroupName][v2.Name][v2.ServiceHash] = v2
				}
			}
			services = servicesTemp
			time.Sleep(time.Duration(checkTime) * time.Second)
		}
	}()
}

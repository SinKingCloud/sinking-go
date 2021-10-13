package sinking_sdk_go

import (
	"strings"
	"sync"
	"time"
)

var (
	// services 服务列表 AppName.EnvName.GroupName.Name.ServiceHash
	services     = make(map[string]map[string]map[string]map[string]map[string]*Service)
	servicesLock sync.Mutex
	// serviceIndex 轮询获取服务地址下标
	serviceIndex = make(map[string]int)
)

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

// GetService 获取随机节点(负载均衡)
func (r *Register) GetService(name string) (string, bool) {
	key := Md5Encode(r.AppName + r.EnvName + r.GroupName + name)
	addr := services[r.AppName][r.EnvName][r.GroupName][name]
	if addr == nil {
		return "", false
	}
	serviceIndex[key]++
	if serviceIndex[key] >= len(addr) {
		serviceIndex[key] = 0
	}
	i := -1
	for _, v := range addr {
		i++
		if i == serviceIndex[key] {
			return v.Addr, true
		}
	}
	return "", false
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
				if list == nil || list.Code != 200 {
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
			servicesLock.Lock()
			services = servicesTemp
			servicesLock.Unlock()
			time.Sleep(time.Duration(checkTime) * time.Second)
		}
	}()
}

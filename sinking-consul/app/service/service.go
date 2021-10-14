package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"sync"
	"time"
)

// Services 服务列表
var (
	Services     = make(map[string]*Service)
	ServicesLock sync.Mutex
)

// Service 服务列表
type Service struct {
	Name          string `json:"name"`            //服务名称
	AppName       string `json:"app_name"`        //所属应用
	EnvName       string `json:"env_name"`        //环境标识
	GroupName     string `json:"group_name"`      //分组名称
	Addr          string `json:"addr"`            //服务地址(规则ip:port)
	ServiceHash   string `json:"service_hash"`    //标识hash(规则md5(Addr))
	LastHeartTime int64  `json:"last_heart_time"` //上次心跳时间
	Status        int    `json:"status"`          //服务状态(0:正常/1:异常)
}

// RegisterService 服务注册
func RegisterService(name string, appName string, envName string, groupName string, addr string) {
	info := &Service{
		Name:          name,
		AppName:       appName,
		EnvName:       envName,
		GroupName:     groupName,
		Addr:          addr,
		ServiceHash:   encode.Md5Encode(appName + envName + groupName + addr),
		LastHeartTime: time.Now().Unix(),
		Status:        0,
	}
	ServicesLock.Lock()
	Services[info.ServiceHash] = info
	ServicesLock.Unlock()
}

// ChangeServiceStatus 更改服务状态
func ChangeServiceStatus(hash string, status int) bool {
	if Services[hash] == nil {
		return false
	}
	ServicesLock.Lock()
	Services[hash].Status = status
	ServicesLock.Unlock()
	return true
}

// GetServiceList 获取服务列表
func GetServiceList(name string, appName string, envName string, groupName string) []*Service {
	var list []*Service
	for _, v := range Services {
		if name != "" {
			if name != v.Name {
				continue
			}
		}
		if appName != "" {
			if appName != v.AppName {
				continue
			}
		}
		if envName != "" {
			if envName != v.EnvName {
				continue
			}
		}
		if groupName != "" {
			if groupName != v.GroupName {
				continue
			}
		}
		list = append(list, v)
	}
	return list
}

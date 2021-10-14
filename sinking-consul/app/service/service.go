package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"sync"
	"time"
)

// Services 服务列表 AppName.EnvName.GroupName.Name.ServiceHash
var (
	Services     = make(map[string]map[string]map[string]map[string]map[string]*Service)
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
	defer ServicesLock.Unlock()
	if Services[appName] == nil {
		Services[appName] = map[string]map[string]map[string]map[string]*Service{}
	}
	if Services[appName][envName] == nil {
		Services[appName][envName] = map[string]map[string]map[string]*Service{}
	}
	if Services[appName][envName][groupName] == nil {
		Services[appName][envName][groupName] = map[string]map[string]*Service{}
	}
	if Services[appName][envName][groupName][name] == nil {
		Services[appName][envName][groupName][name] = map[string]*Service{}
	}
	Services[appName][envName][groupName][name][info.ServiceHash] = info

}

// ChangeServiceStatus 更改服务状态
func ChangeServiceStatus(name string, appName string, envName string, groupName string, hash string, status int) bool {
	if Services[hash] == nil {
		return false
	}
	ServicesLock.Lock()
	defer ServicesLock.Unlock()
	Services[appName][envName][groupName][name][hash].Status = status

	return true
}

// GetAllServiceList 获取所有服务列表
func GetAllServiceList() []*Service {
	ServicesLock.Lock()
	defer ServicesLock.Unlock()
	var list []*Service
	for _, v := range Services {
		for _, v1 := range v {
			for _, v2 := range v1 {
				for _, v3 := range v2 {
					for _, v4 := range v3 {
						list = append(list, v4)
					}
				}
			}
		}
	}

	return list
}

// GetServiceList 获取服务列表
func GetServiceList(appName string, envName string) []*Service {
	ServicesLock.Lock()
	defer ServicesLock.Unlock()
	list := Services[appName][envName]
	var temp []*Service
	for _, v := range list {
		for _, v1 := range v {
			for _, v2 := range v1 {
				temp = append(temp, v2)
			}
		}
	}
	return temp
}

func CopyService() map[string]map[string]map[string]map[string]map[string]*Service {
	var temp = make(map[string]map[string]map[string]map[string]map[string]*Service)
	ServicesLock.Lock()
	defer ServicesLock.Unlock()
	for k, v := range Services {
		if temp[k] == nil {
			temp[k] = map[string]map[string]map[string]map[string]*Service{}
		}
		for k1, v1 := range v {
			if temp[k][k1] == nil {
				temp[k][k1] = map[string]map[string]map[string]*Service{}
			}
			for k2, v2 := range v1 {
				if temp[k][k1][k2] == nil {
					temp[k][k1][k2] = map[string]map[string]*Service{}
				}
				for k3, v3 := range v2 {
					if temp[k][k1][k2][k3] == nil {
						temp[k][k1][k2][k3] = map[string]*Service{}
					}
					for k4, v4 := range v3 {
						temp[k][k1][k2][k3][k4] = v4
					}
				}
			}
		}
	}
	return temp
}

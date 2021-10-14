package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/str"
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
	ServicesLock.Unlock()
}

// ChangeServiceStatus 更改服务状态
func ChangeServiceStatus(name string, appName string, envName string, groupName string, hash string, status int) bool {
	if Services[hash] == nil {
		return false
	}
	ServicesLock.Lock()
	Services[appName][envName][groupName][name][hash].Status = status
	ServicesLock.Unlock()
	return true
}

// GetAllServiceList 获取所有服务列表
func GetAllServiceList() []*Service {
	ServicesLock.Lock()
	serviceList := Services
	ServicesLock.Unlock()
	var list []*Service
	for _, v := range serviceList {
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
	data := str.DeepCopy(&Services).(map[string]map[string]map[string]map[string]map[string]*Service)
	ServicesLock.Unlock()
	list := data[appName][envName]
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

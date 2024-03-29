package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"sync"
)

// Services 服务列表 AppName.EnvName.GroupName.Name.ServiceHash
var (
	Services          = make(map[string]map[string]map[string]map[string]map[string]*Service)
	LocalServices     = make(map[string]map[string]map[string]map[string]map[string]*Service)
	ServicesLock      sync.Mutex
	LocalServicesLock sync.Mutex
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
func RegisterService(name string, appName string, envName string, groupName string, addr string, lastHeartTime int64, status int) {
	info := &Service{
		Name:          name,
		AppName:       appName,
		EnvName:       envName,
		GroupName:     groupName,
		Addr:          addr,
		ServiceHash:   encode.Md5Encode(appName + envName + groupName + addr),
		LastHeartTime: lastHeartTime,
		Status:        status,
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

// RegisterLocalService 本地服务注册
func RegisterLocalService(name string, appName string, envName string, groupName string, addr string, lastHeartTime int64, status int) {
	info := &Service{
		Name:          name,
		AppName:       appName,
		EnvName:       envName,
		GroupName:     groupName,
		Addr:          addr,
		ServiceHash:   encode.Md5Encode(appName + envName + groupName + addr),
		LastHeartTime: lastHeartTime,
		Status:        status,
	}
	LocalServicesLock.Lock()
	defer LocalServicesLock.Unlock()
	if LocalServices[appName] == nil {
		LocalServices[appName] = map[string]map[string]map[string]map[string]*Service{}
	}
	if LocalServices[appName][envName] == nil {
		LocalServices[appName][envName] = map[string]map[string]map[string]*Service{}
	}
	if LocalServices[appName][envName][groupName] == nil {
		LocalServices[appName][envName][groupName] = map[string]map[string]*Service{}
	}
	if LocalServices[appName][envName][groupName][name] == nil {
		LocalServices[appName][envName][groupName][name] = map[string]*Service{}
	}
	LocalServices[appName][envName][groupName][name][info.ServiceHash] = info
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

// ChangeLocalServiceStatus 更改本地服务状态
func ChangeLocalServiceStatus(name string, appName string, envName string, groupName string, hash string, status int) bool {
	if Services[hash] == nil {
		return false
	}
	LocalServicesLock.Lock()
	defer LocalServicesLock.Unlock()
	LocalServices[appName][envName][groupName][name][hash].Status = status
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

// GetAllLocalServiceList 获取本地所有服务列表
func GetAllLocalServiceList() []*Service {
	LocalServicesLock.Lock()
	defer LocalServicesLock.Unlock()
	var list []*Service
	for _, v := range LocalServices {
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
func GetServiceList(appName string, envName string, groupName string, name string) []*Service {
	ServicesLock.Lock()
	defer ServicesLock.Unlock()
	list := Services[appName][envName][groupName][name]
	var temp []*Service
	for _, v2 := range list {
		temp = append(temp, v2)
	}
	return temp
}

// GetProjectAllServiceList 获取所有项目服务名称列表
func GetProjectAllServiceList(appName string, envName string) map[string]map[string][]*Service {
	ServicesLock.Lock()
	defer ServicesLock.Unlock()
	list := Services[appName][envName]
	temp := make(map[string]map[string][]*Service)
	for k, v := range list {
		if temp[k] == nil {
			temp[k] = map[string][]*Service{}
		}
		for k1, v1 := range v {
			if temp[k][k1] == nil {
				temp[k][k1] = []*Service{}
			}
			for _, v2 := range v1 {
				temp[k][k1] = append(temp[k][k1], v2)
			}
		}
	}
	return temp
}

// GetLocalServiceList 获取本地服务列表
func GetLocalServiceList(appName string, envName string) []*Service {
	LocalServicesLock.Lock()
	defer LocalServicesLock.Unlock()
	list := LocalServices[appName][envName]
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

// CopyService 复制服务列表
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

// CopyLocalService 复制本地服务列表
func CopyLocalService() map[string]map[string]map[string]map[string]map[string]*Service {
	var temp = make(map[string]map[string]map[string]map[string]map[string]*Service)
	LocalServicesLock.Lock()
	defer LocalServicesLock.Unlock()
	for k, v := range LocalServices {
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

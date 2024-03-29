package sinking_sdk_go

import (
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	//serviceKeys 储存的key，顺序存放service
	serviceKeys               = make(map[string][]*Service)
	serviceKeysLock           sync.Mutex
	serviceIndex              = 0
	serviceKeysConsistent     = make(map[string]*Consistent)
	serviceKeysConsistentLock = sync.Mutex{}
)

const (
	Rand = iota //随机
	Poll        //轮询
	Only        //只请求第一个
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
func (r *Register) GetService(groupName string, name string, mode int) (*Service, bool) {
	key := Md5Encode(r.AppName + r.EnvName + groupName + name)
	serviceKeysLock.Lock()
	defer serviceKeysLock.Unlock()
	addr := serviceKeys[key]
	n := len(addr)
	if addr == nil || n <= 0 {
		return nil, false
	}
	if n == 1 || mode == Only {
		return addr[0], true
	}
	if mode == Rand {
		return addr[rand.Intn(n)], true
	} else if mode == Poll {
		serviceIndex++
		if serviceIndex >= n {
			serviceIndex = 0
			return addr[0], true
		} else {
			return addr[serviceIndex], true
		}
	} else {
		return addr[0], true
	}
}

// GetServiceAddrByHash 通过hash获取service
func (r *Register) GetServiceAddrByHash(groupName string, name string, hash string) (string, bool) {
	key := Md5Encode(r.AppName + r.EnvName + groupName + name)
	serviceKeysConsistentLock.Lock()
	defer serviceKeysConsistentLock.Unlock()
	cons := serviceKeysConsistent[key]
	if cons == nil {
		return "", false
	}
	addr, err := cons.Get(hash)
	if err != nil || addr == "" {
		return "", false
	}
	return addr, true
}

// getServices 获取并更新节点
func (r *Register) getServices(sync bool) {
	//设置注册节点
	fun := func() {
		serviceKeysTemp := make(map[string]map[string]*Service)
		test := &RequestServer{
			Server:    r.server,
			TokenName: r.TokenName,
			Token:     r.Token,
		}
		if len(r.useService) > 0 {
			for k, v := range r.useService {
				for _, v1 := range v {
					list := test.getServerList(r.AppName, r.EnvName, k, v1)
					if list != nil && list.Code == 200 {
						for _, v2 := range list.Data {
							if v2.Status == 1 {
								continue
							}
							key := Md5Encode(v2.AppName + v2.EnvName + v2.GroupName + v2.Name)
							if serviceKeysTemp[key] == nil {
								serviceKeysTemp[key] = map[string]*Service{}
							}
							serviceKeysTemp[key][v2.ServiceHash] = v2
						}
					}
				}
			}
		} else {
			list := test.getAllServerList(r.AppName, r.EnvName)
			if list != nil && list.Code == 200 {
				for _, v2 := range list.Data {
					for _, v3 := range v2 {
						for _, v4 := range v3 {
							if v4.Status == 1 {
								continue
							}
							key := Md5Encode(v4.AppName + v4.EnvName + v4.GroupName + v4.Name)
							if serviceKeysTemp[key] == nil {
								serviceKeysTemp[key] = map[string]*Service{}
							}
							serviceKeysTemp[key][v4.ServiceHash] = v4
						}
					}
				}
			}
		}
		serviceTemp := make(map[string][]*Service)
		serviceConsistentTemp := make(map[string]*Consistent)
		for i, v := range serviceKeysTemp {
			var keys []string
			for k := range v {
				keys = append(keys, k)
			}
			//按字典升序排列
			sort.Strings(keys)
			var temp []*Service
			temp2 := NewConsistent()
			for _, k := range keys {
				temp = append(temp, v[k])
				temp2.Add(v[k].Addr)
			}
			serviceTemp[i] = temp
			serviceConsistentTemp[i] = temp2
		}
		serviceKeysConsistentLock.Lock()
		serviceKeysConsistent = serviceConsistentTemp
		serviceKeysConsistentLock.Unlock()
		serviceKeysLock.Lock()
		serviceKeys = serviceTemp
		serviceKeysLock.Unlock()
	}
	if sync {
		go func() {
			for {
				fun()
				time.Sleep(time.Duration(checkTime) * time.Second)
			}
		}()
	} else {
		fun()
	}
}

// changeServerStatus 广播更改节点服务状态
func (r *Register) changeServerStatus(serviceHash string, status int) {
	//下线节点
	go func() {
		servers := strings.Split(r.Servers, ",")
		for _, v := range servers {
			test := &RequestServer{
				Server:    v,
				TokenName: r.TokenName,
				Token:     r.Token,
			}
			test.changeServerStatus(serviceHash, status)
		}
	}()
}

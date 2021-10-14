package sinking_sdk_go

import (
	"strings"
	"sync"
	"time"
)

var (
	OnlineStatus     = true
	OnlineStatusLock sync.Mutex
)

// registerServers 注册节点
func (r *Register) registerServices() {
	//设置注册节点
	go func() {
		for {
			if OnlineStatus {
				servers := strings.Split(r.Servers, ",")
				for _, v := range servers {
					test := &RequestServer{
						Server:    v,
						TokenName: r.TokenName,
						Token:     r.Token,
					}
					test.registerServer(r.Name, r.AppName, r.EnvName, r.GroupName, r.Addr)
				}
			}
			time.Sleep(time.Duration(checkTime) * time.Second)
		}
	}()
}

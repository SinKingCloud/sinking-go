package sinking_sdk_go

import (
	"sync"
	"time"
)

var (
	OnlineStatus     = true
	OnlineStatusLock sync.RWMutex
)

// registerServers 注册节点
func (r *Register) registerServices(sync bool) {
	//设置注册节点
	fun := func() {
		if OnlineStatus {
			test := &RequestServer{
				Server:    r.server,
				TokenName: r.TokenName,
				Token:     r.Token,
			}
			res := test.registerServer(r.Name, r.AppName, r.EnvName, r.GroupName, r.Addr)
			if res == nil || res.Code != 200 {
				r.changeServer(true)
			}
		}
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

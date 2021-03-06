package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/request"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"strconv"
	"sync"
	"time"
)

func registerCluster() {
	go func() {
		clusterList := service.CopyRegisterClusters()
		clusterListLock := sync.Mutex{}
		(&job.Task{Thread: len(clusterList), Producer: func(channel chan string) {
			for {
				for k := range clusterList {
					channel <- k
				}
				time.Sleep(time.Duration(setting.GetSystemConfig().Servers.HeartTime) * time.Second)
			}
		}, Consumer: func(k string) {
			clusterListLock.Lock()
			res := &request.Request{
				Ip:      clusterList[k].Ip,
				Port:    clusterList[k].Port,
				Timeout: 5,
			}
			clusterListLock.Unlock()
			res.Register(setting.GetSystemConfig().App.Ip, strconv.Itoa(setting.GetSystemConfig().App.Port))
		}}).Run()
	}()
}

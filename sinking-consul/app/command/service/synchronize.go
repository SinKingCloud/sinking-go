package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/request"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"time"
)

// synchronize 同步数据
func synchronize() {
	go func() {
		(&job.Task{Thread: service.RegisterClustersNum(), Producer: func(channel chan string) {
			for {
				service.RegisterClusters.Range(func(key, value any) bool {
					channel <- key.(string)
					return true
				})
				time.Sleep(time.Duration(setting.GetSystemConfig().Servers.HeartTime) * time.Second)
			}
		}, Consumer: func(k string) {
			value, ok := service.RegisterClusters.Load(k)
			if ok {
				cluster := value.(*service.Cluster)
				req := &request.Request{
					Ip:      cluster.Ip,
					Port:    cluster.Port,
					Timeout: 5,
				}
				syncService(req)
			}
		}}).Run()
	}()
}

// syncService 同步服务
func syncService(req *request.Request) {
	list, err := req.SetTimeout(15).ServiceList()
	if err != nil {
		return
	}
	for _, v := range list {
		service.RegisterService(v.Name, v.AppName, v.EnvName, v.GroupName, v.Addr, v.LastHeartTime, v.Status)
	}
}

package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/request"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"strconv"
	"time"
)

func registerCluster() {
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
				res := &request.Request{
					Ip:      cluster.Ip,
					Port:    cluster.Port,
					Timeout: 5,
				}
				res.Register(setting.GetSystemConfig().App.Ip, strconv.Itoa(setting.GetSystemConfig().App.Port))
			}
		}}).Run()
	}()
}

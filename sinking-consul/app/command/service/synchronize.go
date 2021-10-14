package service

import (
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"time"
)

func synchronize() {
	go func() {
		clusterList := service.CopyRegisterClusters()
		(&job.Task{Thread: len(clusterList), Producer: func(channel chan string) {
			for {
				for k := range clusterList {
					channel <- k
				}
				time.Sleep(time.Duration(setting.GetSystemConfig().Servers.HeartTime) * time.Second)
			}
		}, Consumer: func(k string) {
			//res := &request.Request{
			//	Ip:   clusterList[k].Ip,
			//	Port: clusterList[k].Port,
			//}
			//res.Register()
			fmt.Println("同步数据", k)
		}}).Run()
	}()
}

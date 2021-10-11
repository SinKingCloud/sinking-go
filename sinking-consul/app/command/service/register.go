package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/request"
	"time"
)

func register() {
	(&job.Task{
		Thread: len(service.RegisterClusters),
		Producer: func(channel chan string) {
			for {
				for k := range service.RegisterClusters {
					channel <- k
				}
				time.Sleep(5 * time.Second)
			}
		},
		Consumer: func(hash string) {
			info := service.RegisterClusters[hash]
			if info != nil {
				res := &request.Request{
					Ip:   info.Ip,
					Port: info.Port,
				}
				if !res.Register() {
					logs.Println(info.Hash, "注册失败,节点已下线")
				}
			}
		},
	}).Run()
}

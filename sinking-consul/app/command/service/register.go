package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/request"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
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
				time.Sleep(30 * time.Second)
			}
		},
		Consumer: func(hash string) {
			info := service.RegisterClusters[hash]
			if info != nil {
				res := &request.Request{
					Ip:        info.Ip,
					Port:      info.Port,
					TokenName: setting.GetConfig().GetString("servers."),
					Token:     setting.GetConfig().GetString(""),
				}
				res.Register()
			}

		},
	}).Run()
}

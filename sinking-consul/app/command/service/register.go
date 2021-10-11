package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
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
			logs.Println(hash)
		},
	}).Run()
}

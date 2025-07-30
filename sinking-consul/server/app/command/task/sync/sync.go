package sync

import (
	"server/app/command/queue/sync"
	"server/app/service"
	"server/app/service/cluster"
	"time"
)

func Init() {
	go func() {
		i := 0
		for {
			service.Cluster.Each(func(key string, value *cluster.Cluster) bool {
				if value.LastHeart+60 < time.Now().Unix() {
					value.OnlineStatus = int(cluster.Offline)
				}
				sync.Instance.SendTask(&sync.Task{
					Type:          sync.RegisterService,
					RemoteAddress: key,
				})
				return true
			})
			if i == 5 {
				service.Cluster.Each(func(key string, _ *cluster.Cluster) bool {
					sync.Instance.SendTask(&sync.Task{
						Type:          sync.SynchronizeData,
						RemoteAddress: key,
					})
					return true
				})
				i = 0
			}
			time.Sleep(10 * time.Second)
			i++
		}
	}()
}

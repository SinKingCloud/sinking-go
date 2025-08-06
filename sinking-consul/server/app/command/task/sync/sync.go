package sync

import (
	"server/app/command/queue/sync"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/node"
	"time"
)

// Init 定时同步数据
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
			service.Node.Each("*", func(value *node.Node) {
				if value.LastHeart+60 < time.Now().Unix() {
					value.OnlineStatus = int(node.Offline)
				}
			})
			if i == 3 {
				service.Cluster.Each(func(key string, value *cluster.Cluster) bool {
					if value.Status == int(cluster.Normal) && value.OnlineStatus == int(cluster.Online) {
						sync.Instance.SendTask(&sync.Task{
							Type:          sync.SynchronizeData,
							RemoteAddress: key,
						})
					}
					return true
				})
				i = 0
			}
			time.Sleep(10 * time.Second)
			i++
		}
	}()
}

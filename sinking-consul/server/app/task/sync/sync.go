package sync

import (
	"server/app/constant"
	"server/app/enum/cluster_status"
	"server/app/enum/node_online_status"
	"server/app/enum/node_status"
	"server/app/model"
	"server/app/queue"
	"server/app/queue/sync"
	"server/app/service"
	"server/app/service/node"
	"server/app/util"
	"time"
)

// Init 定时同步数据
func Init() {
	go func() {
		i := 0
		for {
			for {
				if util.Cache.IsLock(constant.LockSyncData) {
					time.Sleep(time.Second)
				} else {
					break
				}
			}
			service.Cluster.Each(func(key string, value *model.Cluster) bool {
				if value.LastHeart+60 < time.Now().Unix() {
					value.Status = cluster_status.Offline
				}
				queue.Sync.SendTask(&sync.Task{
					Type:          sync.RegisterService,
					RemoteAddress: key,
				})
				return true
			})
			temp1 := make(map[string]uint64)
			service.Node.Each("*", func(value *node.Node) {
				if value.Status == node_status.Normal && value.OnlineStatus == node_online_status.Online && value.LastHeart+60 < time.Now().Unix() {
					value.OnlineStatus = node_online_status.Offline
					temp1[value.Group] = 1
				}
			})
			for k := range temp1 {
				service.Node.SetOperateTime(k)
			}
			if i == 3 {
				service.Cluster.Each(func(key string, value *model.Cluster) bool {
					if value.Status == cluster_status.Online {
						queue.Sync.SendTask(&sync.Task{
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

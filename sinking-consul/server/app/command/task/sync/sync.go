package sync

import (
	"server/app/command/queue/sync"
	"server/app/constant"
	"server/app/enum/cluster_status"
	"server/app/enum/node_online_status"
	"server/app/enum/node_status"
	"server/app/model"
	"server/app/service"
	"server/app/service/node"
	"server/app/util"
	"server/app/util/str"
	"time"
)

// Init 定时同步数据
func Init() {
	go func() {
		strTool := str.NewStringTool()
		strToolMax := uint64(10000000)
		nodeTemp := make(map[string]uint64)
		configTemp := make(map[string]uint64)
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
				sync.Instance.SendTask(&sync.Task{
					Type:          sync.RegisterService,
					RemoteAddress: key,
				})
				return true
			})
			temp1 := make(map[string]uint64)
			service.Node.Each("*", func(value *node.Node) {
				if value.LastHeart+60 < time.Now().Unix() {
					value.OnlineStatus = node_online_status.Offline
				}
				if value.OnlineStatus == node_online_status.Online && value.Status == node_status.Normal {
					temp1[value.Group] += strTool.ToNumber(value.Address, strToolMax)
				}
			})
			for k, v := range temp1 {
				if v2, ok := nodeTemp[k]; !ok || v2 != v {
					service.Node.SetOperateTime(k)
					nodeTemp = temp1
				}
			}
			temp2 := make(map[string]uint64)
			service.Config.Each("*", func(value *model.Config) {
				temp2[value.Group] += strTool.ToNumber(value.Hash, strToolMax)
			})
			for k, v := range temp2 {
				if v2, ok := configTemp[k]; !ok || v2 != v {
					service.Config.SetOperateTime(k)
					configTemp = temp2
				}
			}
			if i == 3 {
				service.Cluster.Each(func(key string, value *model.Cluster) bool {
					if value.Status == cluster_status.Online {
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

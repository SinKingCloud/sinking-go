package sync

import (
	"server/app/command/queue/sync"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/node"
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
			service.Cluster.Each(func(key string, value *cluster.Cluster) bool {
				if value.LastHeart+60 < time.Now().Unix() {
					value.Status = int(cluster.Offline)
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
					value.OnlineStatus = int(node.Offline)
				}
				if value.OnlineStatus == int(node.Online) && value.Status == int(node.Normal) {
					temp1[value.Group] += strTool.ToNumber(value.Address, strToolMax)
				}
			})
			for k, v := range temp1 {
				if v2, ok := nodeTemp[k]; !ok || v2 != v {
					service.Node.SetOperateTime(k)
				}
			}
			temp2 := make(map[string]uint64)
			service.Config.Each("*", func(value *config.Config) {
				temp2[value.Group] += strTool.ToNumber(value.Hash, strToolMax)
			})
			for k, v := range temp2 {
				if v2, ok := configTemp[k]; !ok || v2 != v {
					service.Config.SetOperateTime(k)
				}
			}
			if i == 3 {
				service.Cluster.Each(func(key string, value *cluster.Cluster) bool {
					if value.Status == int(cluster.Online) {
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

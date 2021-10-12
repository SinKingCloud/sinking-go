package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"time"
)

func checkCluster() {
	go func() {
		for {
			//检测集群状态
			for k := range service.Clusters {
				if service.Clusters[k].LastHeartTime.Unix()+int64(setting.GetSystemConfig().Servers.CheckHeartTime) < time.Now().Unix() {
					service.Clusters[k].Status = 1
				}
			}
			//检测服务
			time.Sleep(time.Duration(setting.GetSystemConfig().Servers.HeartTime) * time.Second)
		}
	}()
}

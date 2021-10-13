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
				if service.Clusters[k].LastHeartTime+int64(setting.GetSystemConfig().Servers.CheckHeartTime) < time.Now().Unix() {
					service.ClustersLock.Lock()
					service.Clusters[k].Status = 1
					service.ClustersLock.Unlock()
				}
			}
			//检测服务状态
			for k := range service.Services {
				if service.Services[k].LastHeartTime+int64(setting.GetSystemConfig().Servers.CheckHeartTime) < time.Now().Unix() {
					service.ServicesLock.Lock()
					service.Services[k].Status = 1
					service.ServicesLock.Unlock()
				}
			}
			time.Sleep(time.Duration(setting.GetSystemConfig().Servers.HeartTime) * time.Second)
		}
	}()
}

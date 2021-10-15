package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/job"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/request"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"time"
)

// synchronize 同步数据
func synchronize() {
	go func() {
		clusterList := service.CopyRegisterClusters()
		(&job.Task{Thread: len(clusterList), Producer: func(channel chan string) {
			for {
				for k := range clusterList {
					channel <- k
				}
				time.Sleep(time.Duration(setting.GetSystemConfig().Servers.HeartTime) * time.Second)
			}
		}, Consumer: func(k string) {
			res := &request.Request{
				Ip:      clusterList[k].Ip,
				Port:    clusterList[k].Port,
				Timeout: 5,
			}
			syncService(res)
			syncConfig(res)
		}}).Run()
	}()
}

// syncService 同步服务
func syncService(req *request.Request) {
	list, err := req.SetTimeout(15).ServiceList()
	if err != nil {
		return
	}
	for _, v := range list {
		status := 0
		if v.LastHeartTime+int64(setting.GetSystemConfig().Servers.CheckHeartTime) < time.Now().Unix() {
			status = 1
		}
		service.RegisterService(v.Name, v.AppName, v.EnvName, v.GroupName, v.Addr, v.LastHeartTime, status)
	}
}

// syncConfig 同步配置
func syncConfig(req *request.Request) {
	list, err := req.SetTimeout(10).ConfigList()
	if err != nil {
		return
	}
	for _, v := range list {
		logs.Println(v)
	}
}

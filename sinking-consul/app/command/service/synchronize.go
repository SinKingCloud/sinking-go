package service

import (
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
)

func synchronize() {
	go func() {
		clusterList := service.CopyClusters()
		for k := range clusterList {
			//clusterList[k]
			fmt.Println(k) //同步节点配置及服务
		}
	}()
}

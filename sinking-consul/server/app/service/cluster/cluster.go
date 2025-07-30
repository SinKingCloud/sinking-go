package cluster

import (
	"crypto/tls"
	"net/http"
	"server/app/model"
	"sync"
	"time"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	//实例对象
	obj *Service
	//原子锁
	once sync.Once
	//原子锁
	clusterOnce sync.Once
	// 集群池
	clusterPool sync.Map
	// globalClient 全局请求client
	globalClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// Cluster 集群列表
type Cluster struct {
	*model.Cluster
}

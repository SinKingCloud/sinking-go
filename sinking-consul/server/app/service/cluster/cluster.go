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
	obj = &Service{}
	//原子锁
	clusterOnce = &sync.Once{}
	//正在同步数据的协程数量(原子计数器)
	syncDataCoroutineCount = int64(0)
	// 集群池
	clusterPool = &sync.Map{}
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
	return obj
}

// Cluster 集群列表
type Cluster struct {
	*model.Cluster
}

// ConfigUpdateValidate 配置更新验证器
type ConfigUpdateValidate struct {
	Keys    []*model.Config `json:"keys" default:"" validate:"required,min=1,max=1000" label:"配置列表"`
	Type    string          `json:"type" default:"" validate:"omitempty" label:"配置类型"`
	Content string          `json:"content" default:"" validate:"omitempty" label:"配置内容"`
	Status  string          `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
}

// NodeUpdateValidate 节点更新验证器
type NodeUpdateValidate struct {
	Addresses []string `json:"addresses" default:"" validate:"required,min=1,max=1000,unique" label:"节点列表"`
	Status    string   `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
}

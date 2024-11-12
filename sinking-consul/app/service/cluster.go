package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"sync"
	"time"
)

// Clusters 集群列表
var (
	//Clusters     = make(map[string]*Cluster)
	Clusters = &sync.Map{}
	//ClustersLock sync.Mutex
)

// RegisterClusters 需注册集群列表
var (
	//RegisterClusters     = make(map[string]*Cluster)
	RegisterClusters = &sync.Map{}
	//RegisterClustersLock sync.Mutex
)

// Cluster 集群信息结构
type Cluster struct {
	Hash          string `json:"hash"`            //标识hash(规则md5(ip:port))
	Ip            string `json:"ip"`              //集群ip
	Port          string `json:"port"`            //集群端口
	LastHeartTime int64  `json:"last_heart_time"` //上次心跳时间
	Status        int    `json:"status"`          //集群状态(0:正常/1:异常)
}

// ClustersRegister 集群注册
func ClustersRegister(ip string, port string) {
	info := &Cluster{
		Hash:          encode.Md5Encode(ip + ":" + port),
		Ip:            ip,
		Port:          port,
		LastHeartTime: time.Now().Unix(),
		Status:        0,
	}
	Clusters.Store(info.Hash, info)
}

// ClustersList 集群列表
func ClustersList() []*Cluster {
	var list []*Cluster
	Clusters.Range(func(key, value any) bool {
		list = append(list, value.(*Cluster))
		return true
	})
	return list
}

//// CopyClusters 复制节点数据
//func CopyClusters() map[string]*Cluster {
//	var temp = make(map[string]*Cluster)
//	for k, v := range Clusters {
//		temp[k] = v
//	}
//	return temp
//}
//
//// CopyRegisterClusters 复制节点数据
//func CopyRegisterClusters() map[string]*Cluster {
//	RegisterClustersLock.Lock()
//	defer RegisterClustersLock.Unlock()
//	var temp = make(map[string]*Cluster)
//	for k, v := range RegisterClusters {
//		temp[k] = v
//	}
//	return temp
//}

// RegisterClustersNum 节点数量
func RegisterClustersNum() int {
	num := 0
	RegisterClusters.Range(func(_, _ any) bool {
		num++
		return true
	})
	return num
}

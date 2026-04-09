package node

import "server/app/model"

// Node 集群列表
type Node struct {
	*model.Node
	IsLocal bool `gorm:"-" json:"is_local"`
}

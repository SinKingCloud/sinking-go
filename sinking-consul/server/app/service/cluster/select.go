package cluster

import (
	"server/app/model"
	repositoryCluster "server/app/repository/cluster"
)

// Select 获取数据
func (s *service) Select(where *repositoryCluster.SelectCluster, orderByField string, orderByType string, page int, pageSize int) (list []*model.Cluster, total int64, err error) {
	return s.repository.Select(where, orderByField, orderByType, page, pageSize)
}

package cluster

import (
	"errors"

	"server/app/enum/cluster_status"
	"server/app/model"
	repositoryCluster "server/app/repository/cluster"
)

// Select 获取数据
func (s *service) Select(where *repositoryCluster.SelectCluster, orderByField string, orderByType string, page int, pageSize int) (list []*model.Cluster, total int64, err error) {
	if where != nil && where.Status != nil {
		if _, ok := cluster_status.Map()[*where.Status]; !ok {
			return nil, 0, errors.New("集群状态参数不合法")
		}
	}
	return s.repository.Select(where, orderByField, orderByType, page, pageSize)
}

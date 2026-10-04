package cluster

import (
	"errors"

	"server/app/enum/cluster_status"
	"server/app/model"
	repositoryCluster "server/app/repository/cluster"
	"server/app/util/page"
)

// Select 获取数据
func (s *service) Select(where *repositoryCluster.SelectCluster, queryPage *page.Query) (*page.Result[*model.Cluster], error) {
	if where != nil && where.Status != nil {
		if _, ok := cluster_status.Map()[*where.Status]; !ok {
			return nil, errors.New("集群状态参数不合法")
		}
	}
	return s.repository.Select(where, queryPage)
}

package node

import (
	"errors"

	"server/app/enum/node_online_status"
	"server/app/enum/node_status"
	"server/app/model"
	repositoryNode "server/app/repository/node"
)

func (s *service) SelectInAddress(addresses []string) ([]*model.Node, error) {
	return s.repository.SelectInAddress(addresses)
}

// Select 获取数据
func (s *service) Select(where *repositoryNode.SelectNode, orderByField string, orderByType string, page int, pageSize int) (list []*model.Node, total int64, err error) {
	if where != nil {
		if where.Status != nil {
			if _, ok := node_status.Map()[*where.Status]; !ok {
				return nil, 0, errors.New("节点状态参数不合法")
			}
		}
		if where.OnlineStatus != nil {
			if _, ok := node_online_status.Map()[*where.OnlineStatus]; !ok {
				return nil, 0, errors.New("节点在线状态参数不合法")
			}
		}
	}
	return s.repository.Select(where, orderByField, orderByType, page, pageSize)
}

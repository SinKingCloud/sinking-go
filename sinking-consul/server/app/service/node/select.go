package node

import (
	"server/app/model"
	repositoryNode "server/app/repository/node"
)

func (s *service) SelectInAddress(addresses []string) ([]*model.Node, error) {
	return s.repository.SelectInAddress(addresses)
}

// Select 获取数据
func (s *service) Select(where *repositoryNode.SelectNode, orderByField string, orderByType string, page int, pageSize int) (list []*model.Node, total int64, err error) {
	return s.repository.Select(where, orderByField, orderByType, page, pageSize)
}

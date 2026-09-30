package node

import (
	"errors"

	"server/app/enum/node_online_status"
	"server/app/enum/node_status"
	repositoryNode "server/app/repository/node"
)

// UpdateAll 更新
func (s *service) UpdateAll(data *repositoryNode.UpdateNode) (err error) {
	if data == nil {
		return errors.New("更新数据不能为空")
	}
	if data.Status != nil {
		if _, ok := node_status.Map()[*data.Status]; !ok {
			return errors.New("节点状态不合法")
		}
	}
	if data.OnlineStatus != nil {
		if _, ok := node_online_status.Map()[*data.OnlineStatus]; !ok {
			return errors.New("节点在线状态不合法")
		}
	}
	return s.repository.UpdateAll(data)
}

// UpdateByAddresses 通过节点地址更新
func (s *service) UpdateByAddresses(addresses []string, data *repositoryNode.UpdateNode) (err error) {
	if data == nil {
		return errors.New("更新数据不能为空")
	}
	if data.Status != nil {
		if _, ok := node_status.Map()[*data.Status]; !ok {
			return errors.New("节点状态不合法")
		}
	}
	if data.OnlineStatus != nil {
		if _, ok := node_online_status.Map()[*data.OnlineStatus]; !ok {
			return errors.New("节点在线状态不合法")
		}
	}
	err = s.repository.UpdateByAddresses(addresses, data)
	list, err := s.repository.SelectInAddress(addresses)
	if err == nil {
		result := make([]*Node, 0, len(list))
		for _, item := range list {
			if item != nil {
				result = append(result, &Node{
					Node:    item,
					IsLocal: false,
				})
			}
		}
		s.Sets(result)
	}
	return err
}

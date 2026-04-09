package node

import repositoryNode "server/app/repository/node"

// UpdateAll 更新
func (s *service) UpdateAll(data *repositoryNode.UpdateNode) (err error) {
	return s.repository.UpdateAll(data)
}

// UpdateByAddresses 通过节点地址更新
func (s *service) UpdateByAddresses(addresses []string, data *repositoryNode.UpdateNode) (err error) {
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

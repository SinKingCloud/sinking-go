package node

import "server/app/model"

// Save 保存数据
func (s *service) Save() error {
	var nodes []*model.Node
	s.Each("*", func(value *Node) {
		if value != nil && value.Node != nil {
			nodes = append(nodes, value.Node)
		}
	})
	if len(nodes) == 0 {
		return nil
	}
	return s.repository.Save(nodes)
}

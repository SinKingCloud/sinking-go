package node

import (
	"errors"
	"server/app/constant"
	"server/app/enum/node_online_status"
	"server/app/enum/node_status"
	"server/app/model"
)

// Create 插入数据
func (s *service) Create(data *model.Node) (err error) {
	if data == nil {
		return errors.New("创建数据不能为空")
	}
	if _, ok := node_status.Map()[data.Status]; !ok {
		return errors.New("节点状态不合法")
	}
	if _, ok := node_online_status.Map()[data.OnlineStatus]; !ok {
		return errors.New("节点在线状态不合法")
	}
	key := constant.LockNodeCreate
	if !s.cache.Lock(key, constant.LockTimeNodeCreate) {
		return errors.New("系统繁忙,请稍后重试")
	}
	defer s.cache.UnLock(key)
	info, err := s.repository.FindByGroupAndAddress(data.Group, data.Address)
	if err == nil && info != nil {
		return errors.New("配置已存在")
	}
	err = s.repository.Create(data)
	if err == nil {
		info, err = s.repository.FindByGroupAndAddress(data.Group, data.Address)
		if err == nil && info != nil {
			s.Set(info.Group, info.Address, &Node{Node: info})
		}
	}
	return err
}

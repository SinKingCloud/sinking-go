package node

import (
	"errors"

	"server/app/enum/node_online_status"
)

// CountByOnlineStatus 统计online_status数量
func (s *service) CountByOnlineStatus(onlineStatus int) (total int64, err error) {
	if _, ok := node_online_status.Map()[onlineStatus]; !ok {
		return 0, errors.New("节点在线状态参数不合法")
	}
	return s.repository.CountByOnlineStatus(onlineStatus)
}

// CountAll 统计status数量
func (s *service) CountAll() (total int64, err error) {
	return s.repository.CountAll()
}

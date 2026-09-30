package cluster

import (
	"errors"

	"server/app/enum/cluster_status"
)

// CountByStatus 统计status数量
func (s *service) CountByStatus(status int) (total int64, err error) {
	if _, ok := cluster_status.Map()[status]; !ok {
		return 0, errors.New("集群状态参数不合法")
	}
	return s.repository.CountByStatus(status)
}

// CountAll 统计status数量
func (s *service) CountAll() (total int64, err error) {
	return s.repository.CountAll()
}

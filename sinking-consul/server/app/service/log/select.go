package log

import (
	"server/app/model"
	"server/app/repository/log"
)

// Select 获取数据
func (s *service) Select(where *log.SelectLog, orderByField string, orderByType string, page int, pageSize int) ([]*model.Log, int64, error) {
	return s.repositoryLog.Select(where, orderByField, orderByType, page, pageSize)
}

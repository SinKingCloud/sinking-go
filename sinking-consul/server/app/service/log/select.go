package log

import (
	"errors"

	"server/app/enum/log_type"
	"server/app/model"
	"server/app/repository/log"
)

// Select 获取数据
func (s *service) Select(where *log.SelectLog, orderByField string, orderByType string, page int, pageSize int) ([]*model.Log, int64, error) {
	if where != nil && where.Type != nil {
		if _, ok := log_type.Map()[*where.Type]; !ok {
			return nil, 0, errors.New("日志类型参数不合法")
		}
	}
	return s.repositoryLog.Select(where, orderByField, orderByType, page, pageSize)
}

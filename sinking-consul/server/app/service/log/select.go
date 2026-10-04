package log

import (
	"errors"

	"server/app/enum/log_type"
	"server/app/model"
	"server/app/repository/log"
	"server/app/util/page"
)

// Select 获取数据
func (s *service) Select(where *log.SelectLog, queryPage *page.Query) (*page.Result[*model.Log], error) {
	if where != nil && where.Type != nil {
		if _, ok := log_type.Map()[*where.Type]; !ok {
			return nil, errors.New("日志类型参数不合法")
		}
	}
	return s.repositoryLog.Select(where, queryPage)
}

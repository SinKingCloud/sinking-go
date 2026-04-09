package log

import (
	"server/app/model"
	"server/app/repository/log"
)

// Service 接口
type Service interface {
	Create(ip string, types int, title string, content string)                                                                 //插入数据
	Select(where *log.SelectLog, orderByField string, orderByType string, page int, pageSize int) ([]*model.Log, int64, error) //查询数据
}

// service 服务
type service struct {
	repositoryLog log.Interface
}

// NewService 实例化service
func NewService(repositoryLog log.Interface) *service {
	return &service{
		repositoryLog: repositoryLog,
	}
}

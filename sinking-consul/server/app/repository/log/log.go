package log

import (
	"server/app/model"
	"server/app/util/database"
)

// Interface 配置仓储接口
type Interface interface {
	Create(data *model.Log) (err error)                                                                                                   //插入数据
	Select(where *SelectLog, orderByField string, orderByType string, page int, pageSize int) (list []*model.Log, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}

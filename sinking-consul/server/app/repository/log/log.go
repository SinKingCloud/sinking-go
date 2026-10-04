package log

import (
	"server/app/model"
	"server/app/util/database"
	"server/app/util/page"
	"server/app/util/repository"
)

// Interface 配置仓储接口
type Interface interface {
	Create(data *model.Log) (err error)                                               //插入数据
	Select(where *SelectLog, queryPage *page.Query) (*page.Result[*model.Log], error) //查询数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*model.Log]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.Log](db)}
}

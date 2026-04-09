package config

import (
	"server/app/model"
	"server/app/util/database"
)

// Interface 配置仓储接口
type Interface interface {
	CountByStatus(status int) (total int64, err error)                                                                                    //根据status统计
	CountAll() (total int64, err error)                                                                                                   //统计全部
	Create(data *model.Config) (err error)                                                                                                //创建数据
	DeleteByGroupAndName(keys []*model.Config) error                                                                                      //通过集群和名称删除
	FindByGroupAndName(group string, name string) (*model.Config, error)                                                                  //通过group name查询父信息
	UpdateByGroupAndName(keys []*model.Config, config *UpdateConfig) error                                                                //根据 code 更新权限
	SelectAll() (list []*Config, err error)                                                                                               //查询所有
	SelectInGroupAndName(keys []*model.Config) (list []*Config, err error)                                                                //据group name查询
	Select(where *SelectConfig, orderByField string, orderByType string, page int, pageSize int) (list []*Config, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}

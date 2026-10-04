package cluster

import (
	"server/app/model"
	"server/app/util/database"
	"server/app/util/page"
	"server/app/util/repository"
)

// Interface 配置仓储接口
type Interface interface {
	Create(data *model.Cluster) (err error)                                                   //插入数据
	Save(clusters []*model.Cluster) error                                                     //保存数据
	CountByStatus(status int) (total int64, err error)                                        //根据状态统计
	CountAll() (total int64, err error)                                                       //统计所有
	DeleteAll() (err error)                                                                   //删除所有数据
	FindByAddress(address string) (data *model.Cluster, err error)                            //通过集群地址查询
	SelectAll() (list []*model.Cluster, err error)                                            //查询所有
	Select(where *SelectCluster, queryPage *page.Query) (*page.Result[*model.Cluster], error) //查询数据
}

// Repository 仓储
type Repository struct {
	*repository.Repository[*model.Cluster]
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Repository: repository.NewRepository[*model.Cluster](db)}
}

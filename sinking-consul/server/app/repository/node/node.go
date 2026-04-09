package node

import (
	"server/app/model"
	"server/app/util/database"
)

// Interface 配置仓储接口
type Interface interface {
	CountByOnlineStatus(onlineStatus int) (total int64, err error)                                                                          //根据在线状态统计
	CountAll() (total int64, err error)                                                                                                     //统计全部
	DeleteByAddress(addresses []string) error                                                                                               //通过address删除
	Create(data *model.Node) (err error)                                                                                                    //插入数据
	FindByGroupAndAddress(group string, address string) (data *model.Node, err error)                                                       //通过address查询信息
	UpdateAll(node *UpdateNode) error                                                                                                       //更新所有
	UpdateByAddresses(addresses []string, node *UpdateNode) error                                                                           //通过节点更新
	SelectAll() (list []*model.Node, err error)                                                                                             //查询所有
	SelectInAddress(addresses []string) (list []*model.Node, err error)                                                                     //根据节点地址查询
	Select(where *SelectNode, orderByField string, orderByType string, page int, pageSize int) (list []*model.Node, total int64, err error) //查询数据
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}

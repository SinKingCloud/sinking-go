package model

import (
	"gorm.io/gorm"
	"server/app/util/str"
	"time"
)

// Node 服务列表
type Node struct {
	Group        string       `gorm:"column:group" json:"group"`
	Name         string       `gorm:"column:name" json:"name"`
	Address      string       `gorm:"column:address" json:"address"`
	OnlineStatus int          `gorm:"column:online_status" json:"online_status"`
	Status       int          `gorm:"column:status" json:"status"`
	LastHeart    int64        `gorm:"column:last_heart" json:"last_heart"`
	IsDelete     int          `gorm:"column:is_delete" json:"is_delete"`
	CreateTime   str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime   str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Node) TableName() string {
	return "cloud_nodes"
}

// BeforeCreate 创建前
func (t *Node) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Node) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}

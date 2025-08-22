package model

import (
	"gorm.io/gorm"
	"server/app/util/str"
	"time"
)

// Cluster 集群列表
type Cluster struct {
	Address    string       `gorm:"column:address" json:"address"`
	Status     int          `gorm:"column:status" json:"status"`
	LastHeart  int64        `gorm:"column:last_heart" json:"last_heart"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Cluster) TableName() string {
	return "cloud_clusters"
}

// BeforeCreate 创建前
func (t *Cluster) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Cluster) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}

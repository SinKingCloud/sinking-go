package model

import (
	"gorm.io/gorm"
	"server/app/util/str"
	"time"
)

// Service 服务列表
type Service struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Group      string       `gorm:"column:group" json:"group"`
	Name       string       `gorm:"column:name" json:"name"`
	Address    string       `gorm:"column:address" json:"address"`
	Status     int          `gorm:"column:status" json:"status"`
	LastHeart  int64        `gorm:"column:last_heart" json:"last_heart"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Service) TableName() string {
	return "cloud_services"
}

// BeforeCreate 创建前
func (t *Service) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Service) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}

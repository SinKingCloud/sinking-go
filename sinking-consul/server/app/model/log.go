package model

import (
	"gorm.io/gorm"
	"server/app/util/str"
	"time"
)

// Log 日志表
type Log struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Type       int          `gorm:"column:type" json:"type"`
	Ip         string       `gorm:"column:ip" json:"ip"`
	Title      string       `gorm:"column:title" json:"title"`
	Content    string       `gorm:"column:content" json:"content"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Log) TableName() string {
	return "cloud_logs"
}

// BeforeCreate 创建前
func (t *Log) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Log) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}

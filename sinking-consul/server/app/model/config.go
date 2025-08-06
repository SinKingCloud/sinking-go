package model

import (
	"gorm.io/gorm"
	"server/app/util/str"
	"time"
)

// Config 配置表
type Config struct {
	Group      string       `gorm:"column:group" json:"group"`
	Name       string       `gorm:"column:name" json:"name"`
	Type       string       `gorm:"column:type" json:"type"`
	Hash       string       `gorm:"column:hash" json:"hash"`
	Content    string       `gorm:"column:content" json:"content"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Config) TableName() string {
	return "cloud_configs"
}

// BeforeCreate 创建前
func (t *Config) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Config) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}

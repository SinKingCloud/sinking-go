package model

import (
	"gorm.io/gorm"
	"time"
)

type Config struct {
	Id         int64    `gorm:"column:id" json:"id"`
	EnvId      int64    `gorm:"column:env_id" json:"env_id"`
	Title      string   `gorm:"column:title" json:"title"`
	Name       string   `gorm:"column:name" json:"name"`
	GroupName  string   `gorm:"column:group_name" json:"group_name"`
	Content    string   `gorm:"column:content" json:"content"`
	Hash       string   `gorm:"column:hash" json:"hash"`
	Type       string   `gorm:"column:type" json:"type"`
	Status     int      `gorm:"column:status" json:"status"`
	CreateTime DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int64    `gorm:"column:is_delete" json:"is_delete"`
}

func (Config) TableName() string {
	return DbPrefix + "configs"
}

// BeforeCreate 创建前
func (t *Config) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = DateTime(time.Now())
	t.IsDelete = 0
	return nil
}

// AfterFind 查询前
func (u *Config) AfterFind(tx *gorm.DB) (err error) {
	u.IsDelete = 0
	return
}

// BeforeUpdate 更新前
func (t *Config) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = DateTime(time.Now())
	return nil
}

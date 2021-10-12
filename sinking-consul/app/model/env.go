package model

import (
	"gorm.io/gorm"
	"time"
)

type Env struct {
	Id         int64    `gorm:"column:id"`
	AppId      int64    `gorm:"column:app_id"`
	Title      string   `gorm:"column:title" json:"title"`
	Name       string   `gorm:"column:name" json:"name"`
	CreateTime DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int64    `gorm:"column:is_delete" json:"is_delete"`
}

func (Env) TableName() string {
	return DbPrefix + "apps"
}

// BeforeCreate 创建前
func (t *Env) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = DateTime(time.Now())
	t.IsDelete = 0
	return nil
}

// AfterFind 查询前
func (u *Env) AfterFind(tx *gorm.DB) (err error) {
	u.IsDelete = 0
	return
}

// BeforeUpdate 更新前
func (t *Env) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = DateTime(time.Now())
	return nil
}

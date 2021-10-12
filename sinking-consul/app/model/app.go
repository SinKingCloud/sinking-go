package model

import (
	"gorm.io/gorm"
	"time"
)

type App struct {
	Id         int64    `gorm:"column:id"`
	Title      string   `gorm:"column:title" json:"title"`
	Name       string   `gorm:"column:name" json:"name"`
	CreateTime DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int64    `gorm:"column:is_delete" json:"is_delete"`
}

func (App) TableName() string {
	return DbPrefix + "apps"
}

// BeforeCreate 创建前
func (t *App) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = DateTime(time.Now())
	t.IsDelete = 0
	return nil
}

// AfterFind 查询前
func (u *App) AfterFind(tx *gorm.DB) (err error) {
	u.IsDelete = 0
	return
}

// BeforeUpdate 更新前
func (t *App) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = DateTime(time.Now())
	return nil
}

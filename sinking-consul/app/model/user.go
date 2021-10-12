package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id         int64    `gorm:"column:id"`
	RoleId     int64    `gorm:"column:role_id" json:"role_id"`
	User       string   `gorm:"column:user" json:"user"`
	Pwd        string   `gorm:"column:pwd" json:"pwd"`
	Name       string   `gorm:"column:name" json:"name"`
	LoginIp    string   `gorm:"column:login_ip" json:"login_ip"`
	LoginTime  DateTime `gorm:"column:login_time" json:"login_time"`
	CreateTime DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int64    `gorm:"column:is_delete" json:"is_delete"`
}

func (User) TableName() string {
	return DbPrefix + "users"
}

// BeforeCreate 创建前
func (t *User) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = DateTime(time.Now())
	t.IsDelete = 0
	return nil
}

// AfterFind 查询前
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.IsDelete = 0
	return
}

// BeforeUpdate 更新前
func (t *User) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = DateTime(time.Now())
	return nil
}

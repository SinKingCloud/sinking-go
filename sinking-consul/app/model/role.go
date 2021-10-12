package model

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cachePrefix"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/cache"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Role struct {
	Id         int64    `gorm:"column:id"`
	Name       string   `gorm:"column:name" json:"name"`
	Auths      string   `gorm:"column:auths" json:"auths"`
	CreateTime DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int64    `gorm:"column:is_delete" json:"is_delete"`
}

func (r *Role) FindByIdCache() *Role {
	data := cache.Remember(cachePrefix.Role+strconv.FormatInt(r.Id, 10), func() interface{} {
		var info *Role
		Db.Where("id=?", r.Id).First(&info)
		return info
	}, 600*time.Second)
	return data.(*Role)
}

func (Role) TableName() string {
	return DbPrefix + "roles"
}

// BeforeCreate 创建前
func (t *Role) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = DateTime(time.Now())
	t.IsDelete = 0
	return nil
}

// AfterFind 查询前
func (u *Role) AfterFind(tx *gorm.DB) (err error) {
	u.IsDelete = 0
	return
}

// BeforeUpdate 更新前
func (t *Role) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = DateTime(time.Now())
	return nil
}
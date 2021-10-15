package model

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cachePrefix"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cacheTime"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/cache"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type User struct {
	Id         int64    `gorm:"column:id" json:"id"`
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

func (r *User) FindByIdCache() *User {
	data := cache.Remember(cachePrefix.User+strconv.FormatInt(r.Id, 10), func() interface{} {
		var info *User
		Db.Where("id=? and is_delete=0", r.Id).First(&info)
		return info
	}, cacheTime.Time*time.Second)
	return data.(*User)
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

// BeforeUpdate 更新前
func (t *User) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = DateTime(time.Now())
	return nil
}

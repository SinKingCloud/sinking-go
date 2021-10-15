package model

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cachePrefix"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cacheTime"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/cache"
	"gorm.io/gorm"
	"time"
)

type Config struct {
	Id         int64    `gorm:"column:id" json:"id"`
	AppName    string   `gorm:"column:app_name" json:"app_name"`
	EnvName    string   `gorm:"column:env_name" json:"env_name"`
	GroupName  string   `gorm:"column:group_name" json:"group_name"`
	Name       string   `gorm:"column:name" json:"name"`
	Title      string   `gorm:"column:title" json:"title"`
	Content    string   `gorm:"column:content" json:"content"`
	Hash       string   `gorm:"column:hash" json:"hash"`
	Type       string   `gorm:"column:type" json:"type"`
	Status     int      `gorm:"column:status" json:"status"`
	CreateTime DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int64    `gorm:"column:is_delete" json:"is_delete"`
}

func (r *Config) SelectByNameCache() []*service.Config {
	data := cache.Remember(cachePrefix.Config+r.AppName+r.EnvName, func() interface{} {
		var info []*service.Config
		Db.Model(&Config{}).Where("app_name=? and env_name=? and status=0 and is_delete=0", r.AppName, r.EnvName).Find(&info)
		return info
	}, cacheTime.Time*time.Second)
	return data.([]*service.Config)
}

func (r *Config) SelectAllCache() []*service.Config {
	data := cache.Remember(cachePrefix.Config, func() interface{} {
		var info []*service.Config
		Db.Model(&Config{}).Find(&info)
		return info
	}, cacheTime.Time*time.Second)
	return data.([]*service.Config)
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

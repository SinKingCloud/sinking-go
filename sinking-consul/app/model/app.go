package model

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cachePrefix"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cacheTime"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/cache"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type App struct {
	Id         int64    `gorm:"column:id" json:"id"`
	Title      string   `gorm:"column:title" json:"title"`
	Name       string   `gorm:"column:name" json:"name"`
	CreateTime DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime DateTime `gorm:"column:update_time" json:"update_time"`
	IsDelete   int64    `gorm:"column:is_delete" json:"is_delete"`
}

func (r *App) FindByIdCache() *App {
	data := cache.Remember(cachePrefix.App+strconv.FormatInt(r.Id, 10), func() interface{} {
		var info *App
		Db.Where("id=? and is_delete=0", r.Id).First(&info)
		return info
	}, 600*time.Second)
	return data.(*App)
}

func (r *App) FindByNameCache() *App {
	data := cache.Remember(cachePrefix.App+r.Name, func() interface{} {
		var info *App
		Db.Where("name=? and is_delete=0", r.Name).First(&info)
		return info
	}, cacheTime.Time*time.Second)
	return data.(*App)
}

func (r *App) SelectAllCache() []*App {
	data := cache.Remember(cachePrefix.App, func() interface{} {
		var info []*App
		Db.Model(&App{}).Find(&info)
		return info
	}, cacheTime.Time*time.Second)
	return data.([]*App)
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

// BeforeUpdate 更新前
func (t *App) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = DateTime(time.Now())
	return nil
}

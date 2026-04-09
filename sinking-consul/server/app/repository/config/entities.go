package config

import "server/app/util/str"

type UpdateConfig struct {
	Group   interface{}
	Name    interface{}
	Type    interface{}
	Hash    interface{}
	Content interface{}
	Status  interface{}
}

type SelectConfig struct {
	Group           string
	Name            string
	Type            string
	Hash            string
	Content         string
	Status          string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

type Config struct {
	Group      string       `gorm:"column:group" json:"group"`
	Name       string       `gorm:"column:name" json:"name"`
	Type       string       `gorm:"column:type" json:"type"`
	Hash       string       `gorm:"column:hash" json:"hash"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

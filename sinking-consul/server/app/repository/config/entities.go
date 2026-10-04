package config

import "server/app/util/str"

type UpdateConfig struct {
	Group   *string
	Name    *string
	Type    *string
	Hash    *string
	Content *string
	Status  *int
}

type SelectConfig struct {
	Keyword         *string
	Group           *string
	Name            *string
	Type            *string
	Hash            *string
	Content         *string
	Status          *int
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
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

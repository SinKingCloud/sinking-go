package service

type Config struct {
	AppName   string `gorm:"column:app_name" json:"app_name"`
	EnvName   string `gorm:"column:env_name" json:"env_name"`
	GroupName string `gorm:"column:group_name" json:"group_name"`
	Name      string `gorm:"column:name" json:"name"`
	Content   string `gorm:"column:content" json:"content"`
	Hash      string `gorm:"column:hash" json:"hash"`
	Type      string `gorm:"column:type" json:"type"`
	Status    int    `gorm:"column:status" json:"status"`
}

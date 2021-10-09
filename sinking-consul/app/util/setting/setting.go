package setting

import "github.com/spf13/viper"

var (
	system *viper.Viper
)

// SetSetting 设置数据
func SetSetting(setting *viper.Viper) {
	system = setting
}

// GetConfig 获取设置数据
func GetConfig() *viper.Viper {
	return system
}

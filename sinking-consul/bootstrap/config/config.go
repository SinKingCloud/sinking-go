package config

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/spf13/viper"
)

func LoadConfig(configPath string, configName string, configType string) {
	config := viper.New()
	config.AddConfigPath(configPath)
	config.SetConfigName(configName)
	config.SetConfigType(configType)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
		return
	}
	setting.SetSetting(config)
}

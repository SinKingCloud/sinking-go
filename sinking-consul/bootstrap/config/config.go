package config

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/spf13/viper"
	"strings"
	"time"
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
	var conf setting.SystemConfig
	if err := config.Unmarshal(&conf); err != nil {
		panic(err)
		return
	}
	setting.SetSystemSetting(&conf)
	//加载注册节点
	loadRegisterServers()
}

func loadRegisterServers() {
	//设置注册节点
	servers := strings.Split(setting.GetSystemConfig().Servers.Cluster, ",")
	for _, v := range servers {
		server := strings.Split(v, ":")
		if len(server) == 2 {
			info := &service.Cluster{
				Hash:          encode.Md5Encode(server[0] + ":" + server[1]),
				Ip:            server[0],
				Port:          server[1],
				LastHeartTime: model.DateTime(time.Now()),
				Status:        0,
			}
			service.RegisterClusters[info.Hash] = info
		}
	}
}

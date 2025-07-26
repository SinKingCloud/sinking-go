package bootstrap

import (
	"github.com/spf13/viper"
	"server/app/constant"
	"server/app/util"
	"server/app/util/file"
	"strings"
)

// LoadConf 加载本地配置
func LoadConf() {
	path := constant.ConfPath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	configType := "yml"
	fileName := constant.ConfFile
	disk := file.NewDisk(path)
	_ = disk.AutoCreate(fileName + "." + configType)
	config := viper.New()
	config.AutomaticEnv() //读取环境变量
	config.AddConfigPath(path)
	config.SetConfigName(fileName)
	config.SetConfigType(configType)
	config.WatchConfig()
	if err := config.ReadInConfig(); err != nil {
		panic(err)
		return
	}
	//赋值到util
	util.Conf = config
}

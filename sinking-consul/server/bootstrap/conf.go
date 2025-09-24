package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	"server/app/constant"
	"server/app/util"
	"server/app/util/file"
	"strings"
)

func LoadConf() {
	path := constant.ConfPath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	configType := "yml"
	fileName := constant.ConfFile
	disk := file.NewDisk(path)
	if err := disk.AutoCreate(fileName + "." + configType); err != nil {
		panic(fmt.Errorf("创建配置文件失败: %w", err))
		return
	}
	config := viper.New()
	config.AddConfigPath(path)
	config.SetConfigName(fileName)
	config.SetConfigType(configType)
	config.WatchConfig()
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	keys := []string{
		constant.ServerPort,
		constant.ServerHost,
		constant.ServerMode,
		constant.AuthAccount,
		constant.AuthPassword,
		constant.AuthApiToken,
		constant.AuthExpire,
		constant.ClusterLocal,
		constant.ClusterNodes,
	}
	for _, key := range keys {
		if !config.IsSet(key) {
			_ = config.BindEnv(key)
		}
	}
	util.Conf = config
}

package setting

import (
	"github.com/spf13/viper"
)

// Service 接口
type Service interface {
	GetByGroup(key string) interface{}  //通过group获取配置
	GetString(key string) string        //获取配置
	GetInt(key string) int              //获取配置
	GetStringSlice(key string) []string //获取配置
	Set(key string, value string) error //设置配置
	Sets(list []*Config) error          //批量设置配置
}

// service 服务
type service struct {
	conf *viper.Viper
}

// NewService 实例化service
func NewService(conf *viper.Viper) *service {
	return &service{
		conf: conf,
	}
}

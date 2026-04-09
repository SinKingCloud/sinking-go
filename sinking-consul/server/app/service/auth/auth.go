package auth

import (
	"server/app/service/setting"
	"server/app/util/cache"
)

// Service 接口
type Service interface {
	GetCaptcha(key string) (map[string]interface{}, error) //获取验证码
	CheckCaptcha(key string, x int, y int) bool            //判断验证码
	CheckAccount(account string, pwd string) error         //判断账户密码
	GenLoginToken(ip string) (s string, e error)           //生成登陆token
	ClearLoginToken() error                                //清理登陆token
	ChangePassword(password string) error                  //修改密码
	CheckLoginToken(token string) error                    //判断登陆token
	CheckApiToken(token string) error                      //判断API token
}

// service 服务
type service struct {
	cache   cache.Interface
	setting setting.Service
}

// NewService 实例化service
func NewService(setting setting.Service, cache cache.Interface) *service {
	return &service{
		setting: setting,
		cache:   cache,
	}
}

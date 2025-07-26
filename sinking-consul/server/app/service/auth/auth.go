package auth

import (
	"errors"
	"server/app/constant"
	"server/app/service/config"
	"server/app/util/jwt"
	"server/app/util/str"
	"strconv"
	"sync"
	"time"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// CheckAccount 判断账号密码
func (*Service) CheckAccount(account string, pwd string) error {
	sUser := config.GetIns().Get(constant.LoginGroup, constant.LoginAccount)
	sPwd := config.GetIns().Get(constant.LoginGroup, constant.LoginPassword)
	if sUser == "" || sPwd == "" || sUser != account || !str.NewStringTool(sPwd).CheckPassword(pwd) {
		return errors.New("用户名或密码错误")
	}
	return nil
}

// GenLoginToken 生成jwtToken
func (*Service) GenLoginToken(types string, ip string) (s string, e error) {
	token := str.NewStringTool(strconv.FormatInt(time.Now().UnixMilli(), 10)).Md5()
	loginTime := str.DateTime(time.Now())
	if token != "" {
		e = config.GetIns().Set(constant.LoginGroup, constant.LoginToken+"."+types, token)
	} else {
		e = errors.New("生成token失败")
		return
	}
	if e != nil {
		return "", e
	}
	expire, _ := strconv.Atoi(config.GetIns().Get(constant.LoginGroup, constant.LoginExpire))
	if expire > 0 && expire <= 600 {
		expire = 600
	}
	s = jwt.GetToken(&jwt.User{
		LoginToken: token,
		LoginIp:    ip,
		LoginTime:  loginTime,
	}, expire)
	return s, nil
}

// ClearLoginToken 清理jwtToken
func (*Service) ClearLoginToken(types string) error {
	err := config.GetIns().Set(constant.LoginGroup, constant.LoginToken+"."+types, "")
	if err != nil {
		return err
	}
	return nil
}

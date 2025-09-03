package auth

import (
	"errors"
	"server/app/constant"
	"server/app/util"
	"server/app/util/jwt"
	"server/app/util/str"
	"strconv"
	"time"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj = &Service{}
)

// GetIns 获取单例
func GetIns() *Service {
	return obj
}

// CheckAccount 判断账号密码
func (*Service) CheckAccount(account string, pwd string) error {
	if account == "" || pwd == "" {
		return errors.New("用户名或密码不能为空")
	}
	sUser := util.Conf.GetString(constant.AuthAccount)
	sPwd := util.Conf.GetString(constant.AuthPassword)
	if sUser == "" || sPwd == "" {
		sUser = account
		util.Conf.Set(constant.AuthAccount, account)
		sPwd, _ = str.NewStringTool().BcryptHash(pwd)
		util.Conf.Set(constant.AuthPassword, sPwd)
		if util.Conf.WriteConfig() != nil {
			return errors.New("初始化账号密码失败")
		}
	}
	if ok, _ := str.NewStringTool().BcryptVerify(pwd, sPwd); !ok || sUser != account {
		return errors.New("用户名或密码错误")
	}
	return nil
}

// GenLoginToken 生成jwtToken
func (*Service) GenLoginToken(ip string) (s string, e error) {
	token := str.NewStringTool().Md5(strconv.FormatInt(time.Now().UnixMilli(), 10))
	if token != "" {
		util.Conf.Set(constant.AuthLoginToken, token)
		e = util.Conf.WriteConfig()
	} else {
		e = errors.New("生成token失败")
		return
	}
	if e != nil {
		return "", e
	}
	loginTime := str.DateTime(time.Now())
	expire := util.Conf.GetInt(constant.AuthExpire)
	if expire <= 86400 {
		expire = 86400
	}
	s = jwt.GetToken(&jwt.User{
		LoginToken: token,
		LoginIp:    ip,
		LoginTime:  loginTime,
	}, expire)
	return s, nil
}

// ClearLoginToken 清理jwtToken
func (*Service) ClearLoginToken() error {
	util.Conf.Set(constant.AuthLoginToken, "")
	if util.Conf.WriteConfig() != nil {
		return errors.New("清理登录token失败")
	}
	return nil
}

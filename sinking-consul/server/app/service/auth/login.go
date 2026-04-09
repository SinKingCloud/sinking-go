package auth

import (
	"errors"
	"server/app/constant"
	"server/app/service/setting"
	"server/app/util/jwt"
	"server/app/util/str"
	"strconv"
	"time"
)

// CheckAccount 判断账号密码
func (c *service) CheckAccount(account string, pwd string) error {
	if account == "" || pwd == "" {
		return errors.New("用户名或密码不能为空")
	}
	sUser := c.setting.GetString(constant.AuthAccount)
	sPwd := c.setting.GetString(constant.AuthPassword)
	if sUser == "" || sPwd == "" {
		sUser = account
		sPwd, _ = str.NewStringTool().BcryptHash(pwd)
		err := c.setting.Sets([]*setting.Config{
			{Key: constant.AuthAccount, Value: account},
			{Key: constant.AuthPassword, Value: sPwd},
		})
		if err != nil {
			return errors.New("初始化账号密码失败")
		}
	}
	if ok, _ := str.NewStringTool().BcryptVerify(pwd, sPwd); !ok || sUser != account {
		return errors.New("用户名或密码错误")
	}
	return nil
}

// GenLoginToken 生成jwtToken
func (c *service) GenLoginToken(ip string) (s string, e error) {
	token := str.NewStringTool().Md5(strconv.FormatInt(time.Now().UnixMilli(), 10))
	if token != "" {
		e = c.setting.Set(constant.AuthLoginToken, token)
	} else {
		e = errors.New("生成token失败")
		return
	}
	if e != nil {
		return "", e
	}
	loginTime := str.DateTime(time.Now())
	expire := c.setting.GetInt(constant.AuthExpire)
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
func (c *service) ClearLoginToken() error {
	if c.setting.Set(constant.AuthLoginToken, "") != nil {
		return errors.New("清理登录token失败")
	}
	return nil
}

// ChangePassword 修改密码
func (c *service) ChangePassword(password string) error {
	sPwd, _ := str.NewStringTool().BcryptHash(password)
	err := c.setting.Set(constant.AuthPassword, sPwd)
	if err != nil {
		return errors.New("修改失败")
	}
	return nil
}

// CheckLoginToken 判断登陆token
func (c *service) CheckLoginToken(token string) error {
	loginToken := c.setting.GetString(constant.AuthLoginToken)
	if loginToken == "" {
		return errors.New("您的账户已注销登陆,请重新登陆")
	}
	if token != loginToken {
		return errors.New("您的账户已在其他设备登陆,请重新登陆")
	}
	return nil
}

// CheckApiToken 判断API token
func (c *service) CheckApiToken(token string) error {
	if token == "" {
		return errors.New("鉴权token缺失,请检查请求")
	}
	if token != c.setting.GetString(constant.AuthApiToken) {
		return errors.New("鉴权token无效,请确认token是否正确")
	}
	return nil
}

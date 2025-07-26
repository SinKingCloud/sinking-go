package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"server/app/constant"
	"server/app/util/str"
	"time"
)

// MyClaims jwt载体
type MyClaims struct {
	User *User
	jwt.RegisteredClaims
}

// User 会员用户表
type User struct {
	LoginToken string       ` json:"login_token"`
	LoginIp    string       ` json:"login_ip"`
	LoginTime  str.DateTime ` json:"login_time"`
}

// getKey 获取加密key
func getKey() []byte {
	return []byte(constant.JwtKey)
}

// GetToken 生成token user 用户信息
func GetToken(user *User, expireTime int) string {
	if expireTime == 0 {
		expireTime = constant.JwtExpireTime
	}
	setClaim := MyClaims{
		User: user,
	}
	if expireTime >= 0 {
		setClaim.RegisteredClaims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime) * time.Second)), //过期时间（时间戳）
		}
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaim) //生成token
	token, err := reqClaim.SignedString(getKey())                   //转换为字符串
	if err != nil {
		return ""
	}
	return token
}

// CheckToken 验证token
func CheckToken(token string) *MyClaims {
	setToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getKey(), nil
	})
	if err != nil {
		return nil
	}
	if key, success := setToken.Claims.(*MyClaims); setToken.Valid && success {
		return key
	} else {
		return nil
	}
}

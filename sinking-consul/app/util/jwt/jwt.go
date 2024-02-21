package jwt

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/golang-jwt/jwt"
	"time"
)

type MyClaims struct {
	User model.User
	jwt.StandardClaims
}

func getKey() []byte {
	return []byte(setting.GetSystemConfig().App.JwtToken)
}

// SetToken 生成token
func SetToken(user model.User) string {
	expireTime := time.Now().Add(86400 * time.Second) //过期时间
	setClaim := MyClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间（时间戳）
		},
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

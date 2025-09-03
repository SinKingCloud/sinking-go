package constant

import "time"

const (
	CacheNameWithCaptcha = "Captcha_"       //验证码
	CacheTimeWithCaptcha = 10 * time.Minute //验证码缓存时间 10分钟

	CacheNameWithLocalIp = "LocalIp"        //本机IP
	CacheTimeWithLocalIp = 10 * time.Minute //本机IP缓存时间 10分钟
)

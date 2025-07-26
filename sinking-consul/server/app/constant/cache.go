package constant

import "time"

const (
	CacheNameWithSysConfig = "Sys"             //系统配置
	CacheTimeWithSysConfig = 600 * time.Second //系统配置配置储存时间

	CacheNameWithCaptcha = "Captcha_" //验证码
	CacheTimeWithCaptcha = 600 * 1000 //验证码缓存时间
)

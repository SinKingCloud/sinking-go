package auth

import (
	"encoding/json"
	"errors"
	"github.com/wenlng/go-captcha/v2/slide"
	"server/app/constant"
	"server/app/util"
	"server/app/util/captcha"
)

// GetCaptcha 获取验证码
func (c *Service) GetCaptcha(key string) (map[string]interface{}, error) {
	result, masterImage, tileImage, height, width, err := captcha.GenSlide()
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return nil, errors.New("获取验证码失败")
	}
	util.Cache.SetWithExpire(constant.CacheNameWithCaptcha+key, string(bytes), constant.CacheTimeWithCaptcha)
	ret := map[string]interface{}{
		"key":          key,
		"image_base64": masterImage,
		"width":        width,
		"height":       height,
		"tile_base64":  tileImage,
		"tile_width":   result.Width,
		"tile_height":  result.Height,
		"tile_x":       result.DX,
		"tile_y":       result.DY,
	}
	return ret, nil
}

// CheckCaptcha 判断验证码是否正确
func (c *Service) CheckCaptcha(key string, x int, y int) bool {
	value := util.Cache.Get(constant.CacheNameWithCaptcha + key)
	if value == nil {
		return false
	}
	util.Cache.Delete(constant.CacheNameWithCaptcha + key)
	var dct *slide.Block
	if err := json.Unmarshal([]byte(value.(string)), &dct); err != nil {
		return false
	}
	return captcha.CheckSlide(x, y, dct)
}

package auth

import (
	"encoding/json"
	"errors"
	"server/app/constant"
	"server/app/util/captcha"

	"github.com/wenlng/go-captcha/v2/slide"
)

// GetCaptcha 获取验证码
func (c *service) GetCaptcha(key string) (map[string]interface{}, error) {
	result, masterImage, tileImage, height, width, err := captcha.GenSlide()
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return nil, errors.New("获取验证码失败")
	}
	c.cache.SetWithExpire(constant.CacheNameWithCaptcha+key, string(bytes), constant.CacheTimeWithCaptcha)
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
func (c *service) CheckCaptcha(key string, x int, y int) bool {
	value := c.cache.Get(constant.CacheNameWithCaptcha + key)
	if value == nil {
		return false
	}
	c.cache.Delete(constant.CacheNameWithCaptcha + key)
	var dct *slide.Block
	if err := json.Unmarshal([]byte(value.(string)), &dct); err != nil {
		return false
	}
	return captcha.CheckSlide(x, y, dct)
}

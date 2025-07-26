package captcha

import (
	"github.com/afocus/captcha"
	"image/color"
	"sync"
)

var (
	c    *captcha.Captcha
	once sync.Once
)

// getIns 获取单例
func getIns() *captcha.Captcha {
	once.Do(func() {
		c = captcha.New()
		if err := c.AddFontFromBytes(NewTtf()); err != nil {
			panic(err.Error())
		}
		c.SetSize(256, 128)
		c.SetDisturbance(64)
		c.SetFrontColor(color.RGBA{R: 0, G: 81, B: 235, A: 255})
	})
	return c
}

// Image 对象
type Image struct {
	*captcha.Image
}

// NewCaptchaCode 生成验证码图片
// code 验证码
func NewCaptchaCode(code string) *Image {
	img := getIns().CreateCustom(code)
	return &Image{img}
}

// NewCaptchaRandCode 生成随机验证码图片
// num 验证码个数
func NewCaptchaRandCode(num int) (*Image, string) {
	img, code := getIns().Create(num, captcha.ALL)
	return &Image{img}, code
}

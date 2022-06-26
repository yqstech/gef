/**
 * @Author: 云起时
 * @Email: limingxiang@yqstech.com
 * @Description:
 * @File: base64Captcha
 * @Version: 1.0.0
 * @Date: 2022/3/11 9:09 上午
 */

package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
	"math/rand"
)

var store = base64Captcha.DefaultMemStore
var ColorWight = color.RGBA{
	R: 255,
	G: 255,
	B: 255,
	A: 255,
}

func GetCaptchaBase64(CaptchaType string, height, width, length int, bgColor color.RGBA) (string, string, error) {
	CaptchaTypes := map[int]string{
		0: "string", 1: "math", 2: "digit",
	}
	if CaptchaType == "auto" {
		t := rand.Intn(2)
		CaptchaType = CaptchaTypes[t]
	}
	//定义字体库和字体列表
	fontStorage := base64Captcha.DefaultEmbeddedFonts
	fonts := []string{
		"3Dumb.ttf",
		"ApothecaryFont.ttf",
		"Comismsh.ttf",
		"DENNEthree-dee.ttf",
		"DeborahFancyDress.ttf",
		"Flim-Flam.ttf",
		"RitaSmith.ttf",
	}
	//获取驱动
	var driver base64Captcha.Driver
	switch CaptchaType {
	case "string":
		driver = base64Captcha.NewDriverString(height, width, 5, 3, length, "1234567890qwertyuioplkjhgfdsazxcvbnm", &bgColor, fontStorage, fonts)
	case "math":
		driver = base64Captcha.NewDriverMath(height, width, 5, 2, &bgColor, fontStorage, fonts)
	default:
		driver = base64Captcha.NewDriverDigit(height, width, length, 0.7, 80)
	}
	//创建验证码
	c := base64Captcha.NewCaptcha(driver, store)
	return c.Generate()
}

func Verify(captchaId, code string) bool {
	return store.Verify(captchaId, code, true)
}

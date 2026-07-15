package api

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model/response"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

func Captcha(c *gin.Context) {
	// 数字验证码：降低扭曲与噪点，配合前端浅色底更易识别
	driver := base64Captcha.NewDriverDigit(
		global.GVA_CONFIG.Captcha.ImgHeight,
		global.GVA_CONFIG.Captcha.ImgWidth,
		global.GVA_CONFIG.Captcha.KeyLong,
		0.45,
		40,
	)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, _, err := cp.Generate(); err != nil {
		global.GVA_LOG.Error("验证码获取失败!", zap.Any("err", err))
		response.FailWithMessage("验证码获取失败", c)
	} else {
		response.OkWithDetailed(response.SysCaptchaResponse{
			CaptchaId: id,
			PicPath:   b64s,
		}, "验证码获取成功", c)
	}
}

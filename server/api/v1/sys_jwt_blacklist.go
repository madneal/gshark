package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

func JsonInBlacklist(c *gin.Context) {
	token := c.Request.Header.Get("x-token")
	jwt := model.JwtBlacklist{Jwt: token}
	if err := service.JsonInBlacklist(jwt); err != nil {
		global.GVA_LOG.Error("jwt作废失败!", zap.Any("err", err))
		response.FailWithMessage("jwt作废失败", c)
	} else {
		response.OkWithMessage("jwt作废成功", c)
	}
}

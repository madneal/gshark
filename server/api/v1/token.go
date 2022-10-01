package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

func CreateToken(c *gin.Context) {
	var token model.Token
	_ = c.ShouldBindJSON(&token)
	if err := service.CreateToken(token); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func DeleteToken(c *gin.Context) {
	var token model.Token
	_ = c.ShouldBindJSON(&token)
	if err := service.DeleteToken(token); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func DeleteTokenByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteTokenByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

func UpdateToken(c *gin.Context) {
	var token model.Token
	_ = c.ShouldBindJSON(&token)
	if err := service.UpdateToken(token); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

func FindToken(c *gin.Context) {
	var token model.Token
	_ = c.ShouldBindQuery(&token)
	if err, retoken := service.GetToken(token.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"retoken": retoken}, c)
	}
}

func GetTokenList(c *gin.Context) {
	var pageInfo request.TokenSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetTokenInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

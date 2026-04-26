package api

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model/response"
	"go.uber.org/zap"
)

func bindJSON(c *gin.Context, value interface{}) bool {
	if err := c.ShouldBindJSON(value); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return false
	}
	return true
}

func bindQuery(c *gin.Context, value interface{}) bool {
	if err := c.ShouldBindQuery(value); err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return false
	}
	return true
}

func respondMutation(c *gin.Context, err error, logMessage, failMessage, successMessage string) {
	if err != nil {
		global.GVA_LOG.Error(logMessage, zap.Any("err", err))
		response.FailWithMessage(failMessage, c)
		return
	}
	response.OkWithMessage(successMessage, c)
}

func respondPage(c *gin.Context, err error, list interface{}, total int64, page, pageSize int) {
	if err != nil {
		global.GVA_LOG.Error("获取失败", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

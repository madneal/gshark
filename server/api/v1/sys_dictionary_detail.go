package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
	"github.com/madneal/gshark/utils"
	"go.uber.org/zap"
)

func CreateSysDictionaryDetail(c *gin.Context) {
	var detail model.SysDictionaryDetail
	_ = c.ShouldBindJSON(&detail)
	if err := service.CreateSysDictionaryDetail(detail); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func DeleteSysDictionaryDetail(c *gin.Context) {
	var detail model.SysDictionaryDetail
	_ = c.ShouldBindJSON(&detail)
	if err := service.DeleteSysDictionaryDetail(detail); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func UpdateSysDictionaryDetail(c *gin.Context) {
	var detail model.SysDictionaryDetail
	_ = c.ShouldBindJSON(&detail)
	if err := service.UpdateSysDictionaryDetail(&detail); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

func FindSysDictionaryDetail(c *gin.Context) {
	var detail model.SysDictionaryDetail
	_ = c.ShouldBindQuery(&detail)
	if err := utils.Verify(detail, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, resysDictionaryDetail := service.GetSysDictionaryDetail(detail.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithDetailed(gin.H{"resysDictionaryDetail": resysDictionaryDetail}, "查询成功", c)
	}
}

func GetSysDictionaryDetailList(c *gin.Context) {
	var pageInfo request.SysDictionaryDetailSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetSysDictionaryDetailInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
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

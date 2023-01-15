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

func CreateFilter(c *gin.Context) {
	var filter model.Filter
	_ = c.ShouldBindJSON(&filter)
	if err := service.CreateFilter(filter); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func DeleteFilter(c *gin.Context) {
	var filter model.Filter
	_ = c.ShouldBindJSON(&filter)
	if err := service.DeleteFilter(filter); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func DeleteFilterByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteFilterByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

func UpdateFilter(c *gin.Context) {
	var filter model.Filter
	_ = c.ShouldBindJSON(&filter)
	if err := service.UpdateFilter(filter); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

func FindFilter(c *gin.Context) {
	var filter model.Filter
	_ = c.ShouldBindQuery(&filter)
	if err, filter := service.GetFilter(filter.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"filter": filter}, c)
	}
}

func GetFilterList(c *gin.Context) {
	var pageInfo request.FilterSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetFilterInfoList(pageInfo); err != nil {
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

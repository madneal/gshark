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

func CreateApi(c *gin.Context) {
	var api model.SysApi
	_ = c.ShouldBindJSON(&api)
	if err := utils.Verify(api, utils.ApiVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.CreateApi(api); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func DeleteApi(c *gin.Context) {
	var api model.SysApi
	_ = c.ShouldBindJSON(&api)
	if err := utils.Verify(api.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.DeleteApi(api); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func GetApiList(c *gin.Context) {
	var pageInfo request.SearchApiParams
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, list, total := service.GetAPIInfoList(pageInfo.SysApi, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc); err != nil {
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

func GetApiById(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, api := service.GetApiById(idInfo.Id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(response.SysAPIResponse{Api: api}, c)
	}
}

func UpdateApi(c *gin.Context) {
	var api model.SysApi
	_ = c.ShouldBindJSON(&api)
	if err := utils.Verify(api, utils.ApiVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.UpdateApi(api); err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Any("err", err))
		response.FailWithMessage("修改失败", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

func GetAllApis(c *gin.Context) {
	if err, apis := service.GetAllApis(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.SysAPIListResponse{Apis: apis}, "获取成功", c)
	}
}

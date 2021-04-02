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

// @Tags Subdomain
// @Summary 创建Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "创建Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /subdomain/createSubdomain [post]
func CreateSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	_ = c.ShouldBindJSON(&subdomain)
	if err := service.CreateSubdomain(subdomain); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags Subdomain
// @Summary 删除Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "删除Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /subdomain/deleteSubdomain [delete]
func DeleteSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	_ = c.ShouldBindJSON(&subdomain)
	if err := service.DeleteSubdomain(subdomain); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Subdomain
// @Summary 批量删除Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /subdomain/deleteSubdomainByIds [delete]
func DeleteSubdomainByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteSubdomainByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags Subdomain
// @Summary 更新Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "更新Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /subdomain/updateSubdomain [put]
func UpdateSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	_ = c.ShouldBindJSON(&subdomain)
	if err := service.UpdateSubdomain(subdomain); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Subdomain
// @Summary 用id查询Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "用id查询Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /subdomain/findSubdomain [get]
func FindSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	_ = c.ShouldBindQuery(&subdomain)
	if err, resubdomain := service.GetSubdomain(subdomain.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"resubdomain": resubdomain}, c)
	}
}

// @Tags Subdomain
// @Summary 分页获取Subdomain列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SubdomainSearch true "分页获取Subdomain列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /subdomain/getSubdomainList [get]
func GetSubdomainList(c *gin.Context) {
	var pageInfo request.SubdomainSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetSubdomainInfoList(pageInfo); err != nil {
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

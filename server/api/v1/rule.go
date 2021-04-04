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

// @Tags Rule
// @Summary 创建Rule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Rule true "创建Rule"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /rule/createRule [post]
func CreateRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindJSON(&rule)
	if err := service.CreateRule(rule); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags Rule
// @Summary 删除Rule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Rule true "删除Rule"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /rule/deleteRule [delete]
func DeleteRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindJSON(&rule)
	if err := service.DeleteRule(rule); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Rule
// @Summary 批量删除Rule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Rule"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /rule/deleteRuleByIds [delete]
func DeleteRuleByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteRuleByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags Rule
// @Summary 更新Rule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Rule true "更新Rule"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /rule/updateRule [put]
func UpdateRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindJSON(&rule)
	if err := service.UpdateRule(rule); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Rule
// @Summary 用id查询Rule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Rule true "用id查询Rule"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /rule/findRule [get]
func FindRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindQuery(&rule)
	if err, rerule := service.GetRule(rule.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rerule": rerule}, c)
	}
}

// @Tags Rule
// @Summary 分页获取Rule列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.RuleSearch true "分页获取Rule列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /rule/getRuleList [get]
func GetRuleList(c *gin.Context) {
	var pageInfo request.RuleSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetRuleInfoList(pageInfo); err != nil {
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

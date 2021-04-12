package v1

import (
	"github.com/madneal/gshark/global"
    "github.com/madneal/gshark/model"
    "github.com/madneal/gshark/model/request"
    "github.com/madneal/gshark/model/response"
    "github.com/madneal/gshark/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// @Tags Filter
// @Summary 创建Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "创建Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /filter/createFilter [post]
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

// @Tags Filter
// @Summary 删除Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "删除Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /filter/deleteFilter [delete]
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

// @Tags Filter
// @Summary 批量删除Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /filter/deleteFilterByIds [delete]
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

// @Tags Filter
// @Summary 更新Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "更新Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /filter/updateFilter [put]
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

// @Tags Filter
// @Summary 用id查询Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "用id查询Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /filter/findFilter [get]
func FindFilter(c *gin.Context) {
	var filter model.Filter
	_ = c.ShouldBindQuery(&filter)
	if err, refilter := service.GetFilter(filter.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"refilter": refilter}, c)
	}
}

// @Tags Filter
// @Summary 分页获取Filter列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.FilterSearch true "分页获取Filter列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /filter/getFilterList [get]
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

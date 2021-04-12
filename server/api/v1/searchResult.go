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

// @Tags SearchResult
// @Summary 创建SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "创建SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /searchResult/createSearchResult [post]
func CreateSearchResult(c *gin.Context) {
	var searchResult model.SearchResult
	_ = c.ShouldBindJSON(&searchResult)
	if err := service.CreateSearchResult(searchResult); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags SearchResult
// @Summary 删除SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "删除SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /searchResult/deleteSearchResult [delete]
func DeleteSearchResult(c *gin.Context) {
	var searchResult model.SearchResult
	_ = c.ShouldBindJSON(&searchResult)
	if err := service.DeleteSearchResult(searchResult); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags SearchResult
// @Summary 批量删除SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /searchResult/deleteSearchResultByIds [delete]
func DeleteSearchResultByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteSearchResultByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

func UpdateSearchResultByIds(c *gin.Context) {
	var batchUpdateReq request.BatchUpdateReq
	_ = c.ShouldBindJSON(&batchUpdateReq)
	if err := service.UpdateSearchResultByIds(batchUpdateReq); err != nil {
		global.GVA_LOG.Error("批量更新状态失败！", zap.Any("err", err))
		response.FailWithMessage("批量更新状态失败", c)
	} else {
		response.OkWithMessage("批量更新状态成功", c)
	}
}

// @Tags SearchResult
// @Summary 更新SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "更新SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /searchResult/updateSearchResult [put]
func UpdateSearchResult(c *gin.Context) {
	var updateReq request.UpdateReq
	_ = c.ShouldBindJSON(&updateReq)
	if err := service.UpdateSearchResult(updateReq); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags SearchResult
// @Summary 用id查询SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "用id查询SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /searchResult/findSearchResult [get]
func FindSearchResult(c *gin.Context) {
	var searchResult model.SearchResult
	_ = c.ShouldBindQuery(&searchResult)
	if err, researchResult := service.GetSearchResult(searchResult.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"researchResult": researchResult}, c)
	}
}

// @Tags SearchResult
// @Summary 分页获取SearchResult列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SearchResultSearch true "分页获取SearchResult列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /searchResult/getSearchResultList [get]
func GetSearchResultList(c *gin.Context) {
	var pageInfo request.SearchResultSearch
	_ = c.ShouldBindQuery(&pageInfo)
	//if err != nil {
	//	global.GVA_LOG.Error("GetSearchResultList bind query error", zap.Any("err", err))
	//}
	if err, list, total := service.GetSearchResultInfoList(pageInfo); err != nil {
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

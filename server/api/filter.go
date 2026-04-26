package api

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
)

func CreateFilter(c *gin.Context) {
	var filter model.Filter
	if !bindJSON(c, &filter) {
		return
	}
	respondMutation(c, service.CreateFilter(filter), "创建失败!", "创建失败", "创建成功")
}

func DeleteFilter(c *gin.Context) {
	var filter model.Filter
	if !bindJSON(c, &filter) {
		return
	}
	respondMutation(c, service.DeleteFilter(filter), "删除失败!", "删除失败", "删除成功")
}

func DeleteFilterByIds(c *gin.Context) {
	var IDS request.IdsReq
	if !bindJSON(c, &IDS) {
		return
	}
	respondMutation(c, service.DeleteFilterByIds(IDS), "批量删除失败!", "批量删除失败", "批量删除成功")
}

func UpdateFilter(c *gin.Context) {
	var filter model.Filter
	if !bindJSON(c, &filter) {
		return
	}
	respondMutation(c, service.UpdateFilter(filter), "更新失败!", "更新失败", "更新成功")
}

func FindFilter(c *gin.Context) {
	var filter model.Filter
	if !bindQuery(c, &filter) {
		return
	}
	if err, filter := service.GetFilter(filter.ID); err != nil {
		respondMutation(c, err, "查询失败!", "查询失败", "")
	} else {
		response.OkWithData(gin.H{"filter": filter}, c)
	}
}

func GetFilterList(c *gin.Context) {
	var pageInfo request.FilterSearch
	if !bindQuery(c, &pageInfo) {
		return
	}
	err, list, total := service.GetFilterInfoList(pageInfo)
	respondPage(c, err, list, total, pageInfo.Page, pageInfo.PageSize)
}

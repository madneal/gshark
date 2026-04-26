package api

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
)

func CreateSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	if !bindJSON(c, &subdomain) {
		return
	}
	respondMutation(c, service.CreateSubdomain(subdomain), "创建失败!", "创建失败", "创建成功")
}

func DeleteSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	if !bindJSON(c, &subdomain) {
		return
	}
	respondMutation(c, service.DeleteSubdomain(subdomain), "删除失败!", "删除失败", "删除成功")
}

func DeleteSubdomainByIds(c *gin.Context) {
	var IDS request.IdsReq
	if !bindJSON(c, &IDS) {
		return
	}
	respondMutation(c, service.DeleteSubdomainByIds(IDS), "批量删除失败!", "批量删除失败", "批量删除成功")
}

func UpdateSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	if !bindJSON(c, &subdomain) {
		return
	}
	respondMutation(c, service.UpdateSubdomain(subdomain), "更新失败!", "更新失败", "更新成功")
}

func FindSubdomain(c *gin.Context) {
	var subdomain model.Subdomain
	if !bindQuery(c, &subdomain) {
		return
	}
	if err, resubdomain := service.GetSubdomain(subdomain.ID); err != nil {
		respondMutation(c, err, "查询失败!", "查询失败", "")
	} else {
		response.OkWithData(gin.H{"resubdomain": resubdomain}, c)
	}
}

func GetSubdomainList(c *gin.Context) {
	var pageInfo request.SubdomainSearch
	if !bindQuery(c, &pageInfo) {
		return
	}
	err, list, total := service.GetSubdomainInfoList(pageInfo)
	respondPage(c, err, list, total, pageInfo.Page, pageInfo.PageSize)
}

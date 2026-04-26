package api

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
)

func CreateToken(c *gin.Context) {
	var token model.Token
	if !bindJSON(c, &token) {
		return
	}
	respondMutation(c, service.CreateToken(token), "创建失败!", "创建失败", "创建成功")
}

func DeleteToken(c *gin.Context) {
	var token model.Token
	if !bindJSON(c, &token) {
		return
	}
	respondMutation(c, service.DeleteToken(token), "删除失败!", "删除失败", "删除成功")
}

func DeleteTokenByIds(c *gin.Context) {
	var IDS request.IdsReq
	if !bindJSON(c, &IDS) {
		return
	}
	respondMutation(c, service.DeleteTokenByIds(IDS), "批量删除失败!", "批量删除失败", "批量删除成功")
}

func UpdateToken(c *gin.Context) {
	var token model.Token
	if !bindJSON(c, &token) {
		return
	}
	respondMutation(c, service.UpdateToken(token), "更新失败!", "更新失败", "更新成功")
}

func FindToken(c *gin.Context) {
	var token model.Token
	if !bindQuery(c, &token) {
		return
	}
	if err, retoken := service.GetToken(token.ID); err != nil {
		respondMutation(c, err, "查询失败!", "查询失败", "")
	} else {
		response.OkWithData(gin.H{"retoken": retoken}, c)
	}
}

func GetTokenList(c *gin.Context) {
	var pageInfo request.TokenSearch
	if !bindQuery(c, &pageInfo) {
		return
	}
	err, list, total := service.GetTokenInfoList(pageInfo)
	respondPage(c, err, list, total, pageInfo.Page, pageInfo.PageSize)
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
)

func CreateRepo(c *gin.Context) {
	var repo model.Repo
	if !bindJSON(c, &repo) {
		return
	}
	respondMutation(c, service.CreateRepo(repo), "创建失败!", "创建失败", "创建成功")
}

func DeleteRepo(c *gin.Context) {
	var repo model.Repo
	if !bindJSON(c, &repo) {
		return
	}
	respondMutation(c, service.DeleteRepo(repo), "删除失败!", "删除失败", "删除成功")
}

func DeleteRepoByIds(c *gin.Context) {
	var IDS request.IdsReq
	if !bindJSON(c, &IDS) {
		return
	}
	respondMutation(c, service.DeleteRepoByIds(IDS), "批量删除失败!", "批量删除失败", "批量删除成功")
}

func UpdateRepo(c *gin.Context) {
	var repo model.Repo
	if !bindJSON(c, &repo) {
		return
	}
	respondMutation(c, service.UpdateRepo(repo), "更新失败!", "更新失败", "更新成功")
}

func FindRepo(c *gin.Context) {
	var repo model.Repo
	if !bindQuery(c, &repo) {
		return
	}
	if err, repo := service.GetRepo(repo.ID); err != nil {
		respondMutation(c, err, "查询失败!", "查询失败", "")
	} else {
		response.OkWithData(gin.H{"repo": repo}, c)
	}
}

func GetRepoList(c *gin.Context) {
	var pageInfo request.RepoSearch
	if !bindQuery(c, &pageInfo) {
		return
	}
	err, list, total := service.GetRepoInfoList(pageInfo)
	respondPage(c, err, list, total, pageInfo.Page, pageInfo.PageSize)
}

package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/search/githubsearch"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"strings"
)

var taskStatus = "stop"

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

func GetTaskStatus(c *gin.Context) {
	response.OkWithMessage(taskStatus, c)
}

func StartSecFilterTask(c *gin.Context) {
	response.Ok(c)
	go func(taskStatus *string) {
		*taskStatus = "running"
		client, err := githubsearch.GetGithubClient()
		if err != nil {
			*taskStatus = "failed"
			global.GVA_LOG.Error("GetGithubClient error", zap.Error(err))
			response.FailWithMessage("初始化 github 客户端失败", c)
			return
		}
		err, repos := service.GetReposByStatus(0)
		if err != nil {
			*taskStatus = "failed"
			global.GVA_LOG.Error("GetReposByStatus error", zap.Error(err))
			return
		}
		err, secKeywordFilters := model.GetFilterByClass("sec_keyword")
		if err != nil {
			*taskStatus = "failed"
			global.GVA_LOG.Error("GetFilterByClass sec_keyword error", zap.Error(err))
			return
		}
		var secKeywords []string
		for _, secKeywordFilter := range secKeywordFilters {
			secKeywords = append(secKeywords, strings.Split(secKeywordFilter.Content, ",")...)
		}
		for _, repo := range repos {
			for _, keyword := range secKeywords {
				query := fmt.Sprintf("repo:%s %s ", repo, keyword)
				results, err := client.SearchCode(query)
				// find results after second filter, then ignore the results by repo
				if len(results) > 0 && *results[0].Total > 0 {
					err = service.IgnoreResultsByRepo(repo)
					if err != nil {
						global.GVA_LOG.Error("IgnoreResultsByRepo error", zap.Error(err))
						continue
					}
				}
				originalKeyword, err := service.GetKeywordByRepo(repo)
				if err != nil {
					global.GVA_LOG.Error("GetKeywordByRepo error", zap.Error(err))
					continue
				}
				if err != nil {
					global.GVA_LOG.Error("Github search code error", zap.Error(err))
					continue
				}
				if results != nil && len(results) > 0 && *results[0].Total > 0 {
					githubsearch.SaveResult(results, originalKeyword, keyword)
				}
			}

		}
		*taskStatus = "done"
	}(&taskStatus)

}

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

func FindSearchResult(c *gin.Context) {
	var searchResult model.SearchResult
	_ = c.ShouldBindQuery(&searchResult)
	if err, searchResult := service.GetSearchResult(searchResult.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"searchResult": searchResult}, c)
	}
}

func GetSearchResultList(c *gin.Context) {
	var pageInfo request.SearchResultSearch
	_ = c.ShouldBindQuery(&pageInfo)
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

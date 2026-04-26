package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/search/githubsearch"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

var taskStatus = "stop"
var statusOptions = map[int]string{
	0: "未处理", // Unprocessed
	1: "已处理", // Processed
	2: "已忽略", // Ignored
}

func CreateSearchResult(c *gin.Context) {
	var searchResult model.SearchResult
	if !bindJSON(c, &searchResult) {
		return
	}
	respondMutation(c, service.CreateSearchResult(searchResult), "创建失败!", "创建失败", "创建成功")
}

func DeleteSearchResult(c *gin.Context) {
	var searchResult model.SearchResult
	if !bindJSON(c, &searchResult) {
		return
	}
	respondMutation(c, service.DeleteSearchResult(searchResult), "删除失败!", "删除失败", "删除成功")
}

func DeleteSearchResultByIds(c *gin.Context) {
	var IDS request.IdsReq
	if !bindJSON(c, &IDS) {
		return
	}
	respondMutation(c, service.DeleteSearchResultByIds(IDS), "批量删除失败!", "批量删除失败", "批量删除成功")
}

func UpdateSearchResultByIds(c *gin.Context) {
	var batchUpdateReq request.BatchUpdateReq
	if !bindJSON(c, &batchUpdateReq) {
		return
	}
	respondMutation(c, service.UpdateSearchResultByIds(batchUpdateReq), "批量更新状态失败！", "批量更新状态失败", "批量更新状态成功")
}

func GetTaskStatus(c *gin.Context) {
	response.OkWithMessage(taskStatus, c)
}

func StartAITask(c *gin.Context) {
	response.Ok(c)
	go func() {
		err, list := service.ListSearchResultByStatus(0)
		if err != nil {
			global.GVA_LOG.Error("ListSearchResultByStatus error", zap.Any("err", err))
			return
		}
		for _, result := range list {
			textMatches := make([]model.TextMatch, 0)
			var content string
			if json.Valid(result.TextMatchesJson) {
				err = json.Unmarshal(result.TextMatchesJson, &textMatches)
				if err != nil {
					global.GVA_LOG.Error("json unmarshal error", zap.Error(err), zap.Any("result", result))
					continue
				}
				for _, textMatch := range textMatches {
					content += *textMatch.Fragment + "\n"
				}
			} else {
				content = string(result.TextMatchesJson)
			}
			if content == "" {
				content = result.Matches
			}

			ans := service.Question("You are a security operation engineer, you are expected to assistant."+
				"please judge if the below content contains sensitive information,including password, credentials,token,etc. "+
				"The sensitive information could be exploited. Just answer yes or no",
				content)
			global.GVA_LOG.Info(strconv.Itoa(int(result.ID)))
			global.GVA_LOG.Info(content)
			global.GVA_LOG.Info(ans)
			if strings.ToLower(ans) == "yes" {
				err = service.UpdateSearchResultById(int(result.ID), 1)
			} else {
				err = service.UpdateSearchResultById(int(result.ID), 2)
			}
			if err != nil {
				global.GVA_LOG.Error("UpdateSearchResultByIds error", zap.Any("err", err))
			}
		}
	}()
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
	if !bindJSON(c, &updateReq) {
		return
	}
	respondMutation(c, service.UpdateSearchResult(updateReq), "更新失败!", "更新失败", "更新成功")
}

func FindSearchResult(c *gin.Context) {
	var searchResult model.SearchResult
	if !bindQuery(c, &searchResult) {
		return
	}
	if err, searchResult := service.GetSearchResult(searchResult.ID); err != nil {
		respondMutation(c, err, "查询失败!", "查询失败", "")
	} else {
		response.OkWithData(gin.H{"searchResult": searchResult}, c)
	}
}

func GetSearchResultList(c *gin.Context) {
	var pageInfo request.SearchResultSearch
	if !bindQuery(c, &pageInfo) {
		return
	}
	err, list, total := service.GetSearchResultInfoList(pageInfo)
	respondPage(c, err, list, total, pageInfo.Page, pageInfo.PageSize)
}

func ExportSearchResult(c *gin.Context) {
	var searchInfo request.SearchResultSearch
	if !bindQuery(c, &searchInfo) {
		return
	}
	searchInfo.PageInfo.Page = 1
	searchInfo.PageInfo.PageSize = 10000
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="search_results.csv"`)
	writer := csv.NewWriter(c.Writer)
	headers := []string{"Repo", "RepoUrl", "Matches", "Keyword", "SecKeyword", "Path",
		"Url", "Status"}
	if err := writer.Write(headers); err != nil {
		response.FailWithMessage("导出失败", c)
		return
	}
	err, list, _ := service.GetSearchResultInfoList(searchInfo)
	if err != nil {
		global.GVA_LOG.Error("GetSearchResultInfoList  err", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	searchResults, _ := list.([]model.SearchResult)
	for _, result := range searchResults {
		row := []string{
			result.Repo,
			result.RepoUrl,
			result.Matches,
			result.Keyword,
			result.SecKeyword,
			result.Path,
			result.Url,
			statusOptions[result.Status],
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to write CSV row",
			})
			return
		}
	}
	writer.Flush()
	if err = writer.Error(); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
}

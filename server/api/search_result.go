package api

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"net/http"
)

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

// StartAITask is kept for clients using the old endpoint. AI triage now runs
// before persistence, so existing rows are intentionally not reprocessed here.
func StartAITask(c *gin.Context) {
	response.FailWithMessage("AI 分析已改为搜索结果入库前执行，请在 system.ai_analysis_enabled 中启用", c)
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
	headers := []string{"Repo", "RepoUrl", "Matches", "Keyword", "Path",
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

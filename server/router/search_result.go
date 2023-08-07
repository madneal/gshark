package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitSearchResultRouter(Router *gin.RouterGroup) {
	SearchResultRouter := Router.Group("searchResult").Use(middleware.OperationRecord())
	{
		SearchResultRouter.POST("createSearchResult", api.CreateSearchResult)             // 新建SearchResult
		SearchResultRouter.DELETE("deleteSearchResult", api.DeleteSearchResult)           // 删除SearchResult
		SearchResultRouter.DELETE("deleteSearchResultByIds", api.DeleteSearchResultByIds) // 批量删除SearchResult
		SearchResultRouter.POST("updateSearchResult", api.UpdateSearchResult)             // 更新SearchResult
		SearchResultRouter.GET("findSearchResult", api.FindSearchResult)                  // 根据ID获取SearchResult
		SearchResultRouter.GET("getSearchResultList", api.GetSearchResultList)            // 获取SearchResult列表
		SearchResultRouter.POST("updateSearchResultStatusByIds", api.UpdateSearchResultByIds)
		SearchResultRouter.POST("startSecFilterTask", api.StartSecFilterTask)
		SearchResultRouter.GET("getTaskStatus", api.GetTaskStatus)
	}
}

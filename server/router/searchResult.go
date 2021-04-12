package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
)

func InitSearchResultRouter(Router *gin.RouterGroup) {
	SearchResultRouter := Router.Group("searchResult").Use(middleware.OperationRecord())
	{
		SearchResultRouter.POST("createSearchResult", v1.CreateSearchResult)             // 新建SearchResult
		SearchResultRouter.DELETE("deleteSearchResult", v1.DeleteSearchResult)           // 删除SearchResult
		SearchResultRouter.DELETE("deleteSearchResultByIds", v1.DeleteSearchResultByIds) // 批量删除SearchResult
		SearchResultRouter.POST("updateSearchResult", v1.UpdateSearchResult)             // 更新SearchResult
		SearchResultRouter.GET("findSearchResult", v1.FindSearchResult)                  // 根据ID获取SearchResult
		SearchResultRouter.GET("getSearchResultList", v1.GetSearchResultList)            // 获取SearchResult列表
		SearchResultRouter.POST("updateSearchResultStatusByIds", v1.UpdateSearchResultByIds)
	}
}

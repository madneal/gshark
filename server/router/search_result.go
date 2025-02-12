package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitSearchResultRouter(Router *gin.RouterGroup) {
	SearchResultRouter := Router.Group("searchResult").Use(middleware.OperationRecord())
	{
		SearchResultRouter.POST("createSearchResult", api.CreateSearchResult)
		SearchResultRouter.DELETE("deleteSearchResult", api.DeleteSearchResult)
		SearchResultRouter.DELETE("deleteSearchResultByIds", api.DeleteSearchResultByIds)
		SearchResultRouter.POST("updateSearchResult", api.UpdateSearchResult)
		SearchResultRouter.GET("findSearchResult", api.FindSearchResult)
		SearchResultRouter.GET("getSearchResultList", api.GetSearchResultList)
		SearchResultRouter.GET("exportSearchResult", api.ExportSearchResult)
		SearchResultRouter.POST("updateSearchResultStatusByIds", api.UpdateSearchResultByIds)
		SearchResultRouter.POST("startSecFilterTask", api.StartSecFilterTask)
		SearchResultRouter.GET("getTaskStatus", api.GetTaskStatus)
		SearchResultRouter.POST("startAITask", api.StartAITask)
	}
}

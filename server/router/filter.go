package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitFilterRouter(Router *gin.RouterGroup) {
	FilterRouter := Router.Group("filter").Use(middleware.OperationRecord())
	{
		FilterRouter.POST("createFilter", api.CreateFilter)             // 新建Filter
		FilterRouter.DELETE("deleteFilter", api.DeleteFilter)           // 删除Filter
		FilterRouter.DELETE("deleteFilterByIds", api.DeleteFilterByIds) // 批量删除Filter
		FilterRouter.PUT("updateFilter", api.UpdateFilter)              // 更新Filter
		FilterRouter.GET("findFilter", api.FindFilter)                  // 根据ID获取Filter
		FilterRouter.GET("getFilterList", api.GetFilterList)            // 获取Filter列表
	}
}

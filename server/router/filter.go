package router

import (
	"github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
	"github.com/gin-gonic/gin"
)

func InitFilterRouter(Router *gin.RouterGroup) {
	FilterRouter := Router.Group("filter").Use(middleware.OperationRecord())
	{
		FilterRouter.POST("createFilter", v1.CreateFilter)   // 新建Filter
		FilterRouter.DELETE("deleteFilter", v1.DeleteFilter) // 删除Filter
		FilterRouter.DELETE("deleteFilterByIds", v1.DeleteFilterByIds) // 批量删除Filter
		FilterRouter.PUT("updateFilter", v1.UpdateFilter)    // 更新Filter
		FilterRouter.GET("findFilter", v1.FindFilter)        // 根据ID获取Filter
		FilterRouter.GET("getFilterList", v1.GetFilterList)  // 获取Filter列表
	}
}

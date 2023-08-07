package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitSysDictionaryDetailRouter(Router *gin.RouterGroup) {
	SysDictionaryDetailRouter := Router.Group("sysDictionaryDetail").Use(middleware.OperationRecord())
	{
		SysDictionaryDetailRouter.POST("createSysDictionaryDetail", api.CreateSysDictionaryDetail)   // 新建SysDictionaryDetail
		SysDictionaryDetailRouter.DELETE("deleteSysDictionaryDetail", api.DeleteSysDictionaryDetail) // 删除SysDictionaryDetail
		SysDictionaryDetailRouter.PUT("updateSysDictionaryDetail", api.UpdateSysDictionaryDetail)    // 更新SysDictionaryDetail
		SysDictionaryDetailRouter.GET("findSysDictionaryDetail", api.FindSysDictionaryDetail)        // 根据ID获取SysDictionaryDetail
		SysDictionaryDetailRouter.GET("getSysDictionaryDetailList", api.GetSysDictionaryDetailList)  // 获取SysDictionaryDetail列表
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitSysDictionaryRouter(Router *gin.RouterGroup) {
	SysDictionaryRouter := Router.Group("sysDictionary").Use(middleware.OperationRecord())
	{
		SysDictionaryRouter.POST("createSysDictionary", api.CreateSysDictionary)   // 新建SysDictionary
		SysDictionaryRouter.DELETE("deleteSysDictionary", api.DeleteSysDictionary) // 删除SysDictionary
		SysDictionaryRouter.PUT("updateSysDictionary", api.UpdateSysDictionary)    // 更新SysDictionary
		SysDictionaryRouter.GET("findSysDictionary", api.FindSysDictionary)        // 根据ID获取SysDictionary
		SysDictionaryRouter.GET("getSysDictionaryList", api.GetSysDictionaryList)  // 获取SysDictionary列表
	}
}

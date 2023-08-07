package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitApiRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("api").Use(middleware.OperationRecord())
	{
		ApiRouter.POST("createApi", api.CreateApi)   // 创建Api
		ApiRouter.POST("deleteApi", api.DeleteApi)   // 删除Api
		ApiRouter.POST("getApiList", api.GetApiList) // 获取Api列表
		ApiRouter.POST("getApiById", api.GetApiById) // 获取单条Api消息
		ApiRouter.POST("updateApi", api.UpdateApi)   // 更新api
		ApiRouter.POST("getAllApis", api.GetAllApis) // 获取所有api
	}
}

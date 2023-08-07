package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
)

func InitAutoCodeRouter(Router *gin.RouterGroup) {
	AutoCodeRouter := Router.Group("autoCode")
	{
		AutoCodeRouter.POST("preview", api.PreviewTemp)   // 获取自动创建代码预览
		AutoCodeRouter.POST("createTemp", api.CreateTemp) // 创建自动化代码
		AutoCodeRouter.GET("getTables", api.GetTables)    // 获取对应数据库的表
		AutoCodeRouter.GET("getDB", api.GetDB)            // 获取数据库
		AutoCodeRouter.GET("getColumn", api.GetColumn)    // 获取指定表所有字段信息
	}
}

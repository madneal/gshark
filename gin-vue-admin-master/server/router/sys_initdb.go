package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
)

func InitInitRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("init")
	{
		ApiRouter.POST("initdb", v1.InitDB)   // 创建Api
		ApiRouter.POST("checkdb", v1.CheckDB) // 创建Api
	}
}

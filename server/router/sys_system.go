package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitSystemRouter(Router *gin.RouterGroup) {
	SystemRouter := Router.Group("system").Use(middleware.OperationRecord())
	{
		SystemRouter.POST("getSystemConfig", api.GetSystemConfig)
		SystemRouter.POST("setSystemConfig", api.SetSystemConfig)
		SystemRouter.POST("getServerInfo", api.GetServerInfo)
		SystemRouter.POST("reloadSystem", api.ReloadSystem)
		SystemRouter.POST("emailTest", api.EmailTest) // 发送测试邮件
		SystemRouter.GET("botTest", api.BotTest)
	}
}

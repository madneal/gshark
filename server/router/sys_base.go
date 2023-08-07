package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("base").Use(middleware.NeedInit())
	{
		BaseRouter.POST("login", api.Login)
		BaseRouter.POST("captcha", api.Captcha)
	}
	return BaseRouter
}

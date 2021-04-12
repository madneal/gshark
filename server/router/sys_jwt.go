package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
)

func InitJwtRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("jwt").Use(middleware.OperationRecord())
	{
		ApiRouter.POST("jsonInBlacklist", v1.JsonInBlacklist) // jwt加入黑名单
	}
}

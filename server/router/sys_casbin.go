package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitCasbinRouter(Router *gin.RouterGroup) {
	CasbinRouter := Router.Group("casbin").Use(middleware.OperationRecord())
	{
		CasbinRouter.POST("updateCasbin", api.UpdateCasbin)
		CasbinRouter.POST("getPolicyPathByAuthorityId", api.GetPolicyPathByAuthorityId)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitTokenRouter(Router *gin.RouterGroup) {
	TokenRouter := Router.Group("token").Use(middleware.OperationRecord())
	{
		TokenRouter.POST("createToken", api.CreateToken)             // 新建Token
		TokenRouter.DELETE("deleteToken", api.DeleteToken)           // 删除Token
		TokenRouter.DELETE("deleteTokenByIds", api.DeleteTokenByIds) // 批量删除Token
		TokenRouter.PUT("updateToken", api.UpdateToken)              // 更新Token
		TokenRouter.GET("findToken", api.FindToken)                  // 根据ID获取Token
		TokenRouter.GET("getTokenList", api.GetTokenList)            // 获取Token列表
	}
}

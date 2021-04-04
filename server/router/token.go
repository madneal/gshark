package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
)

func InitTokenRouter(Router *gin.RouterGroup) {
	TokenRouter := Router.Group("token").Use(middleware.OperationRecord())
	{
		TokenRouter.POST("createToken", v1.CreateToken)             // 新建Token
		TokenRouter.DELETE("deleteToken", v1.DeleteToken)           // 删除Token
		TokenRouter.DELETE("deleteTokenByIds", v1.DeleteTokenByIds) // 批量删除Token
		TokenRouter.PUT("updateToken", v1.UpdateToken)              // 更新Token
		TokenRouter.GET("findToken", v1.FindToken)                  // 根据ID获取Token
		TokenRouter.GET("getTokenList", v1.GetTokenList)            // 获取Token列表
	}
}

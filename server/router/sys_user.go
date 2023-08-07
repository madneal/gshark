package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(middleware.OperationRecord())
	{
		UserRouter.POST("register", api.Register)
		UserRouter.POST("changePassword", api.ChangePassword)     // 修改密码
		UserRouter.POST("getUserList", api.GetUserList)           // 分页获取用户列表
		UserRouter.POST("setUserAuthority", api.SetUserAuthority) // 设置用户权限
		UserRouter.DELETE("deleteUser", api.DeleteUser)           // 删除用户
		UserRouter.PUT("setUserInfo", api.SetUserInfo)            // 设置用户信息
	}
}

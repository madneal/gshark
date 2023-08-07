package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitSubdomainRouter(Router *gin.RouterGroup) {
	SubdomainRouter := Router.Group("subdomain").Use(middleware.OperationRecord())
	{
		SubdomainRouter.POST("createSubdomain", api.CreateSubdomain)             // 新建Subdomain
		SubdomainRouter.DELETE("deleteSubdomain", api.DeleteSubdomain)           // 删除Subdomain
		SubdomainRouter.DELETE("deleteSubdomainByIds", api.DeleteSubdomainByIds) // 批量删除Subdomain
		SubdomainRouter.PUT("updateSubdomain", api.UpdateSubdomain)              // 更新Subdomain
		SubdomainRouter.GET("findSubdomain", api.FindSubdomain)                  // 根据ID获取Subdomain
		SubdomainRouter.GET("getSubdomainList", api.GetSubdomainList)            // 获取Subdomain列表
	}
}

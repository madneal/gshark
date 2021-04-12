package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
)

func InitSubdomainRouter(Router *gin.RouterGroup) {
	SubdomainRouter := Router.Group("subdomain").Use(middleware.OperationRecord())
	{
		SubdomainRouter.POST("createSubdomain", v1.CreateSubdomain)             // 新建Subdomain
		SubdomainRouter.DELETE("deleteSubdomain", v1.DeleteSubdomain)           // 删除Subdomain
		SubdomainRouter.DELETE("deleteSubdomainByIds", v1.DeleteSubdomainByIds) // 批量删除Subdomain
		SubdomainRouter.PUT("updateSubdomain", v1.UpdateSubdomain)              // 更新Subdomain
		SubdomainRouter.GET("findSubdomain", v1.FindSubdomain)                  // 根据ID获取Subdomain
		SubdomainRouter.GET("getSubdomainList", v1.GetSubdomainList)            // 获取Subdomain列表
	}
}

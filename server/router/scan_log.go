package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
)

func InitScanLogRouter(router *gin.RouterGroup) {
	scanLogRouter := router.Group("scanLog")
	{
		scanLogRouter.GET("getScanLogList", api.GetScanLogList)
		scanLogRouter.GET("getScanLogOverview", api.GetScanLogOverview)
	}
}

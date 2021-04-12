package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
)

func InitExcelRouter(Router *gin.RouterGroup) {
	FileUploadAndDownloadGroup := Router.Group("excel")
	{
		FileUploadAndDownloadGroup.POST("/importExcel", v1.ImportExcel)          // 导入Excel
		FileUploadAndDownloadGroup.GET("/loadExcel", v1.LoadExcel)               // 加载Excel数据
		FileUploadAndDownloadGroup.POST("/exportExcel", v1.ExportExcel)          // 导出Excel
		FileUploadAndDownloadGroup.GET("/downloadTemplate", v1.DownloadTemplate) // 下载模板文件
	}
}

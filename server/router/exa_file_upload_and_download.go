package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
)

func InitFileUploadAndDownloadRouter(Router *gin.RouterGroup) {
	FileUploadAndDownloadGroup := Router.Group("fileUploadAndDownload")
	{
		FileUploadAndDownloadGroup.POST("/upload", v1.UploadFile)                                 // 上传文件
		FileUploadAndDownloadGroup.POST("/getFileList", v1.GetFileList)                           // 获取上传文件列表
		FileUploadAndDownloadGroup.POST("/deleteFile", v1.DeleteFile)                             // 删除指定文件
		FileUploadAndDownloadGroup.POST("/breakpointContinue", v1.BreakpointContinue)             // 断点续传
		FileUploadAndDownloadGroup.GET("/findFile", v1.FindFile)                                  // 查询当前文件成功的切片
		FileUploadAndDownloadGroup.POST("/breakpointContinueFinish", v1.BreakpointContinueFinish) // 查询当前文件成功的切片
		FileUploadAndDownloadGroup.POST("/removeChunk", v1.RemoveChunk)                           // 查询当前文件成功的切片
	}
}

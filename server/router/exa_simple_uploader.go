package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
)

func InitSimpleUploaderRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("simpleUploader")
	{
		ApiRouter.POST("upload", v1.SimpleUploaderUpload) // 上传功能
		ApiRouter.GET("checkFileMd5", v1.CheckFileMd5)    // 文件完整度验证
		ApiRouter.GET("mergeFileMd5", v1.MergeFileMd5)    // 合并文件
	}
}

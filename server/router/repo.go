package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitRepoRouter(Router *gin.RouterGroup) {
	RepoRouter := Router.Group("repo").Use(middleware.OperationRecord())
	{
		RepoRouter.POST("createRepo", api.CreateRepo)             // 新建Repo
		RepoRouter.DELETE("deleteRepo", api.DeleteRepo)           // 删除Repo
		RepoRouter.DELETE("deleteRepoByIds", api.DeleteRepoByIds) // 批量删除Repo
		RepoRouter.PUT("updateRepo", api.UpdateRepo)              // 更新Repo
		RepoRouter.GET("findRepo", api.FindRepo)                  // 根据ID获取Repo
		RepoRouter.GET("getRepoList", api.GetRepoList)            // 获取Repo列表
	}
}

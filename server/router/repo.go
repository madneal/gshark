package router

import (
	"github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
	"github.com/gin-gonic/gin"
)

func InitRepoRouter(Router *gin.RouterGroup) {
	RepoRouter := Router.Group("repo").Use(middleware.OperationRecord())
	{
		RepoRouter.POST("createRepo", v1.CreateRepo)   // 新建Repo
		RepoRouter.DELETE("deleteRepo", v1.DeleteRepo) // 删除Repo
		RepoRouter.DELETE("deleteRepoByIds", v1.DeleteRepoByIds) // 批量删除Repo
		RepoRouter.PUT("updateRepo", v1.UpdateRepo)    // 更新Repo
		RepoRouter.GET("findRepo", v1.FindRepo)        // 根据ID获取Repo
		RepoRouter.GET("getRepoList", v1.GetRepoList)  // 获取Repo列表
	}
}

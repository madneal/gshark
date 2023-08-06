package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
)

func InitTaskRouter(Router *gin.RouterGroup) {
	TaskRouter := Router.Group("task").Use(middleware.OperationRecord())
	{
		TaskRouter.GET("getTaskList", v1.GetTaskList)
		TaskRouter.POST("createTask", v1.CreateTask)
		TaskRouter.POST("switchTaskStatus", v1.SwitchTaskStatus)
	}
}

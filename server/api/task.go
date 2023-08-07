package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateTask(c *gin.Context) {
	var req model.CreateTaskReq
	_ = c.ShouldBindJSON(&req)
	task := model.Task{
		TaskStatus: req.TaskStatus,
		TaskType:   req.TaskType,
		Name:       req.Name,
	}
	if err := service.CreateTask(&task); err != nil {
		global.GVA_LOG.Error("CreateTask err", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailWithMessage(fmt.Sprintf("暂无%s类型规则，请至少创建一条规则", task.TaskType), c)
		} else {
			response.FailWithMessage("创建失败", c)
		}
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func GetTaskList(c *gin.Context) {
	if err, list, total := service.GetTaskList(); err != nil {
		global.GVA_LOG.Error("获取失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     1,
			PageSize: 10,
		}, "获取成功", c)
	}
}

func SwitchTaskStatus(c *gin.Context) {
	var switchTaskReq request.SwitchTaskReq
	_ = c.ShouldBindJSON(&switchTaskReq)
	if err := service.SwitchTaskStatus(switchTaskReq.ID, switchTaskReq.Status); err != nil {
		global.GVA_LOG.Error("切换状态失败", zap.Error(err))
		response.FailWithMessage("切换状态失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

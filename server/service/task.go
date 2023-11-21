package service

import (
	"errors"
	"fmt"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"gorm.io/gorm"
)

func CreateTask(task *model.Task) error {
	hasRule, err := CheckRuleByType(task.TaskType)
	if err != nil {
		return err
	}
	if !hasRule {
		err = errors.New(fmt.Sprintf("暂无%s类型规则，请至少创建一条规则", task.TaskType))
		return err
	}
	err = global.GVA_DB.Create(task).Error
	return err
}

func GetTaskList() (error, interface{}, int64) {
	db := global.GVA_DB.Table("task")
	var tasks []model.Task
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return err, tasks, total
	}
	err = db.Limit(10).Offset(0).Find(&tasks).Error
	return err, tasks, total
}

func SwitchTaskStatus(id int, status int) error {
	err := global.GVA_DB.Debug().Table("task").Where("id = ?", id).
		UpdateColumn("task_status", status).Error
	return err
}

func CheckTaskStatus(taskType string) (bool, error) {
	var task model.Task
	result := global.GVA_DB.Table("task").Where("task_type = ? and task_status = 1", taskType).Take(&task)
	err := result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result.RowsAffected > 0, err
}

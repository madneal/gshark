package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
)

func CreateTask(task *model.Task) error {
	err := global.GVA_DB.Create(task).Error
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

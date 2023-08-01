package model

import "github.com/madneal/gshark/global"

type CreateTaskReq struct {
	Name       string `json:"name"`
	TaskType   string `json:"taskType"`
	TaskStatus bool   `json:"taskStatus"`
}

type Task struct {
	global.GVA_MODEL
	TaskType   string `json:"taskType" form:"task_type" gorm:"column:task_type;type:varchar(10);"`
	Name       string `json:"name" form:"string" gorm:"column:name;type:varchar(20);"`
	TaskStatus bool   `json:"taskStatus" form:"task_status" gorm:"column:task_status;type:tinyint(1);"`
}

func (Task) TableName() string {
	return "task"
}

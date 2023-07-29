package model

import "github.com/madneal/gshark/global"

type CreateTaskReq struct {
	Name       string
	TaskType   string
	TaskStatus int
}

type Task struct {
	global.GVA_MODEL
	TaskType   string `json:"task_type" form:"task_type" gorm:"column:task_type;type:varchar(10);"`
	Name       string `json:"name" form:"string" gorm:"column:name;type:varchar(20);"`
	TaskStatus int    `json:"task_status" form:"task_status" gorm:"column:task_status;type:tinyint(1);"`
}

func (Task) TableName() string {
	return "task"
}

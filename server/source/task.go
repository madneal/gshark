package source

import (
	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"gorm.io/gorm"
)

var tasks = []model.Task{
	{TaskType: "github", Name: "github task", TaskStatus: true},
	{TaskType: "gitlab", Name: "github task", TaskStatus: true},
	{TaskType: "postman", Name: "github task", TaskStatus: true},
	{TaskType: "searchcode", Name: "github task", TaskStatus: true},
}

func Init() error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("task_type in ?", []string{"github", "gitlab"}).Find(&[]model.Task{}).RowsAffected > 1 {
			color.Danger.Println("\n[Mysql] --> task 表的初始数据已存在!")
		}
		if err := tx.Create(&tasks).Error; err != nil {
			return err
		}
		color.Info.Println("\n[Mysql] --> task 表的初始数据成功!")
		return nil
	})
}

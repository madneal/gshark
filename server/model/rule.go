// 自动生成模板Rule
package model

import (
	"github.com/madneal/gshark/global"
)

// 如果含有time.Time 请自行import time包
type Rule struct {
	global.GVA_MODEL
	Type    string `json:"type" form:"type" gorm:"column:type;comment:;type:varchar(20);size:20;"`
	Content string `json:"content" form:"content" gorm:"column:content;comment:;type:varchar(100);size:100;"`
	Name    string `json:"name" form:"name" gorm:"column:name;comment:;type:varchar(100);size:100;"`
	Desc    string `json:"desc" form:"desc" gorm:"column:desc;comment:;type:varchar(300);size:300;"`
	Status  int    `json:"status" form:"status" gorm:"column:status;comment:;type:int;size:3;"`
}

func (Rule) TableName() string {
	return "rule"
}


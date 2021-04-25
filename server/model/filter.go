// 自动生成模板Filter
package model

import (
	"github.com/madneal/gshark/global"
)

// 如果含有time.Time 请自行import time包
type Filter struct {
	global.GVA_MODEL
	Extension string `json:"extension" form:"extension" gorm:"column:extension;comment:;type:varchar(1000);size:1000;"`
}

func GetFilterRule() (error, Filter) {
	var filter Filter
	err := global.GVA_DB.First(&filter).Error
	return err, filter
}

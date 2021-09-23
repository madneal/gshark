// 自动生成模板Token
package model

import (
	"github.com/madneal/gshark/global"
	"time"
)

// 如果含有time.Time 请自行import time包
type Token struct {
	global.GVA_MODEL
	Type       string    `json:"type" form:"type" gorm:"column:type;comment:;type:varchar(10);size:10;"`
	Content    string    `json:"content" form:"content" gorm:"column:content;comment:;type:varchar(100);size:100;"`
	Desc       string    `json:"desc" form:"desc" gorm:"column:description;comment:;type:varchar(100);size:100;"`
	LimitTimes int       `json:"limit" form:"limit" gorm:"column:limit_times;comment:;type:int(5);"`
	Remaining  int       `json:"remaining" form:"remaining" gorm:"column:remaining;comment:;type:int(5);"`
	ResetTime  time.Time `json:"resetTime" form:"resetTime" gorm:"column:reset_time;comment:"`
}

func (Token) TableName() string {
	return "token"
}

// 自动生成模板Token
package model

import (
	"gin-vue-admin/global"
      "time"
)

// 如果含有time.Time 请自行import time包
type Token struct {
      global.GVA_MODEL
      Type  string `json:"type" form:"type" gorm:"column:type;comment:;type:varchar(10);size:10;"`
      Content  string `json:"content" form:"content" gorm:"column:content;comment:;type:varchar(100);size:100;"`
      Desc  string `json:"desc" form:"desc" gorm:"column:desc;comment:;type:varchar(100);size:100;"`
      Limit  int `json:"limit" form:"limit" gorm:"column:limit;comment:;type:int;size:5;"`
      Remaining  int `json:"remaining" form:"remaining" gorm:"column:remaining;comment:;type:int;size:5;"`
      ResetTime  time.Time `json:"resetTime" form:"resetTime" gorm:"column:reset_time;comment:"`
}


func (Token) TableName() string {
  return "token"
}


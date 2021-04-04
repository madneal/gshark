// 自动生成模板Subdomain
package model

import (
	"github.com/madneal/gshark/global"
)

// 如果含有time.Time 请自行import time包
type Subdomain struct {
	global.GVA_MODEL
	Subdomain string `json:"subdomain" form:"subdomain" gorm:"column:subdomain;comment:;type:varchar(100);size:100;"`
	Domain    string `json:"domain" form:"domain" gorm:"column:domain;comment:;type:varchar(100);size:100;"`
	Status    int    `json:"status" form:"status" gorm:"column:status;comment:;type:int;size:3;"`
}

func (Subdomain) TableName() string {
	return "subdomain"
}

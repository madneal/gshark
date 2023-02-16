package model

import (
	"github.com/madneal/gshark/global"
)

type Subdomain struct {
	global.GVA_MODEL
	Subdomain string `json:"subdomain" form:"subdomain" gorm:"column:subdomain;comment:;type:varchar(100);size:100;"`
	Domain    string `json:"domain" form:"domain" gorm:"column:domain;comment:;type:varchar(100);size:100;"`
	Status    int    `json:"status" form:"status" gorm:"column:status;comment:;type:int;size:3;"`
	Source    string `json:"source" form:"source" gorm:"column:source;comment:;type:varchar(20);size:20;'"`
}

func (Subdomain) TableName() string {
	return "subdomain"
}

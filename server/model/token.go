package model

import (
	"github.com/madneal/gshark/global"
)

type Token struct {
	global.GVA_MODEL
	Type    string `json:"type" form:"type" gorm:"column:type;comment:;type:varchar(10);size:10;"`
	Content string `json:"content" form:"content" gorm:"column:content;comment:;type:varchar(100);size:100;"`
}

func (Token) TableName() string {
	return "token"
}

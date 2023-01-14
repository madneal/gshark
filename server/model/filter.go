package model

import (
	"github.com/madneal/gshark/global"
)

type Filter struct {
	global.GVA_MODEL
	FilterType  string `json:"filter_type" form:"filter_type" gorm:"column:filter_type;comment:;type:varchar(20);"`
	FilterClass string `json:"filter_class" form:"filter_class" gorm:"column:filter_class;type:varchar(20);"`
	Content     string `json:"content" form:"content" gorm:"column:content;type:varchar(100);"`
}

func GetFilterRule() (error, Filter) {
	var filter Filter
	err := global.GVA_DB.First(&filter).Error
	return err, filter
}

func GetFilterByClass(class string) (error, []Filter) {
	filters := make([]Filter, 0)
	err := global.GVA_DB.Where("filter_class = ?", class).Find(&filters).Error
	return err, filters
}

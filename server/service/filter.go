package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

func CreateFilter(filter model.Filter) (err error) {
	err = global.GVA_DB.Create(&filter).Error
	return err
}

func DeleteFilter(filter model.Filter) (err error) {
	err = global.GVA_DB.Delete(&filter).Error
	return err
}

func DeleteFilterByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Filter{}, "id in ?", ids.Ids).Error
	return err
}

func UpdateFilter(filter model.Filter) (err error) {
	err = global.GVA_DB.Save(&filter).Error
	return err
}

func GetFilter(id uint) (err error, filter model.Filter) {
	err = global.GVA_DB.Where("id = ?", id).First(&filter).Error
	return
}

func GetFilterInfoList(info request.FilterSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.Filter{})
	var filters []model.Filter
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&filters).Error
	return err, filters, total
}

package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateFilter
//@description: 创建Filter记录
//@param: filter model.Filter
//@return: err error

func CreateFilter(filter model.Filter) (err error) {
	err = global.GVA_DB.Create(&filter).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFilter
//@description: 删除Filter记录
//@param: filter model.Filter
//@return: err error

func DeleteFilter(filter model.Filter) (err error) {
	err = global.GVA_DB.Delete(&filter).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteFilterByIds
//@description: 批量删除Filter记录
//@param: ids request.IdsReq
//@return: err error

func DeleteFilterByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Filter{},"id in ?",ids.Ids).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateFilter
//@description: 更新Filter记录
//@param: filter *model.Filter
//@return: err error

func UpdateFilter(filter model.Filter) (err error) {
	err = global.GVA_DB.Save(&filter).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetFilter
//@description: 根据id获取Filter记录
//@param: id uint
//@return: err error, filter model.Filter

func GetFilter(id uint) (err error, filter model.Filter) {
	err = global.GVA_DB.Where("id = ?", id).First(&filter).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetFilterInfoList
//@description: 分页获取Filter记录
//@param: info request.FilterSearch
//@return: err error, list interface{}, total int64

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
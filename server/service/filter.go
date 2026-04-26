package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

func CreateFilter(filter model.Filter) (err error) {
	return Create(&filter)
}

func DeleteFilter(filter model.Filter) (err error) {
	return Delete(&filter)
}

func DeleteFilterByIds(ids request.IdsReq) (err error) {
	return DeleteByIds[model.Filter](ids)
}

func UpdateFilter(filter model.Filter) (err error) {
	return Update(&filter)
}

func GetFilter(id uint) (err error, filter model.Filter) {
	filter, err = GetByID[model.Filter](id)
	return
}

func GetFilterInfoList(info request.FilterSearch) (err error, list interface{}, total int64) {
	db := global.GVA_DB.Model(&model.Filter{})
	var filters []model.Filter
	total, err = Paginate(db, info.Page, info.PageSize, &filters, "")
	return err, filters, total
}

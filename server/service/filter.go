package service

import (
	"errors"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

var supportedFilterClasses = map[string]struct{}{
	"extension": {},
	"keyword":   {},
}

func validateFilterClass(class string) error {
	if _, ok := supportedFilterClasses[class]; !ok {
		return errors.New("unsupported filter class")
	}
	return nil
}

func CreateFilter(filter model.Filter) (err error) {
	if err := validateFilterClass(filter.FilterClass); err != nil {
		return err
	}
	return Create(&filter)
}

func DeleteFilter(filter model.Filter) (err error) {
	return Delete(&filter)
}

func DeleteFilterByIds(ids request.IdsReq) (err error) {
	return DeleteByIds[model.Filter](ids)
}

func UpdateFilter(filter model.Filter) (err error) {
	if err := validateFilterClass(filter.FilterClass); err != nil {
		return err
	}
	return Update(&filter)
}

func GetFilter(id uint) (err error, filter model.Filter) {
	filter, err = GetByID[model.Filter](id)
	if err == nil {
		err = validateFilterClass(filter.FilterClass)
	}
	return
}

func GetFilterInfoList(info request.FilterSearch) (err error, list interface{}, total int64) {
	db := global.GVA_DB.Model(&model.Filter{}).
		Where("filter_class IN ?", []string{"extension", "keyword"})
	var filters []model.Filter
	total, err = Paginate(db, info.Page, info.PageSize, &filters, "")
	return err, filters, total
}

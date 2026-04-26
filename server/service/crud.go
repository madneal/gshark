package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model/request"
	"gorm.io/gorm"
)

func Create[T any](value *T) error {
	return global.GVA_DB.Create(value).Error
}

func Delete[T any](value *T) error {
	return global.GVA_DB.Delete(value).Error
}

func DeleteByIds[T any](ids request.IdsReq) error {
	var values []T
	return global.GVA_DB.Delete(&values, "id in ?", ids.Ids).Error
}

func Update[T any](value *T) error {
	return global.GVA_DB.Save(value).Error
}

func GetByID[T any](id uint) (T, error) {
	var value T
	err := global.GVA_DB.Where("id = ?", id).First(&value).Error
	return value, err
}

func Paginate(db *gorm.DB, page, pageSize int, dest interface{}, order string) (int64, error) {
	limit, offset := pageLimitOffset(page, pageSize)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return total, err
	}
	if order != "" {
		db = db.Order(order)
	}
	return total, db.Limit(limit).Offset(offset).Find(dest).Error
}

func pageLimitOffset(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return pageSize, pageSize * (page - 1)
}

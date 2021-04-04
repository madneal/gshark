package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateSubdomain
//@description: 创建Subdomain记录
//@param: subdomain model.Subdomain
//@return: err error

func CreateSubdomain(subdomain model.Subdomain) (err error) {
	err = global.GVA_DB.Create(&subdomain).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSubdomain
//@description: 删除Subdomain记录
//@param: subdomain model.Subdomain
//@return: err error

func DeleteSubdomain(subdomain model.Subdomain) (err error) {
	err = global.GVA_DB.Delete(&subdomain).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSubdomainByIds
//@description: 批量删除Subdomain记录
//@param: ids request.IdsReq
//@return: err error

func DeleteSubdomainByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Subdomain{}, "id in ?", ids.Ids).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateSubdomain
//@description: 更新Subdomain记录
//@param: subdomain *model.Subdomain
//@return: err error

func UpdateSubdomain(subdomain model.Subdomain) (err error) {
	err = global.GVA_DB.Save(&subdomain).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSubdomain
//@description: 根据id获取Subdomain记录
//@param: id uint
//@return: err error, subdomain model.Subdomain

func GetSubdomain(id uint) (err error, subdomain model.Subdomain) {
	err = global.GVA_DB.Where("id = ?", id).First(&subdomain).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSubdomainInfoList
//@description: 分页获取Subdomain记录
//@param: info request.SubdomainSearch
//@return: err error, list interface{}, total int64

func GetSubdomainInfoList(info request.SubdomainSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.Subdomain{})
	var subdomains []model.Subdomain
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Subdomain.Subdomain != "" {
		db = db.Where("`subdomain` LIKE ?", "%"+info.Subdomain.Subdomain+"%")
	}
	if info.Domain != "" {
		db = db.Where("`domain` LIKE ?", "%"+info.Domain+"%")
	}
	if info.Status != 0 {
		db = db.Where("`status` = ?", info.Status)
	}
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&subdomains).Error
	return err, subdomains, total
}

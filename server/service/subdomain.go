package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

func CreateSubdomain(subdomain model.Subdomain) (err error) {
	err = global.GVA_DB.Create(&subdomain).Error
	return err
}

func DeleteSubdomain(subdomain model.Subdomain) (err error) {
	err = global.GVA_DB.Delete(&subdomain).Error
	return err
}

func DeleteSubdomainByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Subdomain{}, "id in ?", ids.Ids).Error
	return err
}

func UpdateSubdomain(subdomain model.Subdomain) (err error) {
	err = global.GVA_DB.Save(&subdomain).Error
	return err
}

func GetSubdomain(id uint) (err error, subdomain model.Subdomain) {
	err = global.GVA_DB.Where("id = ?", id).First(&subdomain).Error
	return
}

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

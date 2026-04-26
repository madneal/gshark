package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

func CreateSubdomain(subdomain model.Subdomain) (err error) {
	return Create(&subdomain)
}

func DeleteSubdomain(subdomain model.Subdomain) (err error) {
	return Delete(&subdomain)
}

func DeleteSubdomainByIds(ids request.IdsReq) (err error) {
	return DeleteByIds[model.Subdomain](ids)
}

func UpdateSubdomain(subdomain model.Subdomain) (err error) {
	return Update(&subdomain)
}

func GetSubdomain(id uint) (err error, subdomain model.Subdomain) {
	subdomain, err = GetByID[model.Subdomain](id)
	return
}

func GetSubdomainInfoList(info request.SubdomainSearch) (err error, list interface{}, total int64) {
	db := global.GVA_DB.Model(&model.Subdomain{})
	var subdomains []model.Subdomain
	if info.Subdomain.Subdomain != "" {
		db = db.Where("`subdomain` LIKE ?", "%"+info.Subdomain.Subdomain+"%")
	}
	if info.Domain != "" {
		db = db.Where("`domain` LIKE ?", "%"+info.Domain+"%")
	}
	if info.Status != 0 {
		db = db.Where("`status` = ?", info.Status)
	}
	total, err = Paginate(db, info.Page, info.PageSize, &subdomains, "")
	return err, subdomains, total
}

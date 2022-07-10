package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

func CreateToken(token model.Token) (err error) {
	err = global.GVA_DB.Create(&token).Error
	return err
}

func DeleteToken(token model.Token) (err error) {
	err = global.GVA_DB.Delete(&token).Error
	return err
}

func DeleteTokenByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Token{}, "id in ?", ids.Ids).Error
	return err
}

func UpdateToken(token model.Token) (err error) {
	err = global.GVA_DB.Save(&token).Error
	return err
}

func GetToken(id uint) (err error, token model.Token) {
	err = global.GVA_DB.Where("id = ?", id).First(&token).Error
	return
}

func ListTokenByType(tokenType string) (err error, tokens []model.Token) {
	err = global.GVA_DB.Where("type = ?", tokenType).Find(&tokens).Error
	return err, tokens
}

func GetTokenInfoList(info request.TokenSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.Token{})
	var tokens []model.Token
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&tokens).Error
	return err, tokens, total
}

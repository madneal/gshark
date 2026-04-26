package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

func CreateToken(token model.Token) (err error) {
	return Create(&token)
}

func DeleteToken(token model.Token) (err error) {
	return Delete(&token)
}

func DeleteTokenByIds(ids request.IdsReq) (err error) {
	return DeleteByIds[model.Token](ids)
}

func UpdateToken(token model.Token) (err error) {
	return Update(&token)
}

func GetToken(id uint) (err error, token model.Token) {
	token, err = GetByID[model.Token](id)
	return
}

func ListTokenByType(tokenType string) (err error, tokens []model.Token) {
	err = global.GVA_DB.Where("type = ?", tokenType).Order("id desc").
		Find(&tokens).Error
	return err, tokens
}

func GetTokenInfoList(info request.TokenSearch) (err error, list interface{}, total int64) {
	db := global.GVA_DB.Model(&model.Token{})
	var tokens []model.Token
	total, err = Paginate(db, info.Page, info.PageSize, &tokens, "")
	return err, tokens, total
}

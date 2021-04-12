package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateToken
//@description: 创建Token记录
//@param: token model.Token
//@return: err error

func CreateToken(token model.Token) (err error) {
	err = global.GVA_DB.Create(&token).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteToken
//@description: 删除Token记录
//@param: token model.Token
//@return: err error

func DeleteToken(token model.Token) (err error) {
	err = global.GVA_DB.Delete(&token).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteTokenByIds
//@description: 批量删除Token记录
//@param: ids request.IdsReq
//@return: err error

func DeleteTokenByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Token{}, "id in ?", ids.Ids).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateToken
//@description: 更新Token记录
//@param: token *model.Token
//@return: err error

func UpdateToken(token model.Token) (err error) {
	err = global.GVA_DB.Save(&token).Error
	return err
}

func UpdateTokenRate(token model.Token) error {
	err := global.GVA_DB.Update("rate", &token).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetToken
//@description: 根据id获取Token记录
//@param: id uint
//@return: err error, token model.Token

func GetToken(id uint) (err error, token model.Token) {
	err = global.GVA_DB.Where("id = ?", id).First(&token).Error
	return
}

func ListTokenByType(tokenType string) (err error, tokens []model.Token) {
	err = global.GVA_DB.Where(" type = ?", tokenType).Find(&tokens).Error
	return err, tokens
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetTokenInfoList
//@description: 分页获取Token记录
//@param: info request.TokenSearch
//@return: err error, list interface{}, total int64

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

package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateRule
//@description: 创建Rule记录
//@param: rule model.Rule
//@return: err error

func CreateRule(rule model.Rule) (err error) {
	err = global.GVA_DB.Create(&rule).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteRule
//@description: 删除Rule记录
//@param: rule model.Rule
//@return: err error

func DeleteRule(rule model.Rule) (err error) {
	err = global.GVA_DB.Delete(&rule).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteRuleByIds
//@description: 批量删除Rule记录
//@param: ids request.IdsReq
//@return: err error

func DeleteRuleByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Rule{}, "id in ?", ids.Ids).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateRule
//@description: 更新Rule记录
//@param: rule *model.Rule
//@return: err error

func UpdateRule(rule model.Rule) (err error) {
	err = global.GVA_DB.Save(&rule).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetRule
//@description: 根据id获取Rule记录
//@param: id uint
//@return: err error, rule model.Rule

func GetRule(id uint) (err error, rule model.Rule) {
	err = global.GVA_DB.Where("id = ?", id).First(&rule).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetRuleInfoList
//@description: 分页获取Rule记录
//@param: info request.RuleSearch
//@return: err error, list interface{}, total int64

func GetRuleInfoList(info request.RuleSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.Rule{})
	var rules []model.Rule
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Type != "" {
		db = db.Where("`type` = ?", info.Type)
	}
	if info.Content != "" {
		db = db.Where("`content` LIKE ?", "%"+info.Content+"%")
	}
	if info.Name != "" {
		db = db.Where("`name` = ?", info.Name)
	}
	if info.Status != 0 {
		db = db.Where("`status` = ?", info.Status)
	}
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&rules).Error
	return err, rules, total
}

func GetValidRulesByType(typeStr string) (err error, list []model.Rule) {
	err = global.GVA_DB.Where("type = ? and status = 1", typeStr).Find(&list).Error
	return err, list
}

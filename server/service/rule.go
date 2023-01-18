package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"gorm.io/gorm/clause"
)

func CreateRule(rule model.Rule) (err error) {
	err = global.GVA_DB.Clauses(clause.Insert{
		Modifier: "IGNORE",
	}).Create(&rule).Error
	return err
}

func DeleteRule(rule model.Rule) (err error) {
	err = global.GVA_DB.Unscoped().Delete(&rule).Error
	return err
}

func DeleteRuleByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Debug().Delete(&[]model.Rule{}, ids.Ids).Error
	return err
}

func UpdateRule(rule model.Rule) (err error) {
	err = global.GVA_DB.Save(&rule).Error
	return err
}

func GetRule(id uint) (err error, rule model.Rule) {
	err = global.GVA_DB.Where("id = ?", id).First(&rule).Error
	return
}

func GetRuleInfoList(info request.RuleSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&model.Rule{})
	var rules []model.Rule
	if len(info.RuleType) > 0 {
		db = db.Where("find_in_set(?, rule_type)", info.RuleType)
	}
	if info.Content != "" {
		db = db.Where("`content` LIKE ?", "%"+info.Content+"%")
	}
	if info.Name != "" {
		db = db.Where("`name` = ?", info.Name)
	}
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&rules).Error
	return err, rules, total
}

func GetValidRulesByType(typeStr string) (err error, list []model.Rule) {
	err = global.GVA_DB.Where("find_in_set(?, rule_type) and status = 1", typeStr).Find(&list).Error
	return err, list
}

func SwitchRuleStatus(id uint, status int) error {
	err := global.GVA_DB.Table("rule").Where("id = ?", id).UpdateColumn("status", status).Error
	return err
}

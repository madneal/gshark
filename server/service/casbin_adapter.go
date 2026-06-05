package service

import (
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gsharkModel "github.com/madneal/gshark/model"
	"gorm.io/gorm"
)

type GormCasbinAdapter struct {
	db *gorm.DB
}

func NewGormCasbinAdapter(db *gorm.DB) *GormCasbinAdapter {
	return &GormCasbinAdapter{db: db}
}

func (a *GormCasbinAdapter) LoadPolicy(model model.Model) error {
	var lines []gsharkModel.CasbinModel
	if err := a.db.Find(&lines).Error; err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}
	return nil
}

func (a *GormCasbinAdapter) SavePolicy(model model.Model) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&gsharkModel.CasbinModel{}).Error; err != nil {
			return err
		}

		for ptype, ast := range model["p"] {
			for _, rule := range ast.Policy {
				if err := tx.Create(savePolicyLine(ptype, rule)).Error; err != nil {
					return err
				}
			}
		}

		for ptype, ast := range model["g"] {
			for _, rule := range ast.Policy {
				if err := tx.Create(savePolicyLine(ptype, rule)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (a *GormCasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return a.db.Create(savePolicyLine(ptype, rule)).Error
}

func (a *GormCasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return deletePolicyLine(a.db, savePolicyLine(ptype, rule))
}

func (a *GormCasbinAdapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			if err := tx.Create(savePolicyLine(ptype, rule)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *GormCasbinAdapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			if err := deletePolicyLine(tx, savePolicyLine(ptype, rule)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *GormCasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := &gsharkModel.CasbinModel{Ptype: ptype}

	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.AuthorityId = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.Path = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.Method = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}

	return deletePolicyLine(a.db, line)
}

func loadPolicyLine(line gsharkModel.CasbinModel, model model.Model) {
	values := []string{line.Ptype, line.AuthorityId, line.Path, line.Method, line.V3, line.V4, line.V5}
	for len(values) > 1 && values[len(values)-1] == "" {
		values = values[:len(values)-1]
	}
	persist.LoadPolicyLine(strings.Join(values, ", "), model)
}

func savePolicyLine(ptype string, rule []string) *gsharkModel.CasbinModel {
	line := &gsharkModel.CasbinModel{Ptype: ptype}
	if len(rule) > 0 {
		line.AuthorityId = rule[0]
	}
	if len(rule) > 1 {
		line.Path = rule[1]
	}
	if len(rule) > 2 {
		line.Method = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}
	return line
}

func deletePolicyLine(db *gorm.DB, line *gsharkModel.CasbinModel) error {
	query := db.Where("p_type = ?", line.Ptype)
	if line.AuthorityId != "" {
		query = query.Where("v0 = ?", line.AuthorityId)
	}
	if line.Path != "" {
		query = query.Where("v1 = ?", line.Path)
	}
	if line.Method != "" {
		query = query.Where("v2 = ?", line.Method)
	}
	if line.V3 != "" {
		query = query.Where("v3 = ?", line.V3)
	}
	if line.V4 != "" {
		query = query.Where("v4 = ?", line.V4)
	}
	if line.V5 != "" {
		query = query.Where("v5 = ?", line.V5)
	}
	return query.Delete(&gsharkModel.CasbinModel{}).Error
}

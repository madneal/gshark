package models

import (
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/vars"

	"bufio"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Rule struct {
	Id          int64
	Position    string
	Type        string
	Pattern     string
	Caption     string
	Description string `xorm:"text"`
	Status      int    `xorm:"int default 0 notnull"`
}

func NewRule(ruleType, pat, caption, pos, desc string, status int) *Rule {
	return &Rule{Type: ruleType, Pattern: pat, Caption: caption,
		Position: pos, Description: desc, Status: status}
}

func (r *Rule) Insert() (err error) {
	_, err = Engine.Insert(r)
	return err
}

func GetRules() ([]Rule, error) {
	rules := make([]Rule, 0)
	err := Engine.Table("rule").Where("status=1").Find(&rules)
	return rules, err
}

func GetRulesPage(page int) ([]Rule, int, error) {
	rules := make([]Rule, 0)
	totalPages, err := Engine.Table("rule").Count()
	var pages int

	if int(totalPages)%vars.PAGE_SIZE == 0 {
		pages = int(totalPages) / vars.PAGE_SIZE
	} else {
		pages = int(totalPages)/vars.PAGE_SIZE + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}
	err = Engine.Table("rule").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Desc("status").Find(&rules)
	return rules, pages, err
}

func GetRuleById(id int64) (*Rule, bool, error) {
	rule := new(Rule)
	has, err := Engine.ID(id).Get(rule)
	return rule, has, err
}

func EditRuleById(id int64, position, ruleType, pat, caption, desc string, status int) error {
	rule := new(Rule)
	_, has, err := GetRuleById(id)
	if err == nil && has {
		rule.Position = position
		rule.Type = ruleType
		rule.Pattern = pat
		rule.Caption = caption
		rule.Description = desc
		rule.Status = status
		_, err = Engine.ID(id).Update(rule)
	}
	return err
}

func DeleteRulesById(id int64) (err error) {
	rule := new(Rule)
	_, err = Engine.Id(id).Delete(rule)
	return err
}

func EnableRulesById(id int64) (err error) {
	rules := new(Rule)
	has, err := Engine.Id(id).Get(rules)
	if err == nil && has {
		rules.Status = 1
		_, err = Engine.Id(id).Cols("status").Update(rules)
	}
	return err
}

func DisableRulesById(id int64) (err error) {
	rules := new(Rule)
	has, err := Engine.Id(id).Get(rules)
	if err == nil && has {
		rules.Status = 0
		_, err = Engine.Id(id).Cols("status").Update(rules)
	}
	return err
}

func LoadRuleFromFile(filename string) ([]Rule, error) {
	ruleFile, err := os.Open(filename)
	rules := make([]Rule, 0)
	var content []byte
	if err == nil {
		r := bufio.NewReader(ruleFile)
		content, err = ioutil.ReadAll(r)
		if err == nil {
			err = json.Unmarshal(content, &rules)
		}
	}
	return rules, err
}

func LoadBlackListRuleFromFile(filename string) ([]FilterRule, error) {
	file, err := os.Open(filename)
	rules := make([]FilterRule, 0)
	var content []byte
	if err == nil {
		r := bufio.NewReader(file)
		content, err = ioutil.ReadAll(r)
		if err == nil {
			err = yaml.Unmarshal(content, &rules)
		}
	}
	return rules, err
}

func InsertBlacklistRules(filename string) error {
	rules, err := LoadBlackListRuleFromFile(filename)
	if err == nil {
		for _, rule := range rules {
			rule.Insert()
		}
	}
	return err
}

func InsertRules(filename string) error {
	rules, err := LoadRuleFromFile(filename)
	if err == nil {
		for _, rule := range rules {
			rule.Insert()
		}
	}
	return err
}

func GetValidRulesByType(Type string) ([]Rule, error) {
	rules := make([]Rule, 0)
	err := Engine.Table("rule").Where("status=1 and type=?", Type).Find(&rules)
	return rules, err
}

func GetValidRules() ([]Rule, error) {
	rules := make([]Rule, 0)
	err := Engine.Table("rule").Where("status=1").Find(&rules)
	return rules, err
}

func GetRulesCount() int {
	count, err := Engine.Table("rule").Where("type = subdomain").Count()
	if err != nil {
		logger.Log.Error(err)
	}
	return int(count)
}

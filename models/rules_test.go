package models_test

import (
	"github.com/neal1991/gshark/models"

	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadRuleFromFile(t *testing.T) {
	filename := "/data/code/golang/src/github.com/neal1991/gshark/conf/gitrob.json"
	t.Log(models.LoadRuleFromFile(filename))
}

func TestInsertRules(t *testing.T) {
	filename := "/data/code/golang/src/github.com/neal1991/gshark/conf/gitrob.json"
	rules, err := models.GetRules()
	t.Log(rules, err)
	if err == nil && len(rules) == 0 {
		t.Logf("Init rules, err: %v", models.InsertRules(filename))
	}
}

func TestGetRulesPage(t *testing.T) {
	rules, _, err := models.GetRulesPage(0)
	fmt.Println("The length of rules shoulbe larger than 0")
	assert.True(t, len(rules) > 0)
	fmt.Println("There should be no err")
	assert.True(t, err == nil)
}

func TestGetValidRulesByType(t *testing.T) {
	rules, err := models.GetValidRulesByType("github")
	assert.True(t, err == nil, "There should be no error for GetValidRulesByType of github")
	for rule := range rules {
		fmt.Println(rule)
	}
}

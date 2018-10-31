package models_test

import (
	"github.com/neal1991/gshark/models"

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

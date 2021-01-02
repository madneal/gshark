package models_test

import (
	"github.com/madneal/gshark/models"

	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadRuleFromFile(t *testing.T) {
	filename := "/data/code/golang/src/github.com/madneal/gshark/conf/gitrob.json"
	t.Log(models.LoadRuleFromFile(filename))
}

func TestInsertRules(t *testing.T) {
	filename := "/data/code/golang/src/github.com/madneal/gshark/conf/gitrob.json"
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
	for _, rule := range rules {
		fmt.Println(rule.Pattern)
	}
}

func TestConvertTextToRules(t *testing.T) {
	text := "github|keyword \n gitlab|keyword"
	rules, _ := models.ConvertTextToRules(&text)
	assert.Equal(t, 2, len(*rules), "the length should be the same")
	assert.Equal(t, "github", (*rules)[0].Type, "the type of the first element should be github")
	assert.Equal(t, "keyword", (*rules)[1].Pattern, "the pattern of the second element should be keyword")
	text1 := "github"
	_, err := models.ConvertTextToRules(&text1)
	assert.True(t, err != nil, "there should be least 2 argument")
	text2 := "github|keyword|name|desc\ngitlab|key|name|desc"
	rules2, _ := models.ConvertTextToRules(&text2)
	assert.Equal(t, "desc", (*rules2)[0].Description, "the description should be the same")
	assert.Equal(t, "desc", (*rules2)[1].Description, "the description should be the same")
}

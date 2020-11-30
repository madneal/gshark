package appsearch

import (
	"fmt"
	"github.com/madneal/gshark/models"
	"testing"
)

func TestSearchForApp(t *testing.T) {
	rule := new(models.Rule)
	rule.Pattern = "abc "
	rule.Caption = "HUAWEI"
	SearchForApp(*rule)
}

func TestSaveResults(t *testing.T) {
	rule := new(models.Rule)
	rule.Pattern = "abc "
	rule.Caption = "HUAWEI"
	SaveResults(SearchForApp(*rule))
}

func TestGenerateSearchCodeTask(t *testing.T) {
	rulesMap, err := GenerateSearchCodeTask()
	for _, rules := range rulesMap {
		for _, rule := range rules {
			fmt.Println(rule)
		}
	}
	fmt.Println(err)
}

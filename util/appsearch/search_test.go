package appsearch

import (
	"fmt"
	"github.com/neal1991/gshark/models"
	"testing"
)

func TestSearchForApp(t *testing.T) {
	rule := new(models.Rule)
	rule.Pattern = "浦发 "
	rule.Caption = "HUAWEI"
	SearchForApp(*rule)
}

func TestSaveResults(t *testing.T) {
	rule := new(models.Rule)
	rule.Pattern = "浦发 "
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

func TestScheduleTasks(t *testing.T) {
	ScheduleTasks(100)
}

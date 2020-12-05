package githubsearch

import "testing"
import "github.com/madneal/gshark/models"

func TestFiltetString(t *testing.T) {
	codeResult := new(models.CodeResult)
	id := 1
	_, _ = models.Engine.Table("code_result").Where("id=?", id).Get(codeResult)
	if !PassFilters(codeResult, "abc") {
		t.Log("pass the PassFilters function")
	} else {
		t.Error("failed to pass the PassFilters function")
	}
}

func TestScheduleTasks(t *testing.T) {
	ScheduleTasks(100)
}

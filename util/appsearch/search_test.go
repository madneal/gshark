package appsearch

import (
	"testing"
	"github.com/neal1991/gshark/models"
)

func TestSearchForApp(t *testing.T) {
	rule := new(models.Rule)
	rule.Pattern = "浦发 "
	rule.Caption = "HUAWEI"
	SearchForApp(*rule)
}

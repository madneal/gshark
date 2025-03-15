package githubsearch

import (
	"fmt"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/model"
	"testing"
)

func TestSearch(t *testing.T) {
	global.GVA_VP = initialize.Viper("../../config.yaml") // 初始化Viper
	global.GVA_LOG = initialize.Zap()
	global.GVA_DB = initialize.Gorm()
	if global.GVA_DB == nil {
		fmt.Println("init db failed")
		return
	}
	rules := make([]model.Rule, 0)
	rules = append(rules, model.Rule{
		Content: "mihoyo",
	})
	Search(rules)
}

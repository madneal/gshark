package githubsearch

import (
	"fmt"
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	global.GVA_VP = core.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = core.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()                                           // gorm连接数据库
	q, _ := BuildQuery("spdb")
	fmt.Println(q)
}

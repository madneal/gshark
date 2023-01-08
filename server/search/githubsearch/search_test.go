package githubsearch

import (
	"context"
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
	q, _ := BuildQuery("/partner_key.*[\"|'\\s]([a-z0-9]{64})[\"|'\\s]/")
	fmt.Println(q)
}

func TestClient_SearchCode(t *testing.T) {
	global.GVA_VP = core.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = core.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()
	client, _ := GetGithubClient()
	results, _ := client.SearchCode("/partner_key.*[\"|'\\s]([a-z0-9]{64})[\"|'\\s]/")
	fmt.Println(results)
}

func TestClient_GetCommiter(t *testing.T) {
	global.GVA_VP = core.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = core.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()
	client, _ := GetGithubClient()
	email := client.GetCommiter(context.Background(), "madneal", "gshark")
	fmt.Println(email)
}

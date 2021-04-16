package gitlabsearch

import (
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetClient(t *testing.T) {
	initialDataBase()
	client := GetClient()
	assert.Equal(t, true, client != nil, "the client is not nil")

}

func initialDataBase() {
	global.GVA_VP = core.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = core.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()                                           // gorm连接数据库
}

func TestGetProjects(t *testing.T) {
	initialDataBase()
	client := GetClient()
	GetProjects(client)
}

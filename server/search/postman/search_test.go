package postman

import (
	"fmt"
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"testing"
)

func TestRunTask(t *testing.T) {
	InitialDataBase()
	RunTask()
}

func TestClient_SearchAPI(t *testing.T) {
	//InitialDataBase()
	res, err := SearchAPI("mihoyo", "collection")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println((*res)[0])
}

func InitialDataBase() {
	global.GVA_VP = core.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = core.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()                                           // gorm连接数据库
}

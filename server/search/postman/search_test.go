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

}

func TestClient_SearchAPI(t *testing.T) {
	InitialDataBase()
	client := GetPostmanClient()
	res, _ := client.SearchAPI("mihoyo")
	fmt.Println(res)
}

func InitialDataBase() {
	global.GVA_VP = core.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = core.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()                                           // gorm连接数据库
}

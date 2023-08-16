package main

import (
	"github.com/gookit/color"
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/search"
)

func main() {
	global.GVA_VP = core.Viper()      // 初始化Viper
	global.GVA_LOG = core.Zap()       // 初始化zap日志库
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	if global.GVA_DB != nil {
		initialize.MysqlTables(global.GVA_DB)
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	} else {
		color.Danger.Println("数据库连接失败，请确定在 config.yaml 配置正确数据库信息")
	}
	go search.ScanTask()
	core.RunServer()
}

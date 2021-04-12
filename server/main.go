package main

import (
	"github.com/madneal/gshark/cmd"
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
)

// @title Swagger Example API
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /
func main() {
	global.GVA_VP = core.Viper()      // 初始化Viper
	global.GVA_LOG = core.Zap()       // 初始化zap日志库
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	if global.GVA_DB != nil {
		initialize.MysqlTables(global.GVA_DB) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	app := cli.NewApp()
	app.Name = "GShark"
	app.Usage = "Scan for sensitive information easily and effectively."
	app.Commands = []*cli.Command{&cmd.Web, &cmd.Scan}
	app.Flags = append(app.Flags, cmd.Web.Flags...)
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	if err != nil {
		global.GVA_LOG.Error("app start error", zap.Any("err", err))
	}
}

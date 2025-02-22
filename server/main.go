package main

import (
	"github.com/gookit/color"
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/search"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	global.GVA_VP = core.Viper()
	global.GVA_LOG = core.Zap()
	global.GVA_DB = initialize.Gorm()
}

func main() {
	rootCmd := &cobra.Command{
		Use:  "gshark",
		Long: "GShark is a tool to monitor the sensitive information disclosure for multi platforms",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	var configFile string
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml",
		"config file")
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the GShark server",
		Long:  "Start the GShark web server, supports for the management platform",
		Run: func(cmd *cobra.Command, args []string) {
			core.RunServer()
		},
	}
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Start the scan task",
		Long:  "Support the scan task for multi platforms, including: GitHub, GitLab, Postman, searchcode",
		Run: func(cmd *cobra.Command, args []string) {
			search.ScanTask()
		},
	}
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(scanCmd)
	if err := rootCmd.Execute(); err != nil {
		color.Println(err)
		os.Exit(1)
	}
	//if global.GVA_DB != nil {
	//	service.InitDB(request.InitDB{
	//		Host: global.GVA_CONFIG.Mysql.Path,
	//	})
	//	db, _ := global.GVA_DB.DB()
	//	defer db.Close()
	//} else {
	//	color.Danger.Println("数据库连接失败，请确定在 config.yaml 配置正确数据库信息")
	//}

}

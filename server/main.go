package main

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/search"
	"github.com/madneal/gshark/service"
	"github.com/spf13/cobra"
)

func init() {
	global.GVA_VP = initialize.Viper()
	global.GVA_LOG = initialize.Zap()
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

	var (
		initHost          string
		initPort          string
		initUser          string
		initPassword      string
		initDBName        string
		initAdminUser     string
		initAdminPassword string
	)
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize database and seed admin user",
		Long:  "Create database, run migrations, seed data. Admin credentials via --admin-user / --admin-password (default gshark/gshark).",
		Run: func(cmd *cobra.Command, args []string) {
			if global.GVA_DB != nil {
				// Already connected — treat as already initialized unless empty.
				var n int64
				if err := global.GVA_DB.Table("sys_users").Count(&n).Error; err == nil && n > 0 {
					color.Warn.Println("database already initialized (sys_users not empty); skip")
					return
				}
			}
			conf := request.InitDB{
				Host:          initHost,
				Port:          initPort,
				UserName:      initUser,
				Password:      initPassword,
				DBName:        initDBName,
				AdminUserName: initAdminUser,
				AdminPassword: initAdminPassword,
			}
			if err := service.InitDB(conf); err != nil {
				color.Error.Printf("init failed: %v\n", err)
				os.Exit(1)
			}
			adminUser := initAdminUser
			if adminUser == "" {
				adminUser = "gshark"
			}
			color.Success.Println("database initialized successfully")
			fmt.Printf("admin login: %s / (your --admin-password)\n", adminUser)
		},
	}
	initCmd.Flags().StringVar(&initHost, "host", "127.0.0.1", "MySQL host")
	initCmd.Flags().StringVar(&initPort, "port", "3306", "MySQL port")
	initCmd.Flags().StringVar(&initUser, "user", "root", "MySQL username")
	initCmd.Flags().StringVar(&initPassword, "password", "", "MySQL password")
	initCmd.Flags().StringVar(&initDBName, "db", "gshark", "MySQL database name")
	initCmd.Flags().StringVar(&initAdminUser, "admin-user", "gshark", "Admin login username")
	initCmd.Flags().StringVar(&initAdminPassword, "admin-password", "gshark", "Admin login password (plain)")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		color.Println(err)
		os.Exit(1)
	}
}

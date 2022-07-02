package core

import (
	"flag"
	"fmt"
	"github.com/madneal/gshark/config"
	"github.com/madneal/gshark/global"
	_ "github.com/madneal/gshark/packfile"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper(path ...string) *viper.Viper {
	var configFilename string
	if len(path) == 0 {
		flag.StringVar(&configFilename, "c", "", "choose configFilename file.")
		flag.Parse()
		// 优先级: 命令行 > 环境变量 > 默认值
		if configFilename == "" {
			if configEnv := os.Getenv(config.ConfigEnv); configEnv == "" {
				configFilename = config.ConfigFile
				fmt.Printf("您正在使用config的默认值,config的路径为%v\n", config.ConfigFile)
			} else {
				configFilename = configEnv
				fmt.Printf("您正在使用GVA_CONFIG环境变量,config的路径为%v\n", configFilename)
			}
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", configFilename)
		}
	} else {
		configFilename = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", configFilename)
	}

	v := viper.New()
	v.SetConfigFile(configFilename)
	err := v.ReadInConfig()
	fmt.Println(os.Getwd())
	if err != nil {
		panic(fmt.Errorf("Fatal error configFilename file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("configFilename file changed:", e.Name)
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}

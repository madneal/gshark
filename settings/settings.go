package settings

import (
	"fmt"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/vars"

	"gopkg.in/ini.v1"
	"os"
	"strings"
)

var (
	Cfg *ini.File
)

func init() {
	var err error
	dirName, _ := os.Getwd()
	endIndex := strings.Index(dirName, "gshark")
	var source string

	if endIndex > 0 {
		source = dirName[:endIndex] + "gshark/conf/app.ini"
	} else {
		source = "conf/app.ini"
	}
	Cfg, err = ini.Load(source)

	if err != nil {
		fmt.Println("Please check the config file app.ini")
		logger.Log.Panicln(err)
	}

	vars.HTTP_HOST = Cfg.Section("").Key("HTTP_HOST").MustString("127.0.0.1")
	vars.HTTP_PORT = Cfg.Section("").Key("HTTP_PORT").MustInt(8000)

	vars.REPO_PATH = Cfg.Section("").Key("REPO_PATH").MustString("repos")
	vars.MAX_INDEXERS = Cfg.Section("").Key("MAX_INDEXERS").MustInt(vars.DefaultMaxConcurrentIndexers)
	vars.MAX_Concurrency_REPOS = Cfg.Section("").Key("MAX_Concurrency_REPOS").MustInt(100)
	vars.SCKEY = Cfg.Section("").Key("SCKEY").MustString("")
	vars.GITLAB_TOKEN = Cfg.Section("").Key("gitlab_token").MustString("")
}

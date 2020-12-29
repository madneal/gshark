package vars

import (
	"fmt"
	"github.com/madneal/gshark/logger"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

const (
	DefaultMaxConcurrentIndexers = 2
	PageStep                     = 5
	SearchNum                    = 25
	GITLAB                       = "gitlab"
)

var (
	REPO_PATH             string
	MAX_INDEXERS          int
	HTTP_HOST             string
	HTTP_PORT             int
	MAX_Concurrency_REPOS int
	DEBUG_MODE            bool
	PAGE_SIZE             = 10
	SCKEY                 string
	GOBUSTER              string
	SUBDOMAIN_WORDLIST    string
	ENABLE_SUBDOMAIN      bool
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
	Cfg, err = ini.LoadSources(ini.LoadOptions{SpaceBeforeInlineComment: true}, source)

	if err != nil {
		fmt.Println("Please check the config file app.ini")
		logger.Log.Panicln(err)
	}

	HTTP_HOST = Cfg.Section("").Key("HTTP_HOST").MustString("127.0.0.1")
	HTTP_PORT = Cfg.Section("").Key("HTTP_PORT").MustInt(8000)
	REPO_PATH = Cfg.Section("").Key("REPO_PATH").MustString("repos")
	MAX_INDEXERS = Cfg.Section("").Key("MAX_INDEXERS").MustInt(DefaultMaxConcurrentIndexers)
	MAX_Concurrency_REPOS = Cfg.Section("").Key("MAX_Concurrency_REPOS").MustInt(100)
	SCKEY = Cfg.Section("").Key("SCKEY").MustString("")
	GOBUSTER = Cfg.Section("").Key("gobuster_path").MustString("")
	SUBDOMAIN_WORDLIST = Cfg.Section("").Key("subdomain_wordlist_file").MustString("")
}

package models

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/vars"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/core"
	"xorm.io/xorm"

	"fmt"
	"path/filepath"
)

var (
	DATA_TYPE string
	DATA_HOST string
	DATA_PORT int
	DATA_NAME string
	USERNAME  string
	PASSWORD  string
	SSL_MODE  string
	DATA_PATH string

	Engine *xorm.Engine
)

func init() {
	cfg := vars.Cfg
	sec := cfg.Section("database")

	DATA_TYPE = sec.Key("DB_TYPE").MustString("sqlite")
	DATA_HOST = sec.Key("HOST").MustString("127.0.0.1")
	DATA_PORT = sec.Key("PORT").MustInt(3306)
	USERNAME = sec.Key("USER").MustString("username")
	PASSWORD = sec.Key("PASSWD").MustString("password")
	SSL_MODE = sec.Key("SSL_MODE").MustString("disable")
	DATA_PATH = sec.Key("PATH").MustString("conf")
	DATA_NAME = sec.Key("NAME").MustString("xsec")

	err := NewDbEngine()
	if err != nil {
		logger.Log.Error(err)
	} else {
		err := Engine.Sync2(new(Rule))
		err = Engine.Sync2(new(InputInfo))
		err = Engine.Sync2(new(Admin))
		err = Engine.Sync2(new(Repo))
		err = Engine.Sync2(new(GitToken))
		err = Engine.Sync2(new(CodeResult))
		err = Engine.Sync2(new(FilterRule))
		err = Engine.Sync2(new(CodeResultDetail))
		err = Engine.Sync(new(Subdomain))
		if err != nil {
			logger.Log.Error(err)
		}
		InitRules()
		InitAdmin()
	}
}

// init a database instance
func NewDbEngine() (err error) {
	switch DATA_TYPE {
	case "sqlite":
		//cur, _ := filepath.Abs(".")
		dataSourceName := fmt.Sprintf("%v/%v.db", DATA_PATH, DATA_NAME)
		logger.Log.Infof("sqlite db: %v", dataSourceName)
		Engine, err = xorm.NewEngine("sqlite3", dataSourceName)
		Engine.Logger().SetLevel(core.LOG_OFF)
		err = Engine.Ping()

	case "mysql":
		dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
			USERNAME, PASSWORD, DATA_HOST, DATA_PORT, DATA_NAME)

		Engine, err = xorm.NewEngine("mysql", dataSourceName)
		if err != nil {
			fmt.Println(err)
		}
		Engine.Logger().SetLevel(core.LOG_OFF)
		err = Engine.Ping()
	case "postgres":
		dbSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", USERNAME, PASSWORD, DATA_HOST,
			DATA_PORT, DATA_NAME, SSL_MODE)
		Engine, err = xorm.NewEngine("postgres", dbSourceName)
		Engine.Logger().SetLevel(core.LOG_OFF)
		err = Engine.Ping()

	default:
		cur, _ := filepath.Abs(".")
		dataSourceName := fmt.Sprintf("%v/%v/%v.db", cur, DATA_PATH, DATA_NAME)
		logger.Log.Infof("sqlite db: %v", dataSourceName)
		Engine, err = xorm.NewEngine("sqlite3", dataSourceName)
		Engine.Logger().SetLevel(core.LOG_OFF)
		err = Engine.Ping()
	}

	return err
}

func InitRules() {
	cur, _ := filepath.Abs(".")
	ruleFile := fmt.Sprintf("%v/conf/gitrob.json", cur)
	rules, err := GetRules()
	blacklistFile := fmt.Sprintf("%v/conf/blacklist.yaml", cur)
	blacklistRules, err1 := GetFilterRules()
	if err == nil && len(rules) == 0 {
		logger.Log.Infof("Init rules, err: %v", InsertRules(ruleFile))
	} else if err != nil {
		logger.Log.Println(err)
	}

	if err1 == nil && len(blacklistRules) == 0 {
		logger.Log.Infof("Init filter rules, err: %v", InsertBlacklistRules(blacklistFile))
	} else if err1 != nil {
		logger.Log.Println(err1)
	}
}

func InitAdmin() {
	admins, err := ListAdmins()
	if err == nil && len(admins) == 0 {
		username := "gshark"
		password := "gshark"
		role := "admin"
		admin := NewAdmin(username, password, role)
		_, err := admin.Insert()
		if err != nil {
			logger.Log.Error(err)
		}
	}
}

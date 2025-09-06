package service

import (
	"database/sql"
	"fmt"

	"github.com/madneal/gshark/config"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/source"
	"github.com/madneal/gshark/utils"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func writeConfig(viper *viper.Viper, mysql config.Mysql) error {
	global.GVA_CONFIG.Mysql = mysql
	cs := utils.StructToMap(global.GVA_CONFIG)
	for k, v := range cs {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
}

func executeSql(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

func initDB(InitDBFunctions ...model.InitDBFunc) (err error) {
	for _, v := range InitDBFunctions {
		err = v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func InitDB(conf request.InitDB) error {
	BaseMysql := config.Mysql{
		Path:     "",
		Dbname:   "",
		Username: "",
		Password: "",
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	}

	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}

	if conf.Port == "" {
		conf.Port = "3306"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", conf.UserName, conf.Password, conf.Host, conf.Port)
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", conf.DBName)
	if err := executeSql(dsn, "mysql", createSql); err != nil {
		return err
	}

	MysqlConfig := config.Mysql{
		Path:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Dbname:   conf.DBName,
		Username: conf.UserName,
		Password: conf.Password,
		Config:   "charset=utf8mb4&parseTime=True&loc=Local",
	}

	if err := writeConfig(global.GVA_VP, MysqlConfig); err != nil {
		return err
	}
	m := global.GVA_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}

	linkDns := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       linkDns, // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		global.GVA_DB = db
	}

	err := global.GVA_DB.AutoMigrate(
		model.SysUser{},
		model.SysAuthority{},
		model.SysApi{},
		model.SysBaseMenu{},
		model.SysBaseMenuParameter{},
		model.JwtBlacklist{},
		model.SysDictionary{},
		model.SysDictionaryDetail{},
		model.SysOperationRecord{},
		model.Rule{},
		model.Token{},
		model.SearchResult{},
		model.Subdomain{},
		model.Filter{},
		model.Repo{},
		model.Task{},
	)
	if err != nil {
		return err
	}
	err = initDB(
		source.Admin,
		source.Api,
		source.AuthorityMenu,
		source.Authority,
		source.AuthoritiesMenus,
		source.Casbin,
		source.DataAuthorities,
		source.Dictionary,
		source.DictionaryDetail,
		source.BaseMenu)
	if err != nil {
		_ = writeConfig(global.GVA_VP, BaseMysql)
		return err
	}
	return nil
}

func RemoveDB(conf request.InitDB) error {
	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}

	if conf.Port == "" {
		conf.Port = "3306"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", conf.UserName, conf.Password, conf.Host, conf.Port)
	createSql := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", conf.DBName)
	fmt.Println(createSql)
	err := executeSql(dsn, "mysql", createSql)
	return err
}

package model

import (
	"github.com/madneal/gshark/global"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestSearchResult_CheckUrlExists(t *testing.T) {
	global.GVA_DB = InitialDb()
	global.GVA_LOG = zap.New(nil)
	searchResult := &SearchResult{
		Url: "adfasdfasdf",
	}
	result := searchResult.CheckPathExists()
	assert.Equal(t, false, result, "the url should not exists")
	searchResult1 := &SearchResult{
		Url: "https://baidu.com",
	}
	result1 := searchResult1.CheckPathExists()
	assert.Equal(t, true, result1, "the result should exists")
}

func InitialDb() *gorm.DB {
	dsn := "gshark:gshark@tcp(127.0.0.1:3306)/gshark?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), nil); err != nil {
		//global.GVA_LOG.Error("MySQL启动异常", zap.Any("err", err))
		//os.Exit(0)
		//return nil
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxOpenConns(0)
		return db
	}
}

var level zapcore.Level

func Zap() (logger *zap.Logger) {
	switch global.GVA_CONFIG.Zap.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if global.GVA_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

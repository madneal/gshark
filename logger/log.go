package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	Log *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.Formatter = new(prefixed.TextFormatter)
	logger.Level = logrus.DebugLevel
	Log = logger.WithFields(logrus.Fields{"prefix": "GShark"})
}

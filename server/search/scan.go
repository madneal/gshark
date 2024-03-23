package search

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/search/codesearch"
	"github.com/madneal/gshark/search/githubsearch"
	"github.com/madneal/gshark/search/gitlabsearch"
	"github.com/madneal/gshark/search/gobuster"
	"github.com/madneal/gshark/search/postman"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"time"
)

func ScanTask() {
	for {
		if global.GVA_DB == nil {
			return
		}
		var Interval time.Duration = 900
		if enable, err := service.CheckTaskStatus("gitlab"); enable {
			gitlabsearch.RunTask(Interval)
		} else if err != nil {
			global.GVA_LOG.Error("CheckTaskStatus for gitlab err", zap.Error(err))
			break
		}
		if enable, err := service.CheckTaskStatus("codesearch"); enable {
			codesearch.RunTask(Interval)
		} else if err != nil {
			global.GVA_LOG.Error("CheckTaskStatus for codesearch err", zap.Error(err))
			break
		}
		if enable, err := service.CheckTaskStatus("github"); enable {
			githubsearch.RunTask(Interval)
		} else if err != nil {
			global.GVA_LOG.Error("CheckTaskStatus for github err", zap.Error(err))
			break
		}
		if enable, err := service.CheckTaskStatus("gobuster"); enable {
			gobuster.RunTask(Interval)
		} else if err != nil {
			global.GVA_LOG.Error("CheckTaskStatus for gobuster err", zap.Error(err))
			break
		}
		if enable, err := service.CheckTaskStatus("postman"); enable {
			postman.RunTask()
		} else if err != nil {
			global.GVA_LOG.Error("CheckTaskStatus for postman err", zap.Error(err))
			break
		}
	}

}

package util

import (
	"github.com/neal1991/gshark/util/codesearch"
	"github.com/neal1991/gshark/util/githubsearch"

	"github.com/urfave/cli"

	"github.com/neal1991/gshark/logger"
	"github.com/neal1991/gshark/util/appsearch"
	"strings"
	"time"
)

func Scan(ctx *cli.Context) {
	scanMode := "all"
	// seconds
	var Interval time.Duration = 900

	if ctx.IsSet("mode") {
		scanMode = strings.ToLower(ctx.String("mode"))
	}

	if ctx.IsSet("time") {
		Interval = time.Duration(ctx.Int("time"))
	}

	switch scanMode {
	case "github":
		// use go keyword or not
		logger.Log.Println("scan github code")
		githubsearch.ScheduleTasks(Interval)
	case "app":
		logger.Log.Println("scan app results")
		appsearch.ScheduleTasks(Interval)
	case "searchcode":
		logger.Log.Println("scan searchcode results")
		codesearch.ScheduleTasks(Interval)
	case "all":
		logger.Log.Println("all scan mode")
		codesearch.ScheduleTasks(Interval)
		githubsearch.ScheduleTasks(Interval)
		appsearch.ScheduleTasks(Interval)
	default:
		logger.Log.Println("default scan mode")
		go githubsearch.ScheduleTasks(Interval)
		appsearch.ScheduleTasks(Interval)
		codesearch.ScheduleTasks(Interval)
	}
}

package util

import (
	"github.com/neal1991/gshark/util/githubsearch"

	"github.com/urfave/cli"

	"github.com/neal1991/gshark/logger"
	"time"
	"github.com/neal1991/gshark/util/appsearch"
	"strings"
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
		case "github":	 	// use go keyword or not
			logger.Log.Println("scan github code")
			githubsearch.ScheduleTasks(Interval)
		case "app":
			logger.Log.Println("scan app repos")
			appsearch.ScheduleTasks(Interval)
		case "all":
			logger.Log.Println("scan github code and app repos")
			go githubsearch.ScheduleTasks(Interval)
			appsearch.ScheduleTasks(Interval)
		default:
			logger.Log.Println("scan github code and app repos")
			go githubsearch.ScheduleTasks(Interval)
			appsearch.ScheduleTasks(Interval)
	}
}

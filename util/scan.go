package util

import (
	"x-patrol/tasks"
	"x-patrol/util/githubsearch"

	"github.com/urfave/cli"

	"strings"
	"time"
	"x-patrol/logger"
)

func Scan(ctx *cli.Context) {
	var ScanMode = "github"
	var Interval time.Duration = 900

	if ctx.IsSet("mode") {
		ScanMode = strings.ToLower(ctx.String("mode"))
	}

	if ctx.IsSet("time") {
		Interval = time.Duration(ctx.Int("time"))
	}

	switch ScanMode {
	case "github":
		logger.Log.Println("scan github code")
		githubsearch.ScheduleTasks(Interval)
	case "local":
		logger.Log.Println("scan local repos")
		tasks.ScheduleTasks(Interval)
	case "all":
		logger.Log.Println("scan github code and local repos")
		go githubsearch.ScheduleTasks(Interval)
		tasks.ScheduleTasks(Interval)
	default:
		logger.Log.Println("scan github code ")
		go githubsearch.ScheduleTasks(Interval)
	}
}

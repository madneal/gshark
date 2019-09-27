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
		for {
			githubsearch.RunTask(Interval)
		}
	case "app":
		logger.Log.Println("scan app results")
		for {
			appsearch.RunTask(Interval)
		}
	case "searchcode":
		logger.Log.Println("scan searchcode results")
		for {
			codesearch.RunTask(Interval)
		}
	case "all":
		logger.Log.Println("all scan mode")
		for {
			githubsearch.RunTask(Interval)
			codesearch.RunTask(Interval)
			appsearch.RunTask(Interval)
		}
	default:
		logger.Log.Println("default scan mode")
		for {
			githubsearch.RunTask(Interval)
			appsearch.RunTask(Interval)
			codesearch.RunTask(Interval)
		}
	}
}

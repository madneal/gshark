package search

import (
	"github.com/madneal/gshark/gobuster"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/search/appsearch"
	"github.com/madneal/gshark/search/codesearch"
	"github.com/madneal/gshark/search/githubsearch"
	"github.com/madneal/gshark/search/gitlabsearch"
	"github.com/urfave/cli"
	"strings"
	"time"
)

func Scan(ctx *cli.Context) {
	var scanMode string
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
	case "gitlab":
		logger.Log.Println("scan gitlab code")
		for {
			gitlabsearch.RunTask(Interval)
		}
	case "all":
		logger.Log.Println("all scan mode")
		for {
			gitlabsearch.RunTask(Interval)
			codesearch.RunTask(Interval)
			githubsearch.RunTask(Interval)
			appsearch.RunTask(Interval)
		}
	default:
		logger.Log.Println("default scan mode")
		for {
			go githubsearch.RunTask(Interval)
			go gobuster.RunTask(Interval)
		}
	}
}

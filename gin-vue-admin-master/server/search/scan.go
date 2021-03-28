package search

import (
	"gin-vue-admin/search/appsearch"
	"gin-vue-admin/search/codesearch"
	"gin-vue-admin/search/githubsearch"
	"gin-vue-admin/search/gitlabsearch"
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
		for {
			githubsearch.RunTask(Interval)
		}
	case "app":
		for {
			appsearch.RunTask(Interval)
		}
	case "searchcode":
		for {
			codesearch.RunTask(Interval)
		}
	case "gitlab":
		for {
			gitlabsearch.RunTask(Interval)
		}
	case "all":
		for {
			gitlabsearch.RunTask(Interval)
			codesearch.RunTask(Interval)
			githubsearch.RunTask(Interval)
			appsearch.RunTask(Interval)
		}
	default:
		for {
			githubsearch.RunTask(Interval)
			//go gobuster.RunTask(Interval)
		}
	}
}

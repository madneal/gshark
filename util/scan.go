package util

import (
	"github.com/neal1991/gshark/util/githubsearch"

	"github.com/urfave/cli"

	"github.com/neal1991/gshark/logger"
	"github.com/neal1991/gshark/util/appsearch"
	"time"
)

func Scan(ctx *cli.Context) {
	// seconds
	var Interval time.Duration = 900

	if ctx.IsSet("time") {
		Interval = time.Duration(ctx.Int("time"))
	}

	logger.Log.Println("scan github code ")
	// use go keyword or not
	githubsearch.ScheduleTasks(Interval)
	appsearch.ScheduleTasks(Interval)
}

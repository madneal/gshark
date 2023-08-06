package search

import (
	"github.com/madneal/gshark/search/codesearch"
	"github.com/madneal/gshark/search/githubsearch"
	"github.com/madneal/gshark/search/gitlabsearch"
	"github.com/madneal/gshark/search/gobuster"
	"github.com/madneal/gshark/search/postman"
	"time"
)

func ScanTask() {
	var Interval time.Duration = 900
	gitlabsearch.RunTask(Interval)
	codesearch.RunTask(Interval)
	githubsearch.RunTask(Interval)
	gobuster.RunTask(Interval)
	postman.RunTask()
}

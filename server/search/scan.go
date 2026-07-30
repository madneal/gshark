package search

import (
	"fmt"
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/search/codesearch"
	"github.com/madneal/gshark/search/githubsearch"
	"github.com/madneal/gshark/search/gitlabsearch"
	"github.com/madneal/gshark/search/gobuster"
	"github.com/madneal/gshark/search/postman"
)

// providerWatchdog bounds how long ScanTask waits for a single provider's
// RunTask before moving on. It must stay well above the ~900s sleep every
// RunTask already performs at the end of a normal pass, so healthy runs never
// trip it. A provider that exceeds this leaks its goroutine (Go cannot force-
// kill a running goroutine), but the scan loop itself is no longer blocked by it.
var providerWatchdog = 30 * time.Minute

func runProvider(name string, fn func()) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(providerWatchdog):
		global.GVA_LOG.Error(fmt.Sprintf("%s scan did not finish within %s, moving on to the next provider", name, providerWatchdog))
	}
}

func ScanTask() {
	for {
		if global.GVA_DB == nil {
			global.GVA_LOG.Info("have not init db")
			return
		}
		var Interval time.Duration = 900
		runProvider("gitlab", func() { gitlabsearch.RunTask(Interval) })
		runProvider("searchcode", func() { codesearch.RunTask(Interval) })
		runProvider("github", func() { githubsearch.RunTask(Interval) })
		runProvider("gobuster", func() { gobuster.RunTask(Interval) })
		runProvider("postman", func() { postman.RunTask() })
	}
}

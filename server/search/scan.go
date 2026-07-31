package search

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/search/codesearch"
	"github.com/madneal/gshark/search/githubsearch"
	"github.com/madneal/gshark/search/gitlabsearch"
	"github.com/madneal/gshark/search/gobuster"
	"github.com/madneal/gshark/search/postman"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

var (
	providerWatchdog        = 30 * time.Minute
	providerHeartbeat       = 10 * time.Second
	scanInterval            = 15 * time.Minute
	createScanCycle         = service.CreateScanCycle
	startScanLog            = service.StartScanLog
	heartbeatScanLog        = service.HeartbeatScanLog
	finishScanLog           = service.FinishScanLog
	markInterruptedScanLogs = service.MarkInterruptedScanLogs
	scanSleep               = time.Sleep
)

type providerTask struct {
	name string
	run  func() model.ScanOutcome
}

func runProvider(logID uint, name string, fn func() model.ScanOutcome) model.ScanOutcome {
	startedAt := time.Now()
	if err := startScanLog(logID, startedAt); err != nil {
		global.GVA_LOG.Error("start scan log failed", zap.String("provider", name), zap.Error(err))
	}

	done := make(chan model.ScanOutcome, 1)
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				global.GVA_LOG.Error("scan provider panicked", zap.String("provider", name), zap.Any("panic", recovered), zap.ByteString("stack", debug.Stack()))
				done <- model.ScanFailed(fmt.Sprintf("Provider panicked: %v", recovered))
			}
		}()
		done <- fn()
	}()

	watchdog := time.NewTimer(providerWatchdog)
	heartbeat := time.NewTicker(providerHeartbeat)
	defer watchdog.Stop()
	defer heartbeat.Stop()

	var outcome model.ScanOutcome
	for {
		select {
		case outcome = <-done:
			if outcome.Status == "" {
				outcome = model.ScanSuccess("Scan completed")
			}
			goto Finished
		case heartbeatAt := <-heartbeat.C:
			if err := heartbeatScanLog(logID, heartbeatAt); err != nil {
				global.GVA_LOG.Error("update scan heartbeat failed", zap.String("provider", name), zap.Error(err))
			}
		case <-watchdog.C:
			message := fmt.Sprintf("Scan did not finish within %s", providerWatchdog)
			global.GVA_LOG.Error(fmt.Sprintf("%s %s, moving on to the next provider", name, message))
			outcome = model.ScanOutcome{Status: model.ScanStatusTimeout, Message: message}
			goto Finished
		}
	}

Finished:
	finishedAt := time.Now()
	if err := finishScanLog(logID, outcome, startedAt, finishedAt); err != nil {
		global.GVA_LOG.Error("finish scan log failed", zap.String("provider", name), zap.Error(err))
	}
	return outcome
}

func scanTasks() []providerTask {
	return []providerTask{
		{name: "gitlab", run: gitlabsearch.RunTask},
		{name: "searchcode", run: codesearch.RunTask},
		{name: "github", run: githubsearch.RunTask},
		{name: "gobuster", run: gobuster.RunTask},
		{name: "postman", run: postman.RunTask},
	}
}

func runScanCycle() {
	tasks := scanTasks()
	providers := make([]string, 0, len(tasks))
	for _, task := range tasks {
		providers = append(providers, task.name)
	}
	cycleID := uuid.NewString()
	logIDs, err := createScanCycle(cycleID, providers)
	if err != nil {
		global.GVA_LOG.Error("create scan cycle failed", zap.String("cycleId", cycleID), zap.Error(err))
		logIDs = map[string]uint{}
	}
	for _, task := range tasks {
		runProvider(logIDs[task.name], task.name, task.run)
	}
}

func ScanTask() {
	if global.GVA_DB == nil {
		global.GVA_LOG.Info("have not init db")
		return
	}
	if err := markInterruptedScanLogs(time.Now()); err != nil {
		global.GVA_LOG.Error("mark interrupted scan logs failed", zap.Error(err))
	}
	for {
		runScanCycle()
		scanSleep(scanInterval)
	}
}

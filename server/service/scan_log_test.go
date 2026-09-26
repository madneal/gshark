package service

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/madneal/gshark/model"
)

func TestLimitScanLogMessagePreservesShortMessage(t *testing.T) {
	message := "provider completed"
	if got := limitScanLogMessage(message); got != message {
		t.Fatalf("limitScanLogMessage() = %q, want %q", got, message)
	}
}

func TestLimitScanLogMessageTruncatesAtUTF8Boundary(t *testing.T) {
	message := strings.Repeat("\u4e2d", maxScanLogMessageRunes+100)
	got := limitScanLogMessage(message)
	if count := utf8.RuneCountInString(got); count != maxScanLogMessageRunes {
		t.Fatalf("limited message uses %d characters, want %d", count, maxScanLogMessageRunes)
	}
	if !utf8.ValidString(got) {
		t.Fatal("limited message is not valid UTF-8")
	}
	if !strings.HasSuffix(got, scanLogTruncatedSuffix) {
		t.Fatalf("limited message does not end with %q", scanLogTruncatedSuffix)
	}
}

func TestBuildScanLogOverviewSummarizesCycle(t *testing.T) {
	now := time.Now()
	recent := now.Add(-time.Second)
	finished := now.Add(-2 * time.Second)
	logs := []model.ScanLog{
		{Provider: "gitlab", Status: model.ScanStatusSuccess, FinishedAt: &finished, HeartbeatAt: &finished},
		{Provider: "github", Status: model.ScanStatusRunning, HeartbeatAt: &recent},
		{Provider: "postman", Status: model.ScanStatusPending},
	}

	overview := buildScanLogOverview("cycle-1", logs, now)
	if overview.Total != 3 || overview.Completed != 1 || overview.Running != 1 {
		t.Fatalf("unexpected counts: %#v", overview)
	}
	if overview.Progress != 33 {
		t.Fatalf("progress = %d, want 33", overview.Progress)
	}
	if overview.Abnormal != 0 {
		t.Fatalf("abnormal = %d, want 0", overview.Abnormal)
	}
}

func TestBuildScanLogOverviewDetectsStaleAndFailedTasks(t *testing.T) {
	now := time.Now()
	staleHeartbeat := now.Add(-scanHeartbeatStaleAfter - time.Second)
	logs := []model.ScanLog{
		{Provider: "github", Status: model.ScanStatusRunning, HeartbeatAt: &staleHeartbeat},
		{Provider: "postman", Status: model.ScanStatusFailed},
	}

	overview := buildScanLogOverview("cycle-2", logs, now)
	if !overview.Logs[0].Stale {
		t.Fatal("expected running task with an old heartbeat to be marked stale")
	}
	if overview.Abnormal != 2 {
		t.Fatalf("abnormal = %d, want 2", overview.Abnormal)
	}
}

package model

import (
	"time"

	"github.com/madneal/gshark/global"
)

const (
	ScanStatusPending     = "pending"
	ScanStatusRunning     = "running"
	ScanStatusSuccess     = "success"
	ScanStatusSkipped     = "skipped"
	ScanStatusFailed      = "failed"
	ScanStatusTimeout     = "timeout"
	ScanStatusInterrupted = "interrupted"
)

type ScanLog struct {
	global.GVA_MODEL
	CycleID         string     `json:"cycleId" gorm:"column:cycle_id;type:varchar(36);index:idx_scan_log_cycle"`
	Provider        string     `json:"provider" gorm:"column:provider;type:varchar(32);index:idx_scan_log_provider"`
	Status          string     `json:"status" gorm:"column:status;type:varchar(16);index:idx_scan_log_status"`
	Message         string     `json:"message" gorm:"column:message;type:varchar(1000)"`
	ProgressCurrent int        `json:"progressCurrent" gorm:"column:progress_current;default:0"`
	ProgressTotal   int        `json:"progressTotal" gorm:"column:progress_total;default:1"`
	StartedAt       *time.Time `json:"startedAt" gorm:"column:started_at"`
	FinishedAt      *time.Time `json:"finishedAt" gorm:"column:finished_at"`
	HeartbeatAt     *time.Time `json:"heartbeatAt" gorm:"column:heartbeat_at;index:idx_scan_log_heartbeat"`
	DurationMs      int64      `json:"durationMs" gorm:"column:duration_ms;default:0"`
	Stale           bool       `json:"stale" gorm:"-"`
}

func (ScanLog) TableName() string {
	return "scan_log"
}

type ScanOutcome struct {
	Status  string
	Message string
}

func ScanSuccess(message string) ScanOutcome {
	return ScanOutcome{Status: ScanStatusSuccess, Message: message}
}

func ScanSkipped(message string) ScanOutcome {
	return ScanOutcome{Status: ScanStatusSkipped, Message: message}
}

func ScanFailed(message string) ScanOutcome {
	return ScanOutcome{Status: ScanStatusFailed, Message: message}
}

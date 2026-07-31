package service

import (
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"gorm.io/gorm"
)

const scanHeartbeatStaleAfter = 45 * time.Second

type ScanLogOverview struct {
	CycleID        string          `json:"cycleId"`
	Logs           []model.ScanLog `json:"logs"`
	Total          int             `json:"total"`
	Completed      int             `json:"completed"`
	Running        int             `json:"running"`
	Abnormal       int             `json:"abnormal"`
	Progress       int             `json:"progress"`
	LastActivityAt *time.Time      `json:"lastActivityAt"`
}

func CreateScanCycle(cycleID string, providers []string) (map[string]uint, error) {
	ids := make(map[string]uint, len(providers))
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, provider := range providers {
			logEntry := model.ScanLog{
				CycleID:       cycleID,
				Provider:      provider,
				Status:        model.ScanStatusPending,
				Message:       "Waiting for the previous provider to finish",
				ProgressTotal: 1,
			}
			if err := tx.Create(&logEntry).Error; err != nil {
				return err
			}
			ids[provider] = logEntry.ID
		}
		return nil
	})
	return ids, err
}

func StartScanLog(id uint, startedAt time.Time) error {
	if id == 0 {
		return nil
	}
	return global.GVA_DB.Model(&model.ScanLog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           model.ScanStatusRunning,
		"message":          "Scan is running",
		"started_at":       startedAt,
		"heartbeat_at":     startedAt,
		"progress_current": 0,
	}).Error
}

func HeartbeatScanLog(id uint, heartbeatAt time.Time) error {
	if id == 0 {
		return nil
	}
	return global.GVA_DB.Model(&model.ScanLog{}).Where("id = ? AND status = ?", id, model.ScanStatusRunning).
		Updates(map[string]interface{}{"heartbeat_at": heartbeatAt, "updated_at": heartbeatAt}).Error
}

func FinishScanLog(id uint, outcome model.ScanOutcome, startedAt, finishedAt time.Time) error {
	if id == 0 {
		return nil
	}
	return global.GVA_DB.Model(&model.ScanLog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           outcome.Status,
		"message":          outcome.Message,
		"finished_at":      finishedAt,
		"heartbeat_at":     finishedAt,
		"duration_ms":      finishedAt.Sub(startedAt).Milliseconds(),
		"progress_current": 1,
	}).Error
}

func MarkInterruptedScanLogs(interruptedAt time.Time) error {
	return global.GVA_DB.Model(&model.ScanLog{}).
		Where("status IN ?", []string{model.ScanStatusPending, model.ScanStatusRunning}).
		Updates(map[string]interface{}{
			"status":       model.ScanStatusInterrupted,
			"message":      "Scanner stopped before this task completed",
			"finished_at":  interruptedAt,
			"heartbeat_at": interruptedAt,
		}).Error
}

func GetScanLogInfoList(info request.ScanLogSearch) (error, interface{}, int64) {
	db := global.GVA_DB.Model(&model.ScanLog{})
	var logs []model.ScanLog
	if info.Provider != "" {
		db = db.Where("provider = ?", info.Provider)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.CycleID != "" {
		db = db.Where("cycle_id = ?", info.CycleID)
	}
	total, err := Paginate(db, info.Page, info.PageSize, &logs, "id desc")
	markStaleScanLogs(logs, time.Now())
	return err, logs, total
}

func GetScanLogOverview(now time.Time) (ScanLogOverview, error) {
	var latest model.ScanLog
	if err := global.GVA_DB.Order("id desc").First(&latest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ScanLogOverview{Logs: []model.ScanLog{}}, nil
		}
		return ScanLogOverview{}, err
	}

	var logs []model.ScanLog
	if err := global.GVA_DB.Where("cycle_id = ?", latest.CycleID).Order("id asc").Find(&logs).Error; err != nil {
		return ScanLogOverview{}, err
	}
	return buildScanLogOverview(latest.CycleID, logs, now), nil
}

func buildScanLogOverview(cycleID string, logs []model.ScanLog, now time.Time) ScanLogOverview {
	overview := ScanLogOverview{CycleID: cycleID, Logs: logs}
	markStaleScanLogs(overview.Logs, now)
	overview.Total = len(overview.Logs)
	for i := range overview.Logs {
		entry := overview.Logs[i]
		if entry.Status == model.ScanStatusRunning {
			overview.Running++
		}
		if entry.Status != model.ScanStatusPending && entry.Status != model.ScanStatusRunning {
			overview.Completed++
		}
		if entry.Stale || isAbnormalScanStatus(entry.Status) {
			overview.Abnormal++
		}
		activityAt := entry.UpdatedAt
		if entry.HeartbeatAt != nil {
			activityAt = *entry.HeartbeatAt
		}
		if overview.LastActivityAt == nil || activityAt.After(*overview.LastActivityAt) {
			activity := activityAt
			overview.LastActivityAt = &activity
		}
	}
	if overview.Total > 0 {
		overview.Progress = overview.Completed * 100 / overview.Total
	}
	return overview
}

func markStaleScanLogs(logs []model.ScanLog, now time.Time) {
	for i := range logs {
		if logs[i].Status != model.ScanStatusRunning || logs[i].HeartbeatAt == nil {
			continue
		}
		logs[i].Stale = now.Sub(*logs[i].HeartbeatAt) > scanHeartbeatStaleAfter
	}
}

func isAbnormalScanStatus(status string) bool {
	return status == model.ScanStatusFailed || status == model.ScanStatusTimeout || status == model.ScanStatusInterrupted
}

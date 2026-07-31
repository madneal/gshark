package request

type ScanLogSearch struct {
	PageInfo
	Provider string `json:"provider" form:"provider"`
	Status   string `json:"status" form:"status"`
	CycleID  string `json:"cycleId" form:"cycleId"`
}

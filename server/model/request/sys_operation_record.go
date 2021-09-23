package request

import "github.com/madneal/gshark/model"

type SysOperationRecordSearch struct {
	model.SysOperationRecord
	PageInfo
}

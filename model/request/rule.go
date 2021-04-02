package request

import "github.com/madneal/gshark/model"

type RuleSearch struct {
	model.Rule
	PageInfo
}

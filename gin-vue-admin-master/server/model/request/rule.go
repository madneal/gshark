package request

import "gin-vue-admin/model"

type RuleSearch struct {
	model.Rule
	PageInfo
}

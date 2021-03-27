package request

import "gin-vue-admin/model"

type SubdomainSearch struct {
	model.Subdomain
	PageInfo
}

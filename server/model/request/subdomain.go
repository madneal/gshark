package request

import "github.com/madneal/gshark/model"

type SubdomainSearch struct {
	model.Subdomain
	PageInfo
}

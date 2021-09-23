package response

import "github.com/madneal/gshark/model"

type SysAuthorityResponse struct {
	Authority model.SysAuthority `json:"authority"`
}

type SysAuthorityCopyResponse struct {
	Authority      model.SysAuthority `json:"authority"`
	OldAuthorityId string             `json:"oldAuthorityId"`
}

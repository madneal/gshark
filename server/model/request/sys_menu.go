package request

import "github.com/madneal/gshark/model"

type AddMenuAuthorityInfo struct {
	Menus       []model.SysBaseMenu
	AuthorityId string
}

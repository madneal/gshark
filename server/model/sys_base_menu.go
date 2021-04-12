package model

import (
	"github.com/madneal/gshark/global"
)

type SysBaseMenu struct {
	global.GVA_MODEL
	MenuLevel     uint   `json:"-"`
	ParentId      string `json:"parentId" gorm:"comment:父菜单ID"`
	Path          string `json:"path" gorm:"comment:路由path"`
	Name          string `json:"name" gorm:"comment:路由name"`
	Hidden        bool   `json:"hidden" gorm:"comment:是否在列表隐藏"`
	Component     string `json:"component" gorm:"comment:对应前端文件路径"`
	Sort          int    `json:"sort" gorm:"comment:排序标记"`
	Meta          `json:"meta" gorm:"comment:附加属性"`
	SysAuthoritys []SysAuthority         `json:"authoritys" gorm:"many2many:sys_authority_menus;"`
	Children      []SysBaseMenu          `json:"children" gorm:"-"`
	Parameters    []SysBaseMenuParameter `json:"parameters"`
}

type Meta struct {
	KeepAlive   bool   `json:"keepAlive" gorm:"comment:是否缓存"`
	DefaultMenu bool   `json:"defaultMenu" gorm:"comment:是否是基础路由（开发中）"`
	Title       string `json:"title" gorm:"comment:菜单名"`
	Icon        string `json:"icon" gorm:"comment:菜单图标"`
	CloseTab    bool   `json:"closeTab" gorm:"comment:自动关闭tab"`
}

type SysBaseMenuParameter struct {
	global.GVA_MODEL
	SysBaseMenuID uint
	Type          string `json:"type" gorm:"comment:地址栏携带参数为params还是query"`
	Key           string `json:"key" gorm:"comment:地址栏携带参数的key"`
	Value         string `json:"value" gorm:"comment:地址栏携带参数的值"`
}

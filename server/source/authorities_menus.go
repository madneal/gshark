package source

import (
	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"gorm.io/gorm"
)

var AuthoritiesMenus = new(authoritiesMenus)

type authoritiesMenus struct{}

type AuthorityMenus struct {
	AuthorityId string `gorm:"column:sys_authority_authority_id"`
	BaseMenuId  uint   `gorm:"column:sys_base_menu_id"`
}

var authorityMenus = []AuthorityMenus{
	{"888", 1},
	{"888", 2},
	{"888", 3},
	{"888", 4},
	{"888", 5},
	{"888", 6},
	{"888", 7},
	{"888", 8},
	{"888", 9},
	{"888", 10},
	{"888", 11},
	{"888", 12},
	{"888", 13},
	{"888", 14},
	{"888", 15},
	{"888", 16},
	{"888", 17},
	{"888", 18},
	{"888", 19},
	{"888", 20},
	{"888", 21},
	{"888", 22},
	{"888", 23},
	{"888", 24},
	{"888", 25},
	{"888", 26},
	{"888", 27},
	{"888", 28},
	{"888", 29},
	{"888", 30},
	{"888", 31},
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: sys_authority_menus 表数据初始化
func (a *authoritiesMenus) Init() error {
	return global.GVA_DB.Table("sys_authority_menus").Transaction(func(tx *gorm.DB) error {
		if tx.Where("sys_authority_authority_id IN ('888', '8881', '9528')").Find(&[]AuthorityMenus{}).RowsAffected == 48 {
			color.Danger.Println("\n[Mysql] --> sys_authority_menus 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&authorityMenus).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> sys_authority_menus 表初始数据成功!")
		return nil
	})
}

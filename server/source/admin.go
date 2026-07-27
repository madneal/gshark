package source

import (
	"time"

	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var Admin = new(admin)

type admin struct{}

// Defaults used when init does not pass custom credentials.
var (
	AdminUsername = "gshark"
	AdminPassword = "gshark" // plain text; hashed in Init
	AdminNickName = "超级管理员"
)

// SetAdminCredentials sets the admin seed account for Init().
// Empty values keep the previous/default value.
func SetAdminCredentials(username, password string) {
	if username != "" {
		AdminUsername = username
	}
	if password != "" {
		AdminPassword = password
	}
}

func (a *admin) Init() error {
	user := model.SysUser{
		GVA_MODEL:   global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		UUID:        uuid.NewV4(),
		Username:    AdminUsername,
		Password:    utils.MD5V([]byte(AdminPassword)),
		NickName:    AdminNickName,
		HeaderImg:   "https://s2.loli.net/2023/12/29/FRNA23eJcXjDM4Y.jpg",
		AuthorityId: "888",
	}
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []int{1, 2}).Find(&[]model.SysUser{}).RowsAffected == 2 {
			color.Danger.Println("\n[Mysql] --> sys_users 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		color.Info.Printf("\n[Mysql] --> sys_users 初始管理员: %s\n", AdminUsername)
		return nil
	})
}

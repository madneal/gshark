package models

import (
	"github.com/madneal/gshark/util"
)

type Admin struct {
	Id       int64
	Username string
	Password string
	Role     string
}

func NewAdmin(username, password, role string) *Admin {
	encryptPass := util.MakeMd5(password)
	return &Admin{Username: username, Password: encryptPass, Role: role}
}

func (u *Admin) Insert() (int64, error) {
	return Engine.Insert(u)
}

func ListAdmins() ([]Admin, error) {
	admins := make([]Admin, 0)
	err := Engine.Table("admin").Find(&admins)
	return admins, err
}

func GetAdminById(id int64) (*Admin, bool, error) {
	admin := new(Admin)
	has, err := Engine.ID(id).Get(admin)
	return admin, has, err
}

func EditAdminById(id int64, username, password, role string) error {
	admin := new(Admin)
	has, err := Engine.ID(id).Get(admin)
	if err == nil && has {
		admin.Username = username
		admin.Password = util.MakeMd5(password)
		admin.Role = role
		Engine.ID(id).Update(admin)
	}
	return err
}

func DeleteAdminById(id int64) error {
	admin := new(Admin)
	_, err := Engine.ID(id).Delete(admin)
	return err
}

func Auth(username, password string) (bool, string, error) {
	admin := new(Admin)
	encryptPass := util.MakeMd5(password)
	has, err := Engine.Table("admin").
		Where("username = ? and password = ?", username, encryptPass).Get(admin)
	var role string
	if has {
		role = admin.Role
	}
	return has, role, err
}

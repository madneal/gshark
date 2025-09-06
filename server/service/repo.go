package service

import (
	"errors"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"gorm.io/gorm"
)

func CreateRepo(repo model.Repo) (err error) {
	err = global.GVA_DB.Create(&repo).Error
	return err
}

func DeleteRepo(repo model.Repo) (err error) {
	err = global.GVA_DB.Delete(&repo).Error
	return err
}

func DeleteRepoByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.Repo{}, "id in ?", ids.Ids).Error
	return err
}

func UpdateRepo(repo model.Repo) (err error) {
	err = global.GVA_DB.Save(&repo).Error
	return err
}

func GetRepo(id uint) (err error, repo model.Repo) {
	err = global.GVA_DB.Where("id = ?", id).First(&repo).Error
	return
}

func GetRepoInfoList(info request.RepoSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.Repo{})
	var repos []model.Repo
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&repos).Error
	return err, repos, total
}

func GetRepoByType(typeStr string) (err error, repos []model.Repo) {
	db := global.GVA_DB.Model(&model.Repo{})
	err = db.Where("type = ?", typeStr).Find(&repos).Error
	return err, repos
}

func CheckRepoExist(repo *model.Repo) (err error, result bool) {
	r := global.GVA_DB.Where("project_id = ?", repo.ProjectId).First(repo)
	err = r.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err, false
	}
	return err, true
}

package service

import (
	"errors"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"gorm.io/gorm"
)

func CreateRepo(repo model.Repo) (err error) {
	return Create(&repo)
}

func DeleteRepo(repo model.Repo) (err error) {
	return Delete(&repo)
}

func DeleteRepoByIds(ids request.IdsReq) (err error) {
	return DeleteByIds[model.Repo](ids)
}

func UpdateRepo(repo model.Repo) (err error) {
	return Update(&repo)
}

func GetRepo(id uint) (err error, repo model.Repo) {
	repo, err = GetByID[model.Repo](id)
	return
}

func GetRepoInfoList(info request.RepoSearch) (err error, list interface{}, total int64) {
	db := global.GVA_DB.Model(&model.Repo{})
	var repos []model.Repo
	total, err = Paginate(db, info.Page, info.PageSize, &repos, "")
	return err, repos, total
}

func GetRepoByType(typeStr string) (err error, repos []model.Repo) {
	db := global.GVA_DB.Model(&model.Repo{})
	err = db.Where("type = ?", typeStr).Find(&repos).Error
	return err, repos
}

// ResetRepoStatusByType marks every repo of the given type as unsearched again,
// so a bounded per-cycle crawl (e.g. gitlabsearch's project crawl) can recycle
// its pool once every known repo has been searched.
func ResetRepoStatusByType(typeStr string) (err error) {
	return global.GVA_DB.Model(&model.Repo{}).Where("type = ?", typeStr).Update("status", 0).Error
}

func CheckRepoExist(repo *model.Repo) (err error, result bool) {
	r := global.GVA_DB.Where("project_id = ?", repo.ProjectId).First(repo)
	err = r.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err, false
	}
	return err, true
}

package service

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"go.uber.org/zap"
)

func CreateSearchResult(searchResult model.SearchResult) (err error) {
	return Create(&searchResult)
}

func DeleteSearchResult(searchResult model.SearchResult) (err error) {
	return Delete(&searchResult)
}

func DeleteSearchResultByIds(ids request.IdsReq) (err error) {
	return DeleteByIds[model.SearchResult](ids)
}

func UpdateSearchResultByIds(req request.BatchUpdateReq) (err error) {
	err = global.GVA_DB.Table("search_result").Where("id in ?", req.Ids).
		UpdateColumn("status", req.Status).Error
	return err
}

func IgnoreResultsByRepo(repo string) (err error) {
	err = global.GVA_DB.Table("search_result").
		Where("status = ? and repo = ? and (sec_keyword is null or sec_keyword = '')", global.UnhandledStatus, repo).
		UpdateColumn("status", global.IgnoredStatus).Error
	return err
}

func UpdateSearchResult(updateReq request.UpdateReq) (err error) {
	err = global.GVA_DB.Table("search_result").Where("repo = ?", updateReq.Repo).
		UpdateColumn("status", updateReq.Status).Error
	return err
}

func UpdateSearchResultById(id, status int) (err error) {
	err = global.GVA_DB.Table("search_result").Where("id = ?", id).
		UpdateColumn("status", status).Error
	return err
}

func GetSearchResult(id uint) (err error, searchResult model.SearchResult) {
	searchResult, err = GetByID[model.SearchResult](id)
	return
}

func ListSearchResultByStatus(status int) (err error, list []model.SearchResult) {
	err = global.GVA_DB.Where("status = ?", status).Find(&list).Error
	return err, list
}

func GetSearchResultInfoList(info request.SearchResultSearch) (err error, list interface{}, total int64) {
	db := global.GVA_DB.Model(&model.SearchResult{})
	var searchResults []model.SearchResult
	if info.Query != "" {
		db = db.Where("`repo` LIKE ? or `text_matches_json` LIKE ?",
			"%"+info.Query+"%", "%"+info.Query+"%")
	}
	if info.Keyword != "" {
		db = db.Where("`keyword` = ? or `sec_keyword` = ?", info.Keyword, info.Keyword)
	}
	if info.Status >= 0 {
		db = db.Where("`status` = ?", info.Status)
	}
	if info.OnlySecKeyword {
		db = db.Where("`sec_keyword` != ''")
	}
	total, err = Paginate(db, info.Page, info.PageSize, &searchResults, "id desc")
	return err, searchResults, total
}

func CheckExistOfSearchResult(searchResult *model.SearchResult) bool {
	urlExist := searchResult.CheckPathExists()
	repoExists := searchResult.CheckRepoExists()
	return urlExist || repoExists
}

func SaveSearchResults(searchResults []model.SearchResult) int {
	var insertCount int
	for _, result := range searchResults {
		exist := CheckExistOfSearchResult(&result)
		if exist {
			continue
		}
		err := CreateSearchResult(result)
		if err != nil {
			global.GVA_LOG.Error("save search result error", zap.Any("save searchResult error",
				err))
		} else {
			insertCount++
		}
	}
	return insertCount
}

func SaveSearchResultPointers(searchResults []*model.SearchResult, keyword string) int {
	results := make([]model.SearchResult, 0, len(searchResults))
	for _, result := range searchResults {
		if result == nil {
			continue
		}
		if keyword != "" {
			result.Keyword = keyword
		}
		results = append(results, *result)
	}
	return SaveSearchResults(results)
}

func GetReposByStatus(status int) (error, []string) {
	var results []model.SearchResult
	err := global.GVA_DB.Distinct().Select("repo").Where("status = ?",
		status).Find(&results).Error
	repos := make([]string, 0)
	if err != nil {
		return err, repos
	}
	for _, result := range results {
		repos = append(repos, result.Repo)
	}
	return err, repos
}

func GetKeywordByRepo(repo string) (string, error) {
	var result model.SearchResult
	err := global.GVA_DB.Where("repo = ?", repo).First(&result).Error
	return result.Keyword, err
}

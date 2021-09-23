package service

import (
	"errors"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"gorm.io/gorm"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateSearchResult
//@description: 创建SearchResult记录
//@param: searchResult model.SearchResult
//@return: err error

func CreateSearchResult(searchResult model.SearchResult) (err error) {
	err = global.GVA_DB.Create(&searchResult).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSearchResult
//@description: 删除SearchResult记录
//@param: searchResult model.SearchResult
//@return: err error

func DeleteSearchResult(searchResult model.SearchResult) (err error) {
	err = global.GVA_DB.Delete(&searchResult).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSearchResultByIds
//@description: 批量删除SearchResult记录
//@param: ids request.IdsReq
//@return: err error

func DeleteSearchResultByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.SearchResult{}, "id in ?", ids.Ids).Error
	return err
}

func UpdateSearchResultByIds(req request.BatchUpdateReq) (err error) {
	err = global.GVA_DB.Table("search_result").Where("id in ?", req.Ids).
		UpdateColumn("status", req.Status).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateSearchResult
//@description: 更新SearchResult记录
//@param: searchResult *model.SearchResult
//@return: err error

func UpdateSearchResult(updateReq request.UpdateReq) (err error) {
	err = global.GVA_DB.Table("search_result").Where("repo = ?", updateReq.Repo).
		UpdateColumn("status", updateReq.Status).Error
	return err
}

func UpdateSearchResultById(id, status int) (err error) {
	err = global.GVA_DB.UpdateColumn("status", status).Where("id = ?", id).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSearchResult
//@description: 根据id获取SearchResult记录
//@param: id uint
//@return: err error, searchResult model.SearchResult

func GetSearchResult(id uint) (err error, searchResult model.SearchResult) {
	err = global.GVA_DB.Where("id = ?", id).First(&searchResult).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSearchResultInfoList
//@description: 分页获取SearchResult记录
//@param: info request.SearchResultSearch
//@return: err error, list interface{}, total int64

func GetSearchResultInfoList(info request.SearchResultSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.SearchResult{})
	var searchResults []model.SearchResult
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Query != "" {
		db = db.Where("`repo` LIKE ? or `text_matches_json` LIKE ?", "%"+info.Query+"%", "%"+info.Query+"%")
	}
	if info.Keyword != "" {
		db = db.Where("`keyword` = ?", info.Keyword)
	}
	if info.Status >= 0 {
		db = db.Where("`status` = ?", info.Status)
	}
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&searchResults).Error
	return err, searchResults, total
}

func CheckExistOfSearchResult(searchResult *model.SearchResult) (err error, result bool) {
	queryResult := global.GVA_DB.Where("url = ? or repo = ? and status != 0", searchResult.Url, searchResult.Repo).First(searchResult)
	err = queryResult.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err, false
	}
	return err, true
}

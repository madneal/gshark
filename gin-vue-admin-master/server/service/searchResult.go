package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
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
	err = global.GVA_DB.Delete(&[]model.SearchResult{},"id in ?",ids.Ids).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateSearchResult
//@description: 更新SearchResult记录
//@param: searchResult *model.SearchResult
//@return: err error

func UpdateSearchResult(searchResult model.SearchResult) (err error) {
	err = global.GVA_DB.Save(&searchResult).Error
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
    if info.Repo != "" {
        db = db.Where("`repo` LIKE ?","%"+ info.Repo+"%")
    }
    if info.Keyword != "" {
        db = db.Where("`keyword` = ?",info.Keyword)
    }
    if info.Status != 0 {
        db = db.Where("`status` = ?",info.Status)
    }
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&searchResults).Error
	return err, searchResults, total
}
package models

import (
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/util"
	"github.com/madneal/gshark/vars"
	"time"
)

// AppSearchResult represents a single search result for app market search
type AppSearchResult struct {
	Id          int64
	Name        *string `json:"name,omitempty"`
	Description *string
	Market      *string `json:"market,omitempty"`
	Developer   *string
	Version     *string
	DeployDate  *string
	AppUrl      *string
	Status      int
	CreatedTime time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
}

func (r *AppSearchResult) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *AppSearchResult) Exist() (bool, error) {
	return Engine.Table("app_search_result").Where("name=? and market=?",
		r.Name, r.Market).Exist()
}

func ListAppSearchResultByPage(page int, status int) ([]AppSearchResult, int, int) {
	results := make([]AppSearchResult, 0)
	totalPages, err := Engine.Table("app_search_result").Where("status=?", status).Count()
	var pages int

	page, pages = util.GetPageAndPagesByCount(page, int(totalPages))

	err = Engine.Where("status=?", status).
		Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Desc("id").Find(&results)

	if err != nil {
		logger.Log.Errorf("search failed:%s", err)
	}

	return results, pages, int(totalPages)
}

func ConfirmAppResult(id int64) (err error) {
	err = changeAppResultStatus(1, id)
	return err
}

func IgnoreAppSearchResult(id int64) (err error) {
	err = changeAppResultStatus(2, id)
	return err
}

func changeAppResultStatus(status int, id int64) (err error) {
	_, err = Engine.Table("app_search_result").Exec("update app_search_result set status=?  "+
		"where id=?", status, id)
	return err
}

package models

import "time"

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
	CreatedTime time.Time
	UpdatedTime time.Time
}

func (r *AppSearchResult) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *AppSearchResult) Exist() (bool, error) {
	return Engine.Table("app_search_result").Where("name=? and market=?",
		r.Name, r.Market).Exist()
}

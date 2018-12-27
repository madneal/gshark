package models

import "time"

// APPSearchResult represents a single search result for app market search
type APPSearchResult struct {
	Id          int64
	Name        *string `json:"name,omitempty"`
	Description *string
	Market      *string `json:"market,omitempty"`
	CreatedTime time.Time
	UpdatedTime time.Time
	Developer   *string
	Version     *string
	DeployDate  *string
	APPUrl      *string
	Status      int
}

func (r *APPSearchResult) Insert() (int64, error) {
	return Engine.Insert(r)
}

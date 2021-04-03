// 自动生成模板SearchResult
package model

import (
	"github.com/google/go-github/github"
	"github.com/madneal/gshark/global"
)

// 如果含有time.Time 请自行import time包
type SearchResult struct {
	global.GVA_MODEL
	Repo         string `json:"repo" form:"repo" gorm:"column:repo;comment:;type:varchar(50);size:50;"`
	Repository   *github.Repository
	Matches      string      `json:"matches" form:"matches" gorm:"column:matches;comment:;type:blob;"`
	Keyword      string      `json:"keyword" form:"keyword" gorm:"column:keyword;comment:;type:varchar(100);size:100;"`
	Path         string      `json:"path" form:"path" gorm:"column:path;comment:;type:varchar(100);size:100;"`
	Url          string      `json:"url" form:"url" gorm:"column:url;comment:;type:varchar(500);size:500;"`
	TextmatchMd5 string      `json:"textmatchMd5" form:"textmatchMd5" gorm:"column:textmatch_md5;comment:;type:varchar(100);size:100;"`
	Status       int         `json:"status" form:"status" gorm:"column:status;comment:;type:int;size:3;"`
	TextMatches  []TextMatch `json:"text_matches,omitempty" gorm:"LONGBLOB unique"`
}

// TextMatch represents a text match for a SearchResult
type TextMatch struct {
	Id         int64
	ObjectURL  *string `json:"object_url,omitempty"`
	ObjectType *string `json:"object_type,omitempty"`
	Property   *string `json:"property,omitempty"`
	Fragment   *string `json:"fragment,omitempty"`
	Matches    []Match `xorm:"LONGBLOB"`
}

type Match struct {
	Id      int64
	Text    *string `json:"text,omitempty" xorm:"LONGBLOB"`
	Indices []int   `json:"indices,omitempty" xorm:"json"`
}

func (SearchResult) TableName() string {
	return "search_result"
}

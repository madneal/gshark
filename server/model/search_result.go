package model

import (
	"errors"
	"github.com/madneal/gshark/global"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SearchResult struct {
	global.GVA_MODEL
	Repo            string         `json:"repo" form:"repo" gorm:"column:repo;comment:;type:varchar(200);size:200;"`
	RepoUrl         string         `gorm:"column:repository;type:varchar(200);"`
	Matches         string         `json:"matches" form:"matches" gorm:"column:matches;comment:;type:text;"`
	Keyword         string         `json:"keyword" form:"keyword" gorm:"column:keyword;comment:;type:varchar(100);size:100;"`
	SecKeyword      string         `json:"sec_keyword" form:"sec_keyword" gorm:"column:sec_keyword;comment:;type:varchar(100);size:100;"`
	Path            string         `json:"path" form:"path" gorm:"column:path;comment:;type:varchar(500);size:100;"`
	Url             string         `json:"url" form:"url" gorm:"column:url;comment:;type:varchar(500);size:500;"`
	Status          int            `json:"status" form:"status" gorm:"column:status;comment:;type:int;size:3;"`
	TextMatchesJson datatypes.JSON `json:"text_matches,omitempty" gorm:"type:json;"`
}

type TextMatchesJson struct {
	TextMatch []TextMatch
}

// TextMatch represents a text match for a SearchResult
type TextMatch struct {
	Id         int64
	ObjectURL  *string `json:"object_url,omitempty"`
	ObjectType *string `json:"object_type,omitempty"`
	Property   *string `json:"property,omitempty"`
	Fragment   *string `json:"fragment,omitempty"`
	Matches    []Match `gorm:"json"`
}

type Match struct {
	Id      int64
	Text    *string `json:"text,omitempty" gorm:"json"`
	Indices []int   `json:"indices,omitempty" gorm:"json"`
}

func (SearchResult) TableName() string {
	return "search_result"
}

func (result *SearchResult) CheckPathExists() bool {
	var r SearchResult
	queryResult := global.GVA_DB.Table("search_result").Where("path = ? and repository = ?",
		result.Path, result.RepoUrl).First(&r)
	err := queryResult.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		global.GVA_LOG.Error("CheckPathExists error", zap.Any("err", err))
	}
	return true
}

func (result *SearchResult) CheckRepoExists() bool {
	var r SearchResult
	queryResult := global.GVA_DB.Table("search_result").Where("repo = ? and status > ?",
		result.Repo, global.UnhandledStatus).First(&r)
	err := queryResult.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		global.GVA_LOG.Error("CheckRepoExists err", zap.Error(err))
	}
	return true
}

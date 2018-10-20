package models

import (
	"fmt"
	"github.com/google/go-github/github"
	"time"
	"x-patrol/vars"
)

type Match struct {
	Id      int64
	Text    *string `json:"text,omitempty" xorm:"LONGBLOB"`
	Indices []int   `json:"indices,omitempty" xorm:"json"`
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

// CodeResult represents a single search result.
type CodeResult struct {
	Id          int64
	Name        *string `json:"name,omitempty"`
	Path        *string `json:"path,omitempty"`
	RepoName    string
	SHA         *string            `json:"sha,omitempty" xorm:"sha"`
	HTMLURL     *string            `json:"html_url,omitempty" xorm:"html_url"`
	Repository  *github.Repository `json:"repository,omitempty" xorm:"json"`
	TextMatches []TextMatch        `json:"text_matches,omitempty" xorm:"LONGBLOB"`
	Status      int                // 1 confirmed 2 ignored
	Version     int                `xorm:"version"`
	CreatedTime time.Time          `xorm:"created"`
	UpdatedTime time.Time          `xorm:"updated"`
	RepoPath    *string
	Keyword     *string
}

type MatchedText struct {
	Keyword    *string
	StartIndex int
	EndIndex   int
	Text       *string
}

type CodeResultDetail struct {
	Id int64
	// owner
	OwnerName      *string
	OwnerURl       *string
	Company        *string
	Location       *string
	Email          *string
	Blog           *string
	OwnerCreatedAt string
	Type           *string
	// repo
	RepoName      *string
	RepoUrl       *string
	Lang          *string
	Keyword       *string
	RepoCreatedAt *github.Timestamp
	RepoUpdatedAt *github.Timestamp

	Status       int
	MatchedTexts []*string
}

// CodeSearchResult represents the result of a code search.
type CodeSearchResult struct {
	Total             *int         `json:"total_count,omitempty"`
	IncompleteResults *bool        `json:"incomplete_results,omitempty"`
	CodeResults       []CodeResult `json:"items,omitempty" xorm:"json"`
}

func (r *CodeResult) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *CodeResult) Exist() (bool, error) {
	var c CodeResult
	return Engine.Table("code_result").Where("name=? and sha=?", r.Name, r.SHA).Get(&c)
}

func ListGithubSearchResultPage(page int, status int) ([]CodeResult, int, int) {
	results := make([]CodeResult, 0)
	totalPages, err := Engine.Table("code_result").Where("status=?", status).Count()
	var pages int

	if int(totalPages)%vars.PAGE_SIZE == 0 {
		pages = int(totalPages) / vars.PAGE_SIZE
	} else {
		pages = int(totalPages)/vars.PAGE_SIZE + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}

	err = Engine.Where("status=?", status).Omit("repository").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Desc("id").Find(&results)

	if err != nil {
		fmt.Errorf("search failed:%s", err)
	}

	return results, pages, int(totalPages)
}

func GetPageById(id int64) (int, error) {
	var page int
	result := make([]int64, 0)
	err := Engine.Table("code_result").Cols("id").Where("status=0").Find(&result)
	for i, value := range result {
		if value == id {
			page = ((i + 1) / vars.PAGE_SIZE) + 1
			if page == 0 {
				page = 1
			}
			break
		}
	}
	return page, err
}

func GetCodeResultDetailById(id int64) (*CodeResultDetail, error) {
	codeResultDetail := CodeResultDetail{Id: id}
	has, err := Engine.Table("code_result_detail").ID(id).Get(&codeResultDetail)

	if err == nil && !has {
		omitRepo := false
		_, codeResult, _ := GetReportById(id, omitRepo)
		codeResultDetail = getCodeResultDetail(codeResult)
	}
	return &codeResultDetail, err
}

func getCodeResultDetail(codeResult *CodeResult) CodeResultDetail {
	detail := CodeResultDetail{}
	repo := *codeResult.Repository
	owner := *codeResult.Repository.Owner

	detail.OwnerName = owner.Login
	detail.RepoName = repo.FullName
	detail.RepoUrl = repo.HTMLURL
	detail.Lang = repo.Language
	detail.RepoCreatedAt = repo.CreatedAt
	detail.RepoUpdatedAt = repo.UpdatedAt
	detail.Status = codeResult.Status
	detail.Keyword = codeResult.Keyword

	detail.MatchedTexts = getMatchedTests(codeResult)
	return detail
}

func getMatchedTests(result *CodeResult) []*string {
	textMatches := result.TextMatches
	matchedTexts := make([]*string, 0)
	for _, textMatch := range textMatches {
		matchedTexts = append(matchedTexts, textMatch.Fragment)
	}
	return matchedTexts
}

func GetReportById(id int64, omitRepo bool) (bool, *CodeResult, error) {
	report := new(CodeResult)
	var has bool
	var err error
	if omitRepo {
		has, err = Engine.Id(id).Omit("repository").Get(report)
	} else {
		has, err = Engine.ID(id).Get(report)
	}

	return has, report, err
}

// confirm the whole repository by id
func ConfirmResultById(id int64) (err error) {
	report := new(CodeResult)
	has, err := Engine.Id(id).Get(report)
	if err == nil && has {
		err = ChangeReportsStatusByRepo(id, 1)
	}
	return err
}

func CancelReportById(id int64) (page int, err error) {
	report := new(CodeResult)
	has, err := Engine.Id(id).Omit("repository").Get(report)
	page, err = GetPageById(id)
	if err == nil && has {
		report.Status = 2
		_, err = Engine.Id(id).Cols("status").Update(report)
	}
	return page, err
}

func CancelReportsByRepo(id int64) (err error) {
	var repo string
	has, err := Engine.Table("code_result").Where("id =?", id).Cols("repo_name").Get(&repo)
	if err == nil && has {
		_, err = Engine.Table("code_result").Exec("update code_result set status = 2 where repo_name=?", repo)
	}
	return err
}

func ChangeReportsStatusByRepo(id int64, status int) (err error) {
	var repo string
	has, err := Engine.Table("code_result").Where("id=?", id).Cols("repo_name").Get(&repo)
	if err == nil && has {
		_, err = Engine.Table("code_result").Exec("update code_result set status=? where repo_name=?",
			status, repo)
	}
	return err
}

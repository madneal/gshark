package githubsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"github.com/madneal/gshark/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"

	"time"
)

func Search(rules []model.Rule) {
	client, err := GetGithubClient()
	if err != nil {
		global.GVA_LOG.Error("GetGithubClient err", zap.Error(err))
		return
	}
	var content string
	var counts int
	for _, rule := range rules {
		query, err := BuildQuery(rule.Content)
		if err != nil {
			global.GVA_LOG.Error("BuildQuery error", zap.Error(err))
			continue
		}
		results, err := client.SearchCode(query)
		if err != nil {
			global.GVA_LOG.Error("SearchCode error", zap.Error(err))
			continue
		}
		counts = SaveResult(results, rule.Content, "")
		if counts > 0 {
			content += fmt.Sprintf("%s: %d条<br>", rule.Content, counts)
		}
	}
	if content != "" {
		if global.GVA_CONFIG.Email.Enable {
			err = utils.EmailSend("Github敏感信息报告", content)
			if err != nil {
				global.GVA_LOG.Error("send email error", zap.Any("err", err))
			}
		}
		if global.GVA_CONFIG.Wechat.Enable {
			content = "Github敏感信息报告\n" + content
			err = utils.BotSend(content)
			if err != nil {
				global.GVA_LOG.Error("send wechat error", zap.Any("err", err))
			}
		}
	}
}

func SaveResult(results []*github.CodeSearchResult, keyword, secKeyword string) int {
	searchResults := ConvertToSearchResults(results, keyword, secKeyword)
	insertCount := service.SaveSearchResults(searchResults)
	return insertCount
}

func ConvertToSearchResults(results []*github.CodeSearchResult, keyword, secKeyword string) []model.SearchResult {
	searchResults := make([]model.SearchResult, 0)
	for _, result := range results {
		codeResults := result.CodeResults
		for _, codeResult := range codeResults {
			searchResult := model.SearchResult{
				RepoUrl:    *codeResult.Repository.HTMLURL,
				Repo:       *codeResult.Repository.FullName,
				Keyword:    keyword,
				SecKeyword: secKeyword,
				Url:        *codeResult.HTMLURL,
				Path:       *codeResult.Path,
				Status:     0,
			}
			if len(codeResult.TextMatches) > 0 {
				b, err := json.Marshal(codeResult.TextMatches)
				searchResult.TextMatchesJson = b
				if err != nil {
					global.GVA_LOG.Error("json.marshal error", zap.Error(err))
				}
			}
			searchResults = append(searchResults, searchResult)
		}
	}
	return searchResults
}

func RunTask(duration time.Duration) {
	err, rules := service.GetValidRulesByType("github")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType github err", zap.Error(err))
		return
	}
	Search(rules)
	global.GVA_LOG.Info(fmt.Sprintf("Comple the scan of Github, start to sleep %d seconds", duration))
	time.Sleep(duration * time.Second)
}

func (c *Client) GetCommiter(ctx context.Context, owner, repo string) string {
	commit, _, err := c.Client.Repositories.GetCommit(ctx, owner, repo, "master")
	if err != nil {
		global.GVA_LOG.Error("get github commit err", zap.Error(err))
		return ""
	}
	return commit.Commit.Committer.GetEmail()
}

func (c *Client) SearchCode(query string) ([]*github.CodeSearchResult, error) {
	var allSearchResult []*github.CodeSearchResult
	var err error
	ctx := context.Background()
	listOpt := github.ListOptions{PerPage: 100}
	opt := &github.SearchOptions{TextMatch: true, ListOptions: listOpt}
	global.GVA_LOG.Info("Github scan with the query:", zap.Any("github", query))
	for {
		result, nextPage := searchCodeByOpt(c, ctx, query, *opt)
		if result != nil {
			allSearchResult = append(allSearchResult, result)
		}
		if nextPage <= 0 {
			break
		}
		opt.Page = nextPage
	}
	return allSearchResult, err
}

func BuildQuery(query string) (string, error) {
	query = query + " in:file"
	err, extensionFilters := model.GetFilterByClass("extension")
	if err != nil {
		return query, err
	}
	// if there is no record, does not return err
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query, nil
	}
	str := ""
	for _, extensionFilter := range extensionFilters {
		extensions := strings.Split(extensionFilter.Content, ",")
		filterType := extensionFilter.FilterType
		for _, extension := range extensions {
			if filterType == "blacklist" {
				str += " -extension:" + extension
			} else {
				str += " +extension:" + extension
			}
		}
	}

	err, keywordFilters := model.GetFilterByClass("keyword")

	for _, keywordFilter := range keywordFilters {
		keywords := strings.Split(keywordFilter.Content, ",")
		filterType := keywordFilter.FilterType
		for _, keyword := range keywords {
			if filterType == "black" {
				str += " NOT " + keyword
			} else {
				str += " " + keyword
			}
		}
	}

	builtQuery := query + str
	return builtQuery, err
}

func searchCodeByOpt(c *Client, ctx context.Context, query string, opt github.SearchOptions) (*github.CodeSearchResult,
	int) {
	result, res, err := c.Client.Search.Code(ctx, query, &opt)
	if _, ok := err.(*github.RateLimitError); ok {
		global.GVA_LOG.Warn("Trigger the github rate limit, ready to sleep 5 minutes")
		time.Sleep(5 * time.Minute)
	}
	if err == nil {
		// https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits
		resetTimeStamp := res.Rate.Reset
		time.Sleep(resetTimeStamp.Sub(time.Now()))
		time.Sleep(5 * time.Second)
	}

	if res != nil && res.Rate.Remaining < 10 {
		time.Sleep(45 * time.Second)
	}

	if err == nil {
		global.GVA_LOG.Info("Search for "+query, zap.Any("remaining", res.Rate.Remaining), zap.Any("nextPage",
			res.NextPage), zap.Any("lastPage", res.LastPage))
	} else {
		global.GVA_LOG.Error("Search error", zap.Any("github search error", err))
		time.Sleep(30 * time.Second)
		return nil, 0
	}
	return result, res.NextPage
}

package githubsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v57/github"
	"github.com/gookit/color"
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
	color.Debug.Print(fmt.Sprintf("Github fetch %d rules, begin the scan task\n", len(rules)))
	Search(rules)
	color.Debug.Print(fmt.Sprintf("Comple the scan of Github, start to sleep %d seconds", duration))
	time.Sleep(duration * time.Second)
}

func (c *Client) GetCommiter(ctx context.Context, owner, repo string) string {
	commit, _, err := c.Client.Repositories.GetCommit(ctx, owner, repo, "master", nil)
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
		result, nextPage := c.searchCodeByOpt(ctx, query, *opt)
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

func (c *Client) searchCodeByOpt(ctx context.Context, query string, opt github.SearchOptions) (*github.CodeSearchResult,
	int) {
	result, res, err := c.Client.Search.Code(ctx, query, &opt)

	// Handle error first
	if err != nil {
		var rateLimitError *github.RateLimitError
		if errors.As(err, &rateLimitError) {
			global.GVA_LOG.Warn("Trigger the github rate limit")
			if res != nil {
				resetTimeStamp := res.Rate.Reset
				sleepDuration := resetTimeStamp.Sub(time.Now()) + 10*time.Second
				global.GVA_LOG.Warn(fmt.Sprintf("Ready to sleep for %v", sleepDuration))
				time.Sleep(sleepDuration)
			}
		} else {
			global.GVA_LOG.Error("Search error", zap.Any("github search error", err))
			time.Sleep(30 * time.Second)
			return nil, 0
		}
	}

	// Check if response is nil before accessing its fields
	if res == nil {
		global.GVA_LOG.Error("Received nil response from GitHub API")
		return nil, 0
	}

	// Now safe to access res.Rate
	if res.Rate.Remaining < 3 {
		color.Info.Print("the remaining is less than 3, switch to another token\n")
		newGithubClient, newToken := c.NextClient()
		c.Client = newGithubClient
		c.Token = newToken
	}

	global.GVA_LOG.Info("Search for "+query, zap.Any("remaining", res.Rate.Remaining), zap.Any("nextPage",
		res.NextPage), zap.Any("lastPage", res.LastPage))

	return result, res.NextPage
}

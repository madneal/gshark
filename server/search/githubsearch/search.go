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

	"sync"
	"time"
)

const SearchNum = 10

func GenerateSearchCodeTask() (map[int][]model.Rule, error) {
	result := make(map[int][]model.Rule)
	// get rules with the type of github
	err, rules := service.GetValidRulesByType("github")
	if len(rules) == 0 {
		global.GVA_LOG.Info("Rules of github is empty, please specify one rule for scan at least")
	}
	ruleNum := len(rules)
	batch := ruleNum / SearchNum

	for i := 0; i < batch; i++ {
		result[i] = rules[SearchNum*i : SearchNum*(i+1)]
	}

	if ruleNum%SearchNum != 0 {
		result[batch] = rules[SearchNum*batch : ruleNum]
	}
	return result, err
}

func Search(rules []model.Rule) {
	var wg sync.WaitGroup
	wg.Add(len(rules))
	client, token, err := GetGithubClient()
	var content string
	var counts int
	if err == nil && token != "" {
		for _, rule := range rules {
			go func(rule model.Rule) {
				defer wg.Done()
				results, err := client.SearchCode(rule.Content)
				if err != nil {
					return
				}
				counts = SaveResult(results, &rule.Content)
				if counts > 0 {
					content += fmt.Sprintf("%s: %d条<br>", rule.Content, counts)
				}
			}(rule)
		}
		wg.Wait()
	}
	if counts > 0 {
		err = utils.EmailSend("Github敏感信息报告", content)
		if err != nil {
			global.GVA_LOG.Error("send email error", zap.Any("err", err))
		}
	}
}

func RunSearchTask(mapRules map[int][]model.Rule, err error) {
	if err == nil {
		for _, rules := range mapRules {
			startTime := time.Now()
			Search(rules)
			usedTime := time.Since(startTime).Seconds()
			if usedTime < 60 {
				time.Sleep(time.Duration(60 - usedTime))
			}
		}
	}
}

func SaveResult(results []*github.CodeSearchResult, keyword *string) int {
	searchResults := ConvertToSearchResults(results, keyword)
	insertCount := 0
	for _, result := range searchResults {
		err, exist := service.CheckExistOfSearchResult(&result)
		if exist {
			continue
		}
		err = service.CreateSearchResult(result)
		if err != nil {
			global.GVA_LOG.Error("save search result error", zap.Any("save searchResult error",
				err))
		} else {
			insertCount++
		}
	}
	return insertCount
}

func ConvertToSearchResults(results []*github.CodeSearchResult, keyword *string) []model.SearchResult {
	searchResults := make([]model.SearchResult, 0)
	for _, result := range results {
		codeResults := result.CodeResults
		for _, codeResult := range codeResults {
			searchResult := model.SearchResult{
				RepoUrl: *codeResult.Repository.HTMLURL,
				Repo:    *codeResult.Repository.FullName,
				Keyword: *keyword,
				Url:     *codeResult.HTMLURL,
				Path:    *codeResult.Path,
				Status:  0,
			}
			if len(codeResult.TextMatches) > 0 {
				hash := utils.GenMd5WithSpecificLen(*(codeResult.TextMatches[0].Fragment), 50)
				searchResult.TextmatchMd5 = hash
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
	RunSearchTask(GenerateSearchCodeTask())
	global.GVA_LOG.Info(fmt.Sprintf("Comple the scan of Github, start to sleep %d seconds", duration))
	time.Sleep(duration * time.Second)
}

func (c *Client) SearchCode(keyword string) ([]*github.CodeSearchResult, error) {
	var allSearchResult []*github.CodeSearchResult
	var err error
	ctx := context.Background()
	listOpt := github.ListOptions{PerPage: 100}
	opt := &github.SearchOptions{Sort: "indexed", Order: "desc", TextMatch: true, ListOptions: listOpt}
	query := keyword + " in:file"
	//query, err = BuildQuery(query)
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
	err, filterRule := model.GetFilterRule()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query, err
	}
	str := ""
	if filterRule.Extension != "" {
		extensions := strings.Split(filterRule.Extension, ",")
		for _, extension := range extensions {
			str += " -extension:" + extension
		}
	}
	if filterRule.WhiteExts != "" {
		whiteExts := strings.Split(filterRule.WhiteExts, ",")
		for _, whiteExt := range whiteExts {
			str += " +extension:" + whiteExt
		}
	}
	if filterRule.Keywords != "" {
		keyswords := strings.Split(filterRule.Keywords, ",")
		for _, keyword := range keyswords {
			str += " NOT " + keyword
		}
	}
	builtQuery := query + str
	return builtQuery, err
}

func searchCodeByOpt(c *Client, ctx context.Context, query string, opt github.SearchOptions) (*github.CodeSearchResult,
	int) {
	query, err := BuildQuery(query)
	result, res, err := c.Client.Search.Code(ctx, query, &opt)

	// for best guidelines, wait one second
	// https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits

	time.Sleep(1 * time.Second)
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

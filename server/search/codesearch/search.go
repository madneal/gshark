package codesearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

const (
	searchcodeRequestTimeout = 30 * time.Second
	maxSearchcodePages       = 50 // defensive cap; normal results end via hasResult=false
)

var (
	searchcodeHTTPClient = &http.Client{Timeout: searchcodeRequestTimeout}
	searchcodeBaseURL    = "https://searchcode.com/api/codesearch_I/"
)

func RunTask(duration time.Duration) {
	err, rules := service.GetValidRulesByType("searchcode")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType searchcode err", zap.Error(err))
		return
	}
	if len(rules) == 0 {
		global.GVA_LOG.Info("Rules of search code is empty")
		return
	}
	for _, rule := range rules {
		global.GVA_LOG.Info(fmt.Sprintf("Search for %s in searchcode", rule.Content))
		codeResults := SearchForSearchCode(rule, searchcodeHTTPClient)
		SaveResults(codeResults, &rule.Content)
	}
	global.GVA_LOG.Info("Complete the scan of searchcode")
	time.Sleep(duration * time.Second)
}

func SaveResults(results []*model.SearchResult, keyword *string) {
	if len(results) == 0 {
		return
	}
	stats := service.SaveSearchResultPointersWithStats(results, *keyword)
	global.GVA_LOG.Info(stats.Summary(*keyword, "SearchCode"))
}

func SearchForSearchCode(rule model.Rule, client *http.Client) []*model.SearchResult {
	keyword := rule.Content
	totalCodeResults := make([]*model.SearchResult, 0)
	for page := 0; page < maxSearchcodePages; page++ {
		url := searchcodeBaseURL + "?q=" + keyword + "&p=" + strconv.Itoa(page)
		global.GVA_LOG.Info("search searchcode result for page " + strconv.Itoa(page))
		codeResults, hasResult := GetResult(client, url)
		totalCodeResults = append(totalCodeResults, codeResults...)
		if !hasResult {
			break
		}
	}
	return totalCodeResults
}

func GetResult(client *http.Client, url string) ([]*model.SearchResult, bool) {
	codeResults := make([]*model.SearchResult, 0)

	resp, err := client.Get(url)
	if err != nil || resp == nil {
		global.GVA_LOG.Error("search result of search code error", zap.Any("err", err))
		return codeResults, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		global.GVA_LOG.Error(fmt.Sprintf("Request to %s error, status code: %d", url, resp.StatusCode))
		return codeResults, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Error("read searchcode response body err", zap.Error(err))
		return codeResults, false
	}

	var result model.SearchCodeRes
	if jErr := json.Unmarshal(body, &result); jErr != nil {
		global.GVA_LOG.Error("json unmarshal searchCodeRes err", zap.Error(jErr))
		return codeResults, false
	}

	results := result.Results
	if len(results) == 0 {
		return codeResults, false
	}
	for _, val := range results {
		if strings.Contains(val.Repo, "github") {
			continue
		}
		var lines string
		for _, line := range val.Lines {
			lines += fmt.Sprint(line) + "\n"
		}
		repoPath := val.Repo
		textMatch := new(model.TextMatch)
		textMatch.Fragment = &lines
		textMatchs := make([]model.TextMatch, 0)
		textMatchs = append(textMatchs, *textMatch)
		b, _ := json.Marshal(textMatchs)
		codeResult := model.SearchResult{
			Path:            val.Filename,
			RepoUrl:         val.Location,
			Status:          0,
			Url:             val.Url,
			Repo:            repoPath,
			TextMatchesJson: b,
		}
		codeResults = append(codeResults, &codeResult)
	}
	return codeResults, true
}

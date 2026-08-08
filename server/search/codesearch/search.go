package codesearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
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
	searchcodeHTTPClient  = &http.Client{Timeout: searchcodeRequestTimeout}
	searchcodeBaseURL     = "https://searchcode.com/api/codesearch_I/"
	searchcodeUnavailable atomic.Bool
)

// ErrSearchcodeUnavailable indicates that the legacy global Searchcode API is
// no longer available. The current Searchcode service uses a repository-scoped
// API and is not a drop-in replacement for this provider's global queries.
var ErrSearchcodeUnavailable = errors.New("Searchcode legacy API unavailable")

func RunTask() model.ScanOutcome {
	if searchcodeUnavailable.Load() {
		message := "Searchcode provider skipped: legacy global-search API is unavailable"
		global.GVA_LOG.Warn(message)
		return model.ScanSkipped(message)
	}

	err, rules := service.GetValidRulesByType("searchcode")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType searchcode err", zap.Error(err))
		return model.ScanFailed("Failed to load Searchcode rules: " + err.Error())
	}
	if len(rules) == 0 {
		message := "No enabled Searchcode rules; provider skipped"
		global.GVA_LOG.Info(message)
		return model.ScanSkipped(message)
	}
	var scanErrors []error
	for _, rule := range rules {
		global.GVA_LOG.Info(fmt.Sprintf("Search for %s in searchcode", rule.Content))
		codeResults, err := SearchForSearchCode(rule, searchcodeHTTPClient)
		if err != nil {
			if errors.Is(err, ErrSearchcodeUnavailable) {
				searchcodeUnavailable.Store(true)
				message := "Searchcode provider skipped: legacy global-search API returned 404; migrate the provider or disable Searchcode rules"
				global.GVA_LOG.Warn(message, zap.Error(err))
				return model.ScanSkipped(message)
			}
			scanErrors = append(scanErrors, fmt.Errorf("search %q: %w", rule.Content, err))
			continue
		}
		SaveResults(codeResults, &rule.Content)
	}
	if err := errors.Join(scanErrors...); err != nil {
		return model.ScanFailed("Searchcode scan completed with errors: " + err.Error())
	}
	global.GVA_LOG.Info("Complete the scan of searchcode")
	return model.ScanSuccess(fmt.Sprintf("Completed %d Searchcode rules", len(rules)))
}

func SaveResults(results []*model.SearchResult, keyword *string) {
	if len(results) == 0 {
		return
	}
	stats := service.SaveSearchResultPointersWithStats(results, *keyword)
	global.GVA_LOG.Info(stats.Summary(*keyword, "SearchCode"))
}

func SearchForSearchCode(rule model.Rule, client *http.Client) ([]*model.SearchResult, error) {
	keyword := rule.Content
	totalCodeResults := make([]*model.SearchResult, 0)
	for page := 0; page < maxSearchcodePages; page++ {
		url := searchcodeBaseURL + "?q=" + keyword + "&p=" + strconv.Itoa(page)
		global.GVA_LOG.Info("search searchcode result for page " + strconv.Itoa(page))
		codeResults, hasResult, err := GetResult(client, url)
		if err != nil {
			return totalCodeResults, err
		}
		totalCodeResults = append(totalCodeResults, codeResults...)
		if !hasResult {
			break
		}
	}
	return totalCodeResults, nil
}

func GetResult(client *http.Client, url string) ([]*model.SearchResult, bool, error) {
	codeResults := make([]*model.SearchResult, 0)

	resp, err := client.Get(url)
	if err != nil || resp == nil {
		global.GVA_LOG.Error("search result of search code error", zap.Any("err", err))
		if err == nil {
			err = errors.New("Searchcode returned a nil response")
		}
		return codeResults, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return codeResults, false, fmt.Errorf("%w: request to %s returned status 404", ErrSearchcodeUnavailable, url)
		}
		err = fmt.Errorf("request to %s returned status %d", url, resp.StatusCode)
		global.GVA_LOG.Error(err.Error())
		return codeResults, false, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Error("read searchcode response body err", zap.Error(err))
		return codeResults, false, err
	}

	var result model.SearchCodeRes
	if jErr := json.Unmarshal(body, &result); jErr != nil {
		global.GVA_LOG.Error("json unmarshal searchCodeRes err", zap.Error(jErr))
		return codeResults, false, jErr
	}

	results := result.Results
	if len(results) == 0 {
		return codeResults, false, nil
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
	return codeResults, true, nil
}

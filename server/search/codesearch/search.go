package codesearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

const (
	searchcodeRequestTimeout = 30 * time.Second
	searchcodeMaxResults     = 100
	searchcodeMaxPages       = 50
	searchcodeClientName     = "gshark"
)

var (
	searchcodeHTTPClient = &http.Client{Timeout: searchcodeRequestTimeout}
	searchcodeBaseURL    = "https://api.searchcode.com/api/v1/code_search"
)

type codeSearchRequest struct {
	Repository        string `json:"repository"`
	Query             string `json:"query"`
	Offset            int    `json:"offset"`
	MaxResults        int    `json:"max_results"`
	MaxMatchesPerFile int    `json:"max_matches_per_file"`
	ContextLines      int    `json:"context_lines"`
	SnippetMode       string `json:"snippet_mode"`
}

type codeSearchResponse struct {
	Repository      string             `json:"repository"`
	CommitSHA       string             `json:"commit_sha"`
	TotalMatches    int                `json:"total_matches"`
	ResultsReturned int                `json:"results_returned"`
	HasMore         bool               `json:"has_more"`
	Truncated       bool               `json:"truncated"`
	Results         []codeSearchResult `json:"results"`
}

type codeSearchResult struct {
	File          string            `json:"file"`
	Language      string            `json:"language"`
	MatchesInFile int               `json:"matches_in_file"`
	Matches       []codeSearchMatch `json:"matches"`
}

type codeSearchMatch struct {
	Line          int      `json:"line"`
	Content       string   `json:"content"`
	ContextBefore []string `json:"context_before"`
	ContextAfter  []string `json:"context_after"`
}

type searchcodeAPIError struct {
	Error struct {
		Code              string `json:"code"`
		Message           string `json:"message"`
		RetryAfterSeconds int    `json:"retry_after_seconds"`
	} `json:"error"`
}

func RunTask() model.ScanOutcome {
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

	repositories, err := configuredRepositories()
	if err != nil {
		global.GVA_LOG.Error("Get Searchcode repositories err", zap.Error(err))
		return model.ScanFailed("Failed to load Searchcode repositories: " + err.Error())
	}
	if len(repositories) == 0 {
		message := "No Searchcode repositories configured; add public repository URLs to search.searchcode-repositories or repo records with type=searchcode"
		global.GVA_LOG.Warn(message)
		return model.ScanSkipped(message)
	}

	var scanErrors []error
	for _, repository := range repositories {
		for _, rule := range rules {
			global.GVA_LOG.Info("Search in Searchcode repository", zap.String("repository", repository), zap.String("query", rule.Content))
			codeResults, err := SearchForSearchCode(rule, repository, searchcodeHTTPClient)
			if err != nil {
				scanErrors = append(scanErrors, fmt.Errorf("search %q in %s: %w", rule.Content, repository, err))
				continue
			}
			SaveResults(codeResults, &rule.Content)
		}
	}
	if err := errors.Join(scanErrors...); err != nil {
		return model.ScanFailed("Searchcode scan completed with errors: " + err.Error())
	}
	return model.ScanSuccess(fmt.Sprintf("Completed %d Searchcode rules across %d repositories", len(rules), len(repositories)))
}

func configuredRepositories() ([]string, error) {
	repositories := make([]string, 0, len(global.GVA_CONFIG.Search.SearchcodeRepositories))
	for _, repository := range global.GVA_CONFIG.Search.SearchcodeRepositories {
		repositories = appendRepository(repositories, repository)
	}

	err, records := service.GetRepoByType("searchcode")
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		repositories = appendRepository(repositories, record.Url)
	}
	return repositories, nil
}

func appendRepository(repositories []string, repository string) []string {
	repository = strings.TrimRight(strings.TrimSpace(repository), "/")
	if repository == "" {
		return repositories
	}
	parsed, err := url.Parse(repository)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		global.GVA_LOG.Warn("Ignoring invalid Searchcode repository URL", zap.String("repository", repository))
		return repositories
	}
	for _, existing := range repositories {
		if strings.EqualFold(existing, repository) {
			return repositories
		}
	}
	return append(repositories, repository)
}

func SaveResults(results []*model.SearchResult, keyword *string) {
	if len(results) == 0 {
		return
	}
	stats := service.SaveSearchResultPointersWithStats(results, *keyword)
	global.GVA_LOG.Info(stats.Summary(*keyword, "SearchCode"))
}

func SearchForSearchCode(rule model.Rule, repository string, client *http.Client) ([]*model.SearchResult, error) {
	allResults := make([]*model.SearchResult, 0)
	request := codeSearchRequest{
		Repository:        repository,
		Query:             rule.Content,
		MaxResults:        searchcodeMaxResults,
		MaxMatchesPerFile: 50,
		ContextLines:      2,
		SnippetMode:       "grep",
	}

	for page := 0; page < searchcodeMaxPages; page++ {
		response, err := GetResult(client, request)
		if err != nil {
			return allResults, err
		}
		allResults = append(allResults, convertResults(*response, request.Repository, request.Query)...)
		if !response.HasMore || len(response.Results) == 0 {
			break
		}
		nextOffset := request.Offset + len(response.Results)
		if nextOffset <= request.Offset {
			break
		}
		request.Offset = nextOffset
	}
	return allResults, nil
}

func GetResult(client *http.Client, request codeSearchRequest) (*codeSearchResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal Searchcode request: %w", err)
	}
	endpoint := searchcodeBaseURL + "?client=" + url.QueryEscape(searchcodeClientName)
	httpRequest, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create Searchcode request: %w", err)
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := client.Do(httpRequest)
	if err != nil || response == nil {
		if err == nil {
			err = errors.New("Searchcode returned a nil response")
		}
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read Searchcode response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		var apiError searchcodeAPIError
		if json.Unmarshal(responseBody, &apiError) == nil && apiError.Error.Message != "" {
			return nil, fmt.Errorf("Searchcode request returned status %d (%s): %s", response.StatusCode, apiError.Error.Code, apiError.Error.Message)
		}
		return nil, fmt.Errorf("Searchcode request returned status %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result codeSearchResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("decode Searchcode response: %w", err)
	}
	return &result, nil
}

func convertResults(response codeSearchResponse, repository, keyword string) []*model.SearchResult {
	repositoryURL := strings.TrimRight(repository, "/")
	if response.Repository != "" {
		repositoryURL = strings.TrimRight(response.Repository, "/")
	}
	repoName := repositoryName(repositoryURL)
	results := make([]*model.SearchResult, 0, len(response.Results))
	for _, result := range response.Results {
		matches := make([]model.TextMatch, 0, len(result.Matches))
		for _, match := range result.Matches {
			fragmentLines := make([]string, 0, len(match.ContextBefore)+1+len(match.ContextAfter))
			fragmentLines = append(fragmentLines, match.ContextBefore...)
			fragmentLines = append(fragmentLines, match.Content)
			fragmentLines = append(fragmentLines, match.ContextAfter...)
			fragment := strings.Join(fragmentLines, "\n")
			matches = append(matches, model.TextMatch{Fragment: &fragment})
		}
		textMatches, err := json.Marshal(matches)
		if err != nil {
			global.GVA_LOG.Error("marshal Searchcode matches error", zap.Error(err))
			continue
		}
		results = append(results, &model.SearchResult{
			Repo:            repoName,
			RepoUrl:         repositoryURL,
			Keyword:         keyword,
			Path:            result.File,
			Url:             resultURL(repositoryURL, response.CommitSHA, result.File),
			Status:          0,
			TextMatchesJson: textMatches,
		})
	}
	return results
}

func repositoryName(repository string) string {
	parsed, err := url.Parse(repository)
	if err != nil {
		return repository
	}
	name := strings.Trim(parsed.Path, "/")
	return strings.TrimSuffix(name, ".git")
}

func resultURL(repository, commit, file string) string {
	file = strings.TrimLeft(file, "/")
	ref := commit
	if ref == "" {
		ref = "HEAD"
	}
	return strings.TrimRight(repository, "/") + "/blob/" + url.PathEscape(ref) + "/" + file
}

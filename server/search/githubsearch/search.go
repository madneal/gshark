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

func Search(rules []model.Rule) error {
	if len(rules) == 0 {
		return nil
	}
	client, err := GetGithubClient()
	if err != nil {
		global.GVA_LOG.Error("GetGithubClient err", zap.Error(err))
		return err
	}
	return searchCode(client, rules)
}

func searchCode(client *Client, rules []model.Rule) error {
	if len(rules) == 0 {
		return nil
	}
	var content string
	var scanErrors []error
	for _, rule := range rules {
		query, err := BuildQuery(rule.Content)
		if err != nil {
			global.GVA_LOG.Error("BuildQuery error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("build query for %q: %w", rule.Content, err))
			continue
		}
		var ruleInserted int
		var ruleRepos []string
		var ruleHasMoreRepos bool
		err = client.SearchCodeStream(query, func(page []*github.CodeSearchResult) error {
			stats := SaveResultWithStats(page, rule.Content)
			ruleInserted += stats.Inserted
			var more bool
			ruleRepos, more = appendUniqueRepos(ruleRepos, stats.Repos)
			ruleHasMoreRepos = ruleHasMoreRepos || more
			return nil
		})
		if err != nil {
			global.GVA_LOG.Error("SearchCode error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("search %q: %w", rule.Content, err))
			continue
		}
		global.GVA_LOG.Info(fmt.Sprintf("[GitHub] keyword=%q: inserted=%d", rule.Content, ruleInserted))
		if ruleInserted > 0 {
			repoInfo := ""
			if len(ruleRepos) > 0 {
				repoInfo = fmt.Sprintf(" (repos: %s)", strings.Join(ruleRepos, ", "))
				if ruleHasMoreRepos {
					repoInfo += " +more"
				}
			}
			content += fmt.Sprintf("%s: %d条%s<br>", rule.Content, ruleInserted, repoInfo)
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
	notifyNewResults("Github敏感信息报告", content)
	return errors.Join(scanErrors...)
}

func formatInsertedSummary(keyword string, stats *service.SaveResultStats) string {
	if stats == nil || stats.Inserted == 0 {
		return ""
	}
	repoInfo := ""
	if len(stats.Repos) > 0 {
		if len(stats.Repos) <= 3 {
			repoInfo = fmt.Sprintf(" (repos: %s)", strings.Join(stats.Repos, ", "))
		} else {
			repoInfo = fmt.Sprintf(" (repos: %s +%d more)", strings.Join(stats.Repos[:3], ", "), len(stats.Repos)-3)
		}
	}
	return fmt.Sprintf("%s: %d条%s<br>", keyword, stats.Inserted, repoInfo)
}

func notifyNewResults(title, content string) {
	if content == "" {
		return
	}
	if global.GVA_CONFIG.Email.Enable {
		if err := utils.EmailSend(title, content); err != nil {
			global.GVA_LOG.Error("send email error", zap.Any("err", err))
		}
	}
	if global.GVA_CONFIG.Wechat.Enable {
		if err := utils.BotSend(title + "\n" + content); err != nil {
			global.GVA_LOG.Error("send wechat error", zap.Any("err", err))
		}
	}
}

func SaveResultWithStats(results []*github.CodeSearchResult, keyword string) *service.SaveResultStats {
	searchResults := ConvertToSearchResults(results, keyword)
	return service.SaveSearchResultsWithStats(searchResults)
}

func appendUniqueRepos(existing []string, repos []string) ([]string, bool) {
	const previewLimit = 3
	hasMore := false
	for _, repo := range repos {
		duplicate := false
		for _, current := range existing {
			if current == repo {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}
		if len(existing) < previewLimit {
			existing = append(existing, repo)
		} else {
			hasMore = true
		}
	}
	return existing, hasMore
}

func ConvertToSearchResults(results []*github.CodeSearchResult, keyword string) []model.SearchResult {
	searchResults := make([]model.SearchResult, 0)
	for _, result := range results {
		codeResults := result.CodeResults
		for _, codeResult := range codeResults {
			searchResult := model.SearchResult{
				RepoUrl: *codeResult.Repository.HTMLURL,
				Repo:    *codeResult.Repository.FullName,
				Keyword: keyword,
				Url:     *codeResult.HTMLURL,
				Path:    *codeResult.Path,
				Status:  0,
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

func RunTask() model.ScanOutcome {
	err, codeRules := getGithubRules("github")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType github err", zap.Error(err))
		return model.ScanFailed("Failed to load GitHub rules: " + err.Error())
	}
	err, issueRules := getGithubRules("github_issue")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType github_issue err", zap.Error(err))
		return model.ScanFailed("Failed to load GitHub issue rules: " + err.Error())
	}
	err, gistRules := getGithubRules("gist")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType gist err", zap.Error(err))
		return model.ScanFailed("Failed to load Gist rules: " + err.Error())
	}

	color.Debug.Print(fmt.Sprintf("Github fetch %d code / %d issue / %d gist rules\n",
		len(codeRules), len(issueRules), len(gistRules)))
	if len(codeRules) == 0 && len(issueRules) == 0 && len(gistRules) == 0 {
		message := "No enabled GitHub, issue, or gist rules; provider skipped"
		global.GVA_LOG.Info(message)
		return model.ScanSkipped(message)
	}

	client, err := GetGithubClient()
	if err != nil {
		global.GVA_LOG.Error("GetGithubClient err", zap.Error(err))
		return model.ScanFailed("GitHub client initialization failed: " + err.Error())
	}

	var scanErrors []error
	if err := searchCode(client, codeRules); err != nil {
		scanErrors = append(scanErrors, err)
	}
	if err := searchIssues(client, issueRules); err != nil {
		scanErrors = append(scanErrors, err)
	}
	if err := searchGists(client, gistRules); err != nil {
		scanErrors = append(scanErrors, err)
	}
	if err := errors.Join(scanErrors...); err != nil {
		return model.ScanFailed("GitHub scan completed with errors: " + err.Error())
	}
	color.Debug.Print("Complete the scan of Github\n")
	return model.ScanSuccess(fmt.Sprintf("Completed %d code, %d issue, %d gist rules",
		len(codeRules), len(issueRules), len(gistRules)))
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
	err := c.SearchCodeStream(query, func(page []*github.CodeSearchResult) error {
		allSearchResult = append(allSearchResult, page...)
		return nil
	})
	return allSearchResult, err
}

// SearchCodeStream fetches GitHub search pages one at a time. The callback is
// invoked before the next page is requested so callers can persist each page
// and keep the result set out of memory.
func (c *Client) SearchCodeStream(query string, onPage func([]*github.CodeSearchResult) error) error {
	ctx := context.Background()
	listOpt := github.ListOptions{PerPage: 100}
	opt := &github.SearchOptions{TextMatch: true, ListOptions: listOpt}
	global.GVA_LOG.Info("Github scan with the query:", zap.Any("github", query))
	for {
		result, nextPage := c.searchCodeByOpt(ctx, query, *opt)
		if result == nil {
			return fmt.Errorf("GitHub returned no response for query %q", query)
		}
		if err := onPage([]*github.CodeSearchResult{result}); err != nil {
			return err
		}
		if nextPage <= 0 {
			break
		}
		opt.Page = nextPage
	}
	return nil
}

func BuildQuery(query string) (string, error) {
	query = query + " in:file"
	ext, err := appendExtensionFilters("", " -extension:", " +extension:")
	if err != nil {
		return query + ext, err
	}
	kw, err := appendKeywordFilters("", " NOT ", " ")
	return query + ext + kw, err
}

func BuildIssueQuery(query string) (string, error) {
	if !strings.Contains(query, "in:") {
		query += " in:title,body,comments"
	}
	return appendKeywordFilters(query, " NOT ", " ")
}

func BuildGistQuery(query string) (string, error) {
	ext, err := appendExtensionFilters("", " -extension:", " extension:")
	if err != nil {
		return query + ext, err
	}
	kw, err := appendKeywordFilters("", " -", " ")
	return query + ext + kw, err
}

func appendExtensionFilters(str, denyPrefix, allowPrefix string) (string, error) {
	err, extensionFilters := getFiltersByClass("extension")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return str, nil
		}
		return str, err
	}
	for _, extensionFilter := range extensionFilters {
		extensions := strings.Split(extensionFilter.Content, ",")
		filterType := extensionFilter.FilterType
		for _, extension := range extensions {
			extension = strings.TrimSpace(extension)
			if extension == "" {
				continue
			}
			if filterType == "blacklist" {
				str += denyPrefix + extension
			} else {
				str += allowPrefix + extension
			}
		}
	}
	return str, nil
}

func appendKeywordFilters(str, denyPrefix, allowPrefix string) (string, error) {
	err, keywordFilters := getFiltersByClass("keyword")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return str, err
	}
	for _, keywordFilter := range keywordFilters {
		keywords := strings.Split(keywordFilter.Content, ",")
		filterType := keywordFilter.FilterType
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword == "" {
				continue
			}
			if filterType == "black" || filterType == "blacklist" {
				str += denyPrefix + keyword
			} else {
				str += allowPrefix + keyword
			}
		}
	}
	return str, nil
}

// maxSearchAttempts bounds how many times a single page is retried after
// hitting a rate limit, so a persistently exhausted token set eventually
// gives up instead of looping forever.
const maxSearchAttempts = 3

// sleepFn is overridden in tests so rate-limit backoffs don't actually block.
var sleepFn = time.Sleep

var getGithubRules = service.GetValidRulesByType

var getFiltersByClass = model.GetFilterByClass

func (c *Client) searchCodeByOpt(ctx context.Context, query string, opt github.SearchOptions) (*github.CodeSearchResult,
	int) {
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		result, res, err := c.Client.Search.Code(ctx, query, &opt)
		if err == nil {
			return result, c.noteSearchResponse(query, res)
		}
		if !c.recoverSearchError(query, err, attempt) {
			return nil, 0
		}
	}

	global.GVA_LOG.Error("exhausted retry attempts for Github search", zap.String("query", query))
	return nil, 0
}

func (c *Client) recoverSearchError(query string, err error, attempt int) bool {
	var rateLimitError *github.RateLimitError
	var abuseRateLimitError *github.AbuseRateLimitError
	switch {
	case errors.As(err, &rateLimitError):
		global.GVA_LOG.Warn("Trigger the github rate limit", zap.Int("attempt", attempt), zap.String("query", query))
		if !c.attemptRotate() {
			sleepUntil(rateLimitError.Rate.Reset.Time)
		}
		return true
	case errors.As(err, &abuseRateLimitError):
		global.GVA_LOG.Warn("Trigger the github secondary rate limit", zap.Int("attempt", attempt), zap.String("query", query))
		if !c.attemptRotate() {
			sleepForAbuse(abuseRateLimitError.RetryAfter)
		}
		return true
	default:
		global.GVA_LOG.Error("Search error", zap.Any("github search error", err), zap.String("query", query))
		sleepFn(30 * time.Second)
		return false
	}
}

func (c *Client) noteSearchResponse(query string, res *github.Response) int {
	if res == nil {
		global.GVA_LOG.Error("Received nil response from GitHub API")
		return 0
	}

	if res.Rate.Remaining < 3 {
		color.Info.Print("the remaining is less than 3, switch to another token\n")
		c.attemptRotate()
	}

	global.GVA_LOG.Info("Search for "+query, zap.Any("remaining", res.Rate.Remaining), zap.Any("nextPage",
		res.NextPage), zap.Any("lastPage", res.LastPage))

	return res.NextPage
}

// attemptRotate switches to the next configured github token and reports
// whether it actually changed to a different token. Tests inject c.rotate
// to avoid hitting the database; production code falls back to rotateToken.
func (c *Client) attemptRotate() bool {
	if c.rotate != nil {
		return c.rotate()
	}
	return c.rotateToken()
}

// rotateToken switches to the next configured github token and reports
// whether it actually changed to a different token. When only one token is
// configured, NextClient returns the same token, so callers know to fall
// back to sleeping instead of busy-retrying with the same exhausted token.
func (c *Client) rotateToken() bool {
	previousToken := c.Token
	newGithubClient, newToken := c.NextClient()
	if newGithubClient == nil || newToken == "" || newToken == previousToken {
		return false
	}
	c.Client = newGithubClient
	c.Token = newToken
	return true
}

func sleepUntil(reset time.Time) {
	sleepDuration := time.Until(reset) + 10*time.Second
	if sleepDuration < 0 {
		sleepDuration = 10 * time.Second
	}
	global.GVA_LOG.Warn(fmt.Sprintf("Ready to sleep for %v", sleepDuration))
	sleepFn(sleepDuration)
}

func sleepForAbuse(retryAfter *time.Duration) {
	sleepDuration := 60 * time.Second
	if retryAfter != nil {
		sleepDuration = *retryAfter + 5*time.Second
	}
	global.GVA_LOG.Warn(fmt.Sprintf("Ready to sleep for %v due to secondary rate limit", sleepDuration))
	sleepFn(sleepDuration)
}

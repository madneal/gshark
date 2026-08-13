package githubsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

func searchIssues(client *Client, rules []model.Rule) error {
	if len(rules) == 0 {
		return nil
	}
	var content string
	var scanErrors []error
	for _, rule := range rules {
		query, err := BuildIssueQuery(rule.Content)
		if err != nil {
			global.GVA_LOG.Error("BuildIssueQuery error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("build issue query for %q: %w", rule.Content, err))
			continue
		}
		results, err := client.SearchIssues(query)
		if err != nil {
			global.GVA_LOG.Error("SearchIssues error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("issue search %q: %w", rule.Content, err))
			continue
		}
		stats := service.SaveSearchResultsWithStats(ConvertIssuesToSearchResults(results, rule.Content))
		global.GVA_LOG.Info(stats.Summary(rule.Content, "GitHub Issue"))
		content += formatInsertedSummary(rule.Content, stats)
	}
	notifyNewResults("GitHub Issue/PR 敏感信息报告", content)
	return errors.Join(scanErrors...)
}

func (c *Client) SearchIssues(query string) ([]*github.IssuesSearchResult, error) {
	var all []*github.IssuesSearchResult
	ctx := context.Background()
	opt := &github.SearchOptions{TextMatch: true, ListOptions: github.ListOptions{PerPage: 100}}
	global.GVA_LOG.Info("GitHub issue scan with the query:", zap.String("query", query))
	for {
		result, nextPage := c.searchIssuesByOpt(ctx, query, *opt)
		if result == nil {
			return all, fmt.Errorf("GitHub returned no issue response for query %q", query)
		}
		all = append(all, result)
		if nextPage <= 0 {
			return all, nil
		}
		opt.Page = nextPage
	}
}

func (c *Client) searchIssuesByOpt(ctx context.Context, query string, opt github.SearchOptions) (*github.IssuesSearchResult, int) {
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		result, res, err := c.Client.Search.Issues(ctx, query, &opt)
		if err == nil {
			return result, c.noteSearchResponse(query, res)
		}
		if !c.recoverSearchError(query, err, attempt) {
			return nil, 0
		}
	}
	global.GVA_LOG.Error("exhausted retry attempts for GitHub issue search", zap.String("query", query))
	return nil, 0
}

func ConvertIssuesToSearchResults(results []*github.IssuesSearchResult, keyword string) []model.SearchResult {
	searchResults := make([]model.SearchResult, 0)
	for _, result := range results {
		if result == nil {
			continue
		}
		for _, issue := range result.Issues {
			if issue == nil {
				continue
			}
			repoName, repoURL := repoFromIssue(issue)
			kind := "issues"
			if issue.IsPullRequest() {
				kind = "pull"
			}
			item := model.SearchResult{
				Repo:    repoName,
				RepoUrl: repoURL,
				Keyword: keyword,
				Path:    fmt.Sprintf("%s/%d", kind, issue.GetNumber()),
				Url:     issue.GetHTMLURL(),
				Status:  0,
				Matches: issueSnippet(issue),
			}
			if len(issue.TextMatches) > 0 {
				if encoded, err := json.Marshal(issue.TextMatches); err != nil {
					global.GVA_LOG.Error("json.marshal issue text matches error", zap.Error(err))
				} else {
					item.TextMatchesJson = encoded
				}
			}
			searchResults = append(searchResults, item)
		}
	}
	return searchResults
}

func repoFromIssue(issue *github.Issue) (name, htmlURL string) {
	if issue == nil {
		return "", ""
	}
	if repo := issue.GetRepository(); repo != nil && repo.GetFullName() != "" {
		name = repo.GetFullName()
		htmlURL = repo.GetHTMLURL()
		if htmlURL == "" && name != "" {
			htmlURL = "https://github.com/" + name
		}
		return name, htmlURL
	}
	name = strings.TrimPrefix(issue.GetRepositoryURL(), "https://api.github.com/repos/")
	if name == issue.GetRepositoryURL() {
		name = ""
	}
	if name != "" {
		htmlURL = "https://github.com/" + name
	}
	return name, htmlURL
}

func issueSnippet(issue *github.Issue) string {
	title := strings.TrimSpace(issue.GetTitle())
	body := strings.TrimSpace(issue.GetBody())
	if body == "" {
		return title
	}
	if len(body) > 400 {
		body = body[:400]
	}
	if title == "" {
		return body
	}
	return title + "\n" + body
}

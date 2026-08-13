package githubsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

const (
	maxGistSearchPages     = 5
	maxGistSearchPageBytes = 2 << 20
	gistSearchTimeout      = 30 * time.Second
	maxPublicGistFallback  = 50
)

var gistPathRe = regexp.MustCompile(`href="/(?P<owner>[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?)/(?P<id>[0-9a-f]{20,32})(?:["/#?])`)

var (
	fetchGistSearchPage = defaultFetchGistSearchPage
	gistSearchClient    = &http.Client{Timeout: gistSearchTimeout}
)

type gistHit struct {
	Owner string
	ID    string
}

func searchGists(client *Client, rules []model.Rule) error {
	if len(rules) == 0 {
		return nil
	}
	var content string
	var scanErrors []error
	for _, rule := range rules {
		query, err := BuildGistQuery(rule.Content)
		if err != nil {
			global.GVA_LOG.Error("BuildGistQuery error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("build gist query for %q: %w", rule.Content, err))
			continue
		}
		hits, err := client.SearchGistHits(query)
		if err != nil {
			global.GVA_LOG.Error("SearchGistHits error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("gist search %q: %w", rule.Content, err))
			continue
		}
		results, err := client.LoadGistResults(hits, rule.Content)
		if err != nil {
			global.GVA_LOG.Error("LoadGistResults error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("load gists for %q: %w", rule.Content, err))
		}
		stats := service.SaveSearchResultsWithStats(results)
		global.GVA_LOG.Info(stats.Summary(rule.Content, "Gist"))
		content += formatInsertedSummary(rule.Content, stats)
	}
	notifyNewResults("Gist 敏感信息报告", content)
	return errors.Join(scanErrors...)
}

func (c *Client) SearchGistHits(query string) ([]gistHit, error) {
	seen := make(map[string]struct{})
	hits := make([]gistHit, 0)
	for page := 1; page <= maxGistSearchPages; page++ {
		body, err := fetchGistSearchPage(c.Token, query, page)
		if err != nil {
			if page == 1 {
				global.GVA_LOG.Warn("gist HTML search failed, falling back to public gist stream", zap.Error(err))
				return c.listRecentPublicGists()
			}
			break
		}
		pageHits := parseGistSearchHTML(body)
		if len(pageHits) == 0 {
			if page == 1 {
				global.GVA_LOG.Warn("gist HTML search returned no parseable hits, falling back to public gist stream")
				return c.listRecentPublicGists()
			}
			break
		}
		added := 0
		for _, hit := range pageHits {
			if _, exists := seen[hit.ID]; exists {
				continue
			}
			seen[hit.ID] = struct{}{}
			hits = append(hits, hit)
			added++
		}
		if added == 0 {
			break
		}
	}
	return hits, nil
}

func (c *Client) listRecentPublicGists() ([]gistHit, error) {
	ctx := context.Background()
	opt := &github.GistListOptions{ListOptions: github.ListOptions{PerPage: 100}}
	hits := make([]gistHit, 0, maxPublicGistFallback)
	for {
		gists, res, err := c.getPublicGists(ctx, opt)
		if err != nil {
			return hits, err
		}
		for _, gist := range gists {
			if gist == nil || gist.GetID() == "" {
				continue
			}
			hits = append(hits, gistHit{Owner: gistOwnerLogin(gist), ID: gist.GetID()})
			if len(hits) >= maxPublicGistFallback {
				return hits, nil
			}
		}
		if res == nil || res.NextPage <= 0 {
			return hits, nil
		}
		opt.Page = res.NextPage
	}
}

func (c *Client) getPublicGists(ctx context.Context, opt *github.GistListOptions) ([]*github.Gist, *github.Response, error) {
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		gists, res, err := c.Client.Gists.ListAll(ctx, opt)
		if err == nil {
			c.noteSearchResponse("gists/public", res)
			return gists, res, nil
		}
		if !c.recoverSearchError("gists/public", err, attempt) {
			return nil, res, err
		}
	}
	return nil, nil, errors.New("exhausted retry attempts for public gist list")
}

func (c *Client) LoadGistResults(hits []gistHit, keyword string) ([]model.SearchResult, error) {
	results := make([]model.SearchResult, 0)
	var loadErrors []error
	ctx := context.Background()
	for _, hit := range hits {
		gist, err := c.getGist(ctx, hit.ID)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Errorf("gist %s: %w", hit.ID, err))
			continue
		}
		results = append(results, convertGistToSearchResults(gist, hit, keyword)...)
	}
	return results, errors.Join(loadErrors...)
}

func (c *Client) getGist(ctx context.Context, id string) (*github.Gist, error) {
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		gist, res, err := c.Client.Gists.Get(ctx, id)
		if err == nil {
			c.noteSearchResponse("gists/"+id, res)
			return gist, nil
		}
		if !c.recoverSearchError("gists/"+id, err, attempt) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("exhausted retry attempts for gist %s", id)
}

func convertGistToSearchResults(gist *github.Gist, hit gistHit, keyword string) []model.SearchResult {
	if gist == nil {
		return nil
	}
	owner := gistOwnerLogin(gist)
	if owner == "" {
		owner = hit.Owner
	}
	repo := strings.Trim(owner+"/"+gist.GetID(), "/")
	if repo == "" {
		repo = "gist/" + hit.ID
	}
	htmlURL := gist.GetHTMLURL()
	if htmlURL == "" && gist.GetID() != "" {
		htmlURL = "https://gist.github.com/" + gist.GetID()
	}
	results := make([]model.SearchResult, 0, len(gist.Files))
	for filename, file := range gist.Files {
		name := file.GetFilename()
		if name == "" {
			name = string(filename)
		}
		fragment := gistFileFragment(file.GetContent(), keyword)
		if keyword != "" && !gistMatchesKeyword(name, gist.GetDescription(), file.GetContent(), keyword) {
			continue
		}
		item := model.SearchResult{
			Repo:    repo,
			RepoUrl: htmlURL,
			Keyword: keyword,
			Path:    name,
			Url:     htmlURL,
			Status:  0,
			Matches: firstLine(gist.GetDescription(), fragment),
		}
		if fragment != "" {
			textMatches := []model.TextMatch{{
				Fragment: &fragment,
				Property: github.String("content"),
			}}
			if encoded, err := json.Marshal(textMatches); err != nil {
				global.GVA_LOG.Error("json.marshal gist text matches error", zap.Error(err))
			} else {
				item.TextMatchesJson = encoded
			}
		}
		results = append(results, item)
	}
	if len(results) == 0 && gist.GetDescription() != "" && gistMatchesKeyword("", gist.GetDescription(), "", keyword) {
		results = append(results, model.SearchResult{
			Repo:    repo,
			RepoUrl: htmlURL,
			Keyword: keyword,
			Path:    "gist",
			Url:     htmlURL,
			Status:  0,
			Matches: gist.GetDescription(),
		})
	}
	return results
}

func gistMatchesKeyword(filename, description, content, keyword string) bool {
	needle := strings.ToLower(strings.TrimSpace(keyword))
	if needle == "" {
		return true
	}
	// Drop filter-style tokens so "company -extension:md" still matches "company".
	fields := strings.Fields(needle)
	needles := make([]string, 0, len(fields))
	for _, field := range fields {
		if strings.HasPrefix(field, "-") || strings.HasPrefix(field, "+") || strings.Contains(field, ":") {
			continue
		}
		needles = append(needles, field)
	}
	if len(needles) == 0 {
		return true
	}
	haystack := strings.ToLower(filename + "\n" + description + "\n" + content)
	for _, part := range needles {
		if !strings.Contains(haystack, part) {
			return false
		}
	}
	return true
}

func gistFileFragment(content, keyword string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}
	lower := strings.ToLower(content)
	idx := strings.Index(lower, strings.ToLower(strings.TrimSpace(keyword)))
	if idx < 0 {
		if len(content) > 400 {
			return content[:400]
		}
		return content
	}
	start := idx - 80
	if start < 0 {
		start = 0
	}
	end := idx + 200
	if end > len(content) {
		end = len(content)
	}
	return content[start:end]
}

func firstLine(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func gistOwnerLogin(gist *github.Gist) string {
	if gist == nil || gist.Owner == nil {
		return ""
	}
	return gist.Owner.GetLogin()
}

func parseGistSearchHTML(body string) []gistHit {
	matches := gistPathRe.FindAllStringSubmatch(body, -1)
	seen := make(map[string]struct{})
	hits := make([]gistHit, 0, len(matches))
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		owner, id := match[1], match[2]
		if owner == "search" || owner == "login" || owner == "join" || owner == "settings" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		hits = append(hits, gistHit{Owner: owner, ID: id})
	}
	return hits
}

func defaultFetchGistSearchPage(token, query string, page int) (string, error) {
	endpoint := fmt.Sprintf("https://gist.github.com/search?q=%s&p=%d", url.QueryEscape(strings.TrimSpace(query)), page)
	var lastErr error
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", "gshark")
		req.Header.Set("Accept", "text/html")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := gistSearchClient.Do(req)
		if err != nil {
			lastErr = err
		} else {
			limited := io.LimitReader(resp.Body, maxGistSearchPageBytes+1)
			body, readErr := io.ReadAll(limited)
			resp.Body.Close()
			if readErr != nil {
				lastErr = readErr
			} else if len(body) > maxGistSearchPageBytes {
				return "", fmt.Errorf("gist search page exceeded %d bytes", maxGistSearchPageBytes)
			} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return string(body), nil
			} else {
				lastErr = fmt.Errorf("gist search returned HTTP %d", resp.StatusCode)
				if resp.StatusCode < 500 {
					return "", lastErr
				}
			}
		}
		if attempt < maxSearchAttempts {
			sleepFn(time.Duration(attempt) * 2 * time.Second)
		}
	}
	return "", lastErr
}

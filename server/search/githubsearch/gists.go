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
	gistSearchUserAgent    = "GShark (+https://github.com/madneal/gshark)"
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
	filterErr, extensionFilters := getFiltersByClass("extension")
	if filterErr != nil {
		global.GVA_LOG.Warn("load gist extension filters failed; continuing without local extension filtering", zap.Error(filterErr))
		extensionFilters = []model.Filter{}
	}
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
		var inserted int
		var repos []string
		var hasMoreRepos bool
		err = client.LoadGistResultsStream(hits, rule.Content, extensionFilters, func(results []model.SearchResult) error {
			stats := service.SaveSearchResultsWithStats(results)
			inserted += stats.Inserted
			var more bool
			repos, more = appendUniqueRepos(repos, stats.Repos)
			hasMoreRepos = hasMoreRepos || more
			return nil
		})
		if err != nil {
			global.GVA_LOG.Error("LoadGistResultsStream error", zap.Error(err))
			scanErrors = append(scanErrors, fmt.Errorf("load gists for %q: %w", rule.Content, err))
		}
		stats := service.NewSaveResultStats()
		stats.Inserted = inserted
		stats.Repos = repos
		global.GVA_LOG.Info(fmt.Sprintf("[Gist] keyword=%q: inserted=%d", rule.Content, inserted))
		content += formatInsertedSummaryWithMore(rule.Content, stats, hasMoreRepos)
	}
	notifyNewResults("Gist 敏感信息报告", content)
	return errors.Join(scanErrors...)
}

func (c *Client) SearchGistHits(query string) ([]gistHit, error) {
	seen := make(map[string]struct{})
	hits := make([]gistHit, 0)
	for page := 1; page <= maxGistSearchPages; page++ {
		body, err := fetchGistSearchPage(query, page)
		if err != nil {
			if page == 1 {
				return nil, err
			}
			global.GVA_LOG.Warn("gist HTML search stopped after a later page error", zap.Int("page", page), zap.Error(err))
			break
		}
		pageHits := parseGistSearchHTML(body)
		if len(pageHits) == 0 {
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

func (c *Client) LoadGistResults(hits []gistHit, keyword string) ([]model.SearchResult, error) {
	results := make([]model.SearchResult, 0)
	err := c.LoadGistResultsStream(hits, keyword, nil, func(batch []model.SearchResult) error {
		results = append(results, batch...)
		return nil
	})
	return results, err
}

// LoadGistResultsStream fetches and converts one Gist at a time. The callback
// runs before the next Gist is loaded so callers can persist results without
// retaining the complete search response in memory.
func (c *Client) LoadGistResultsStream(hits []gistHit, keyword string, extensionFilters []model.Filter, onBatch func([]model.SearchResult) error) error {
	if onBatch == nil {
		return fmt.Errorf("Gist result callback is nil")
	}
	if extensionFilters == nil {
		if _, extensionFilters = getFiltersByClass("extension"); extensionFilters == nil {
			extensionFilters = []model.Filter{}
		}
	}
	var loadErrors []error
	ctx := context.Background()
	for _, hit := range hits {
		gist, err := c.getGist(ctx, hit.ID)
		if err != nil {
			loadErrors = append(loadErrors, fmt.Errorf("gist %s: %w", hit.ID, err))
			continue
		}
		results := convertGistToSearchResultsWithFilters(gist, hit, keyword, extensionFilters)
		if len(results) == 0 {
			continue
		}
		if err := onBatch(results); err != nil {
			return err
		}
	}
	return errors.Join(loadErrors...)
}

func (c *Client) getGist(ctx context.Context, id string) (*github.Gist, error) {
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		gist, res, err := c.Client.Gists.Get(ctx, id)
		if err == nil {
			c.noteSearchResponse("gists/"+id, res)
			return gist, nil
		}
		if isNonRetryableAPIError(err) {
			return nil, err
		}
		if !c.recoverSearchError("gists/"+id, err, attempt) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("exhausted retry attempts for gist %s", id)
}

func isNonRetryableAPIError(err error) bool {
	var errResp *github.ErrorResponse
	if !errors.As(err, &errResp) || errResp.Response == nil {
		return false
	}
	code := errResp.Response.StatusCode
	return code != http.StatusTooManyRequests && code >= 400 && code < 500
}

func convertGistToSearchResults(gist *github.Gist, hit gistHit, keyword string) []model.SearchResult {
	_, filters := getFiltersByClass("extension")
	return convertGistToSearchResultsWithFilters(gist, hit, keyword, filters)
}

func convertGistToSearchResultsWithFilters(gist *github.Gist, hit gistHit, keyword string, extensionFilters []model.Filter) []model.SearchResult {
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
		if !gistAllowsExtensionWithFilters(name, extensionFilters) {
			continue
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
	if len(results) == 0 && gistAllowsExtensionWithFilters("", extensionFilters) && gist.GetDescription() != "" && gistMatchesKeyword("", gist.GetDescription(), "", keyword) {
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

var gistQueryQualifiers = []string{"extension:", "language:", "filename:", "user:", "in:"}
var gistQueryTokenRe = regexp.MustCompile(`(?:[+-]?"[^"]+"|[^\s]+)`)

func gistMatchesKeyword(filename, description, content, keyword string) bool {
	needles := gistNeedles(keyword)
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

func gistNeedles(keyword string) []string {
	needles := make([]string, 0)
	for _, field := range gistQueryTokenRe.FindAllString(strings.ToLower(strings.TrimSpace(keyword)), -1) {
		denied := strings.HasPrefix(field, "-")
		field = strings.TrimPrefix(strings.TrimPrefix(field, "+"), "-")
		field = strings.Trim(field, `"`)
		if field == "" || isGistQueryQualifier(field) {
			continue
		}
		if denied {
			continue
		}
		needles = append(needles, field)
	}
	return needles
}

func isGistQueryQualifier(field string) bool {
	for _, prefix := range gistQueryQualifiers {
		if strings.HasPrefix(field, prefix) {
			return true
		}
	}
	return false
}

func gistAllowsExtension(filename string) bool {
	_, filters := getFiltersByClass("extension")
	return gistAllowsExtensionWithFilters(filename, filters)
}

func gistAllowsExtensionWithFilters(filename string, filters []model.Filter) bool {
	if len(filters) == 0 {
		return true
	}
	ext := fileExtension(filename)
	hasAllow := false
	allowed := false
	for _, filter := range filters {
		for _, raw := range strings.Split(filter.Content, ",") {
			want := strings.TrimPrefix(strings.TrimSpace(strings.ToLower(raw)), ".")
			if want == "" {
				continue
			}
			if filter.FilterType == "blacklist" {
				if ext == want {
					return false
				}
				continue
			}
			hasAllow = true
			if ext == want {
				allowed = true
			}
		}
	}
	if hasAllow {
		return allowed
	}
	return true
}

func fileExtension(filename string) string {
	base := filename
	if i := strings.LastIndex(filename, "/"); i >= 0 {
		base = filename[i+1:]
	}
	dot := strings.LastIndex(base, ".")
	if dot < 0 || dot == len(base)-1 {
		return ""
	}
	return strings.ToLower(base[dot+1:])
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

func defaultFetchGistSearchPage(query string, page int) (string, error) {
	endpoint := fmt.Sprintf("https://gist.github.com/search?q=%s&p=%d", url.QueryEscape(strings.TrimSpace(query)), page)
	var lastErr error
	for attempt := 1; attempt <= maxSearchAttempts; attempt++ {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", gistSearchUserAgent)
		req.Header.Set("Accept", "text/html")
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
				if resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode < 500 {
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

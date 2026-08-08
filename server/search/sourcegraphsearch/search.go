package sourcegraphsearch

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

const (
	defaultSourcegraphURL = "https://sourcegraph.com"
	defaultRequestTimeout = 2 * time.Minute
)

var (
	sourcegraphHTTPClient = &http.Client{Timeout: defaultRequestTimeout}
	getSourcegraphRules   = service.GetValidRulesByType
)

// streamMatch is the content/path match shape emitted by Sourcegraph's SSE
// search endpoint. Sourcegraph can emit either lineMatches or chunkMatches;
// both are converted into GShark's existing text-match format.
type streamMatch struct {
	Type         string       `json:"type"`
	Path         string       `json:"path"`
	Repository   string       `json:"repository"`
	Commit       string       `json:"commit"`
	LineMatches  []lineMatch  `json:"lineMatches"`
	ChunkMatches []chunkMatch `json:"chunkMatches"`
}

type lineMatch struct {
	Line       string `json:"line"`
	LineNumber int    `json:"lineNumber"`
}

type chunkMatch struct {
	Content      string `json:"content"`
	ContentStart struct {
		Line int `json:"line"`
	} `json:"contentStart"`
	BestLineMatch int `json:"bestLineMatch"`
}

type progressEvent struct {
	Done              bool `json:"done"`
	MatchCount        int  `json:"matchCount"`
	RepositoriesCount int  `json:"repositoriesCount"`
	Skipped           []struct {
		Reason  string `json:"reason"`
		Title   string `json:"title"`
		Message string `json:"message"`
	} `json:"skipped"`
}

type alertEvent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// RunTask scans all repositories indexed by Sourcegraph for both the new
// sourcegraph rule type and the legacy searchcode rule type. The latter keeps
// existing deployments working without a database migration.
func RunTask() model.ScanOutcome {
	rules, err := loadRules()
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType sourcegraph err", zap.Error(err))
		return model.ScanFailed("Failed to load Sourcegraph rules: " + err.Error())
	}
	if len(rules) == 0 {
		message := "No enabled Sourcegraph rules; provider skipped"
		global.GVA_LOG.Info(message)
		return model.ScanSkipped(message)
	}

	var scanErrors []error
	partial := false
	for _, rule := range rules {
		global.GVA_LOG.Info("Search all indexed repositories in Sourcegraph", zap.String("rule", rule.Content))
		results, warnings, searchErr := SearchForSourcegraph(rule, sourcegraphHTTPClient)
		if searchErr != nil {
			scanErrors = append(scanErrors, fmt.Errorf("search %q: %w", rule.Content, searchErr))
			continue
		}
		if len(warnings) > 0 {
			partial = true
			for _, warning := range warnings {
				global.GVA_LOG.Warn("Sourcegraph returned partial search information", zap.String("rule", rule.Content), zap.String("warning", warning))
			}
		}
		SaveResults(results, &rule.Content)
	}
	if err := errors.Join(scanErrors...); err != nil {
		return model.ScanFailed("Sourcegraph scan completed with errors: " + err.Error())
	}
	message := fmt.Sprintf("Completed %d Sourcegraph rules across indexed repositories", len(rules))
	if partial {
		message += " (upstream reported result limits; see scanner logs)"
	}
	global.GVA_LOG.Info(message)
	return model.ScanSuccess(message)
}

func loadRules() ([]model.Rule, error) {
	seen := make(map[string]struct{})
	rules := make([]model.Rule, 0)
	for _, ruleType := range []string{"sourcegraph", "searchcode"} {
		err, typeRules := getSourcegraphRules(ruleType)
		if err != nil {
			return nil, err
		}
		for _, rule := range typeRules {
			key := rule.Content
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func sourcegraphBaseURL() string {
	configured := strings.TrimSpace(global.GVA_CONFIG.Search.SourcegraphURL)
	if configured == "" {
		return defaultSourcegraphURL
	}
	return strings.TrimRight(configured, "/")
}

func sourcegraphToken() string {
	if configured := strings.TrimSpace(global.GVA_CONFIG.Search.SourcegraphToken); configured != "" {
		return configured
	}
	return strings.TrimSpace(os.Getenv("SOURCEGRAPH_TOKEN"))
}

// SearchForSourcegraph executes one rule against every repository indexed by
// Sourcegraph. count:all is intentional: unlike the retired Searchcode API,
// this is a global search and does not require a repository list.
func SearchForSourcegraph(rule model.Rule, client *http.Client) ([]*model.SearchResult, []string, error) {
	query := globalQuery(rule.Content)
	params := url.Values{}
	params.Set("q", query)
	params.Set("v", "V3")
	params.Set("t", "keyword")
	params.Set("cm", "true")
	params.Set("cl", "2")

	endpoint := sourcegraphBaseURL() + "/.api/search/stream?" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	if token := sourcegraphToken(); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, errors.New("Sourcegraph returned a nil response")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, nil, fmt.Errorf("Sourcegraph request returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	results := make([]*model.SearchResult, 0)
	warnings := make([]string, 0)
	var streamErr error
	if err := parseStream(resp.Body, func(eventName, data string) error {
		switch eventName {
		case "matches":
			var matches []streamMatch
			if err := json.Unmarshal([]byte(data), &matches); err != nil {
				return fmt.Errorf("decode Sourcegraph matches: %w", err)
			}
			for _, match := range matches {
				if result := convertMatch(match); result != nil {
					results = append(results, result)
				}
			}
		case "progress":
			var progress progressEvent
			if err := json.Unmarshal([]byte(data), &progress); err != nil {
				return fmt.Errorf("decode Sourcegraph progress: %w", err)
			}
			for _, skipped := range progress.Skipped {
				warning := skipped.Message
				if warning == "" {
					warning = skipped.Title
				}
				if warning != "" {
					warnings = append(warnings, skipped.Reason+": "+warning)
				}
			}
		case "alert":
			var alert alertEvent
			if err := json.Unmarshal([]byte(data), &alert); err != nil {
				return fmt.Errorf("decode Sourcegraph alert: %w", err)
			}
			if alert.Severity == "error" {
				streamErr = fmt.Errorf("%s: %s", alert.Title, alert.Description)
			}
		}
		return nil
	}); err != nil {
		return results, warnings, err
	}
	if streamErr != nil {
		return results, warnings, streamErr
	}
	return results, uniqueWarnings(warnings), nil
}

func globalQuery(content string) string {
	query := strings.TrimSpace(content)
	if !strings.Contains(query, "count:") {
		query += " count:all"
	}
	if !strings.Contains(query, "fork:") {
		query += " fork:yes"
	}
	if !strings.Contains(query, "archived:") {
		query += " archived:yes"
	}
	return strings.TrimSpace(query)
}

func parseStream(body io.Reader, handle func(eventName, data string) error) error {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 64*1024), 8*1024*1024)
	eventName := ""
	data := ""
	flush := func() error {
		if eventName == "" || data == "" {
			eventName, data = "", ""
			return nil
		}
		err := handle(eventName, data)
		eventName, data = "", ""
		return err
	}
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if line == "" {
			if err := flush(); err != nil {
				return err
			}
			continue
		}
		if strings.HasPrefix(line, "event:") {
			eventName = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}
		if strings.HasPrefix(line, "data:") {
			data = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return flush()
}

func convertMatch(match streamMatch) *model.SearchResult {
	if match.Repository == "" || match.Path == "" || match.Type != "content" {
		return nil
	}
	fragment, lineNumber := matchFragment(match)
	if fragment == "" {
		return nil
	}
	repoURL := canonicalRepoURL(match.Repository)
	resultURL := sourcegraphResultURL(match.Repository, match.Commit, match.Path, lineNumber)
	textMatches := []model.TextMatch{{Fragment: &fragment}}
	encoded, err := json.Marshal(textMatches)
	if err != nil {
		global.GVA_LOG.Error("marshal Sourcegraph text match", zap.Error(err))
		return nil
	}
	return &model.SearchResult{
		Repo:            match.Repository,
		RepoUrl:         repoURL,
		Path:            match.Path,
		Url:             resultURL,
		TextMatchesJson: encoded,
		Status:          0,
	}
}

func matchFragment(match streamMatch) (string, int) {
	if len(match.ChunkMatches) > 0 {
		lines := make([]string, 0, len(match.ChunkMatches))
		lineNumber := 0
		for _, chunk := range match.ChunkMatches {
			if chunk.Content != "" {
				lines = append(lines, chunk.Content)
			}
			if lineNumber == 0 {
				lineNumber = chunk.ContentStart.Line + 1
				if chunk.BestLineMatch > 0 {
					lineNumber = chunk.BestLineMatch + 1
				}
			}
		}
		return strings.Join(lines, "\n"), lineNumber
	}
	lines := make([]string, 0, len(match.LineMatches))
	lineNumber := 0
	for _, line := range match.LineMatches {
		if lineNumber == 0 {
			lineNumber = line.LineNumber
		}
		lines = append(lines, fmt.Sprintf("%d: %s", line.LineNumber, line.Line))
	}
	return strings.Join(lines, "\n"), lineNumber
}

func canonicalRepoURL(repository string) string {
	if strings.HasPrefix(repository, "http://") || strings.HasPrefix(repository, "https://") {
		return repository
	}
	return "https://" + strings.TrimPrefix(repository, "//")
}

func sourcegraphResultURL(repository, commit, path string, line int) string {
	result := sourcegraphBaseURL() + "/" + strings.TrimPrefix(repository, "/") + "/-/blob/" + strings.TrimPrefix(path, "/")
	if commit != "" {
		result += "?L=" + url.QueryEscape(fmt.Sprint(line)) + "&rev=" + url.QueryEscape(commit)
	} else if line > 0 {
		result += "?L=" + url.QueryEscape(fmt.Sprint(line))
	}
	return result
}

func uniqueWarnings(warnings []string) []string {
	seen := make(map[string]struct{}, len(warnings))
	unique := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		if _, exists := seen[warning]; exists {
			continue
		}
		seen[warning] = struct{}{}
		unique = append(unique, warning)
	}
	return unique
}

func SaveResults(results []*model.SearchResult, keyword *string) {
	if len(results) == 0 {
		return
	}
	stats := service.SaveSearchResultPointersWithStats(results, *keyword)
	global.GVA_LOG.Info(stats.Summary(*keyword, "Sourcegraph"))
}

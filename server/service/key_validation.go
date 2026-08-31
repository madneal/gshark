package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
)

const (
	ValidationSkipped    = "skipped"
	ValidationValid      = "valid"
	ValidationInvalid    = "invalid"
	ValidationDeferred   = "deferred"
	keyValidationTimeout = 15 * time.Second
)

var keyValidationClient = &http.Client{Timeout: keyValidationTimeout}

type keyValidationCacheEntry struct {
	status    string
	expiresAt time.Time
}

var keyValidationCache = struct {
	sync.Mutex
	items map[string]keyValidationCacheEntry
}{items: make(map[string]keyValidationCacheEntry)}

var keyPatterns = map[string]*regexp.Regexp{
	"github":      regexp.MustCompile(`(?:gh[pousr]_[A-Za-z0-9_]{16,}|github_pat_[A-Za-z0-9_]{20,})`),
	"gitlab":      regexp.MustCompile(`glpat-[A-Za-z0-9_-]{20,}`),
	"sourcegraph": regexp.MustCompile(`sgp_[A-Za-z0-9_-]{20,}`),
	"postman":     regexp.MustCompile(`PMAK-[A-Za-z0-9_-]{20,}`),
}

// ValidateRuleResult extracts the key selected by a rule's match pattern and
// verifies it against that platform's authenticated API. A deferred result is
// deliberately not persisted: transient errors and rate limits must be
// retried on a later scan instead of being treated as a valid secret.
func ValidateRuleResult(rule model.Rule, result model.SearchResult) (string, error) {
	provider := ruleValidationType(rule)
	if provider == "" {
		return ValidationSkipped, nil
	}
	if _, ok := keyPatterns[provider]; !ok {
		return ValidationDeferred, fmt.Errorf("unsupported key validation provider %q", provider)
	}
	candidates := ruleKeyCandidates(rule, SearchResultContent(result), provider)
	if len(candidates) == 0 {
		return ValidationInvalid, nil
	}
	deferred := false
	for _, candidate := range candidates {
		status, err := validateKeyCached(context.Background(), provider, candidate)
		if err != nil {
			deferred = true
			continue
		}
		if status == ValidationValid {
			return ValidationValid, nil
		}
		if status == ValidationDeferred {
			deferred = true
		}
	}
	if deferred {
		return ValidationDeferred, errors.New("key validation could not be completed")
	}
	return ValidationInvalid, nil
}

func validateKeyCached(ctx context.Context, provider, key string) (string, error) {
	cacheScope := provider + "\x00" + global.GVA_CONFIG.System.GitlabBase + "\x00" + global.GVA_CONFIG.Search.SourcegraphURL + "\x00" + key
	digest := sha256.Sum256([]byte(cacheScope))
	cacheKey := fmt.Sprintf("%x", digest[:])
	now := time.Now()
	keyValidationCache.Lock()
	if cached, ok := keyValidationCache.items[cacheKey]; ok && now.Before(cached.expiresAt) {
		keyValidationCache.Unlock()
		if cached.status == ValidationDeferred {
			return cached.status, errors.New("cached key validation failure")
		}
		return cached.status, nil
	}
	keyValidationCache.Unlock()

	status, err := validateKey(ctx, provider, key)
	ttl := 30 * time.Second
	if status == ValidationValid || status == ValidationInvalid {
		ttl = 10 * time.Minute
	}
	keyValidationCache.Lock()
	keyValidationCache.items[cacheKey] = keyValidationCacheEntry{status: status, expiresAt: now.Add(ttl)}
	if len(keyValidationCache.items) > 10000 {
		for itemKey, item := range keyValidationCache.items {
			if now.After(item.expiresAt) {
				delete(keyValidationCache.items, itemKey)
			}
		}
	}
	keyValidationCache.Unlock()
	return status, err
}

func ruleValidationType(rule model.Rule) string {
	if value := strings.ToLower(strings.TrimSpace(rule.ValidationType)); value != "" {
		if value == "none" {
			return ""
		}
		return value
	}
	text := strings.ToLower(rule.MatchPattern + "\n" + rule.Content)
	switch {
	case strings.Contains(text, "github_pat_") || strings.Contains(text, "ghp_") ||
		strings.Contains(text, "gho_") || strings.Contains(text, "ghu_") ||
		strings.Contains(text, "ghs_") || strings.Contains(text, "ghr_"):
		return "github"
	case strings.Contains(text, "glpat-"):
		return "gitlab"
	case strings.Contains(text, "sgp_"):
		return "sourcegraph"
	case strings.Contains(text, "pmak-"):
		return "postman"
	default:
		return ""
	}
}

func ruleKeyCandidates(rule model.Rule, content, provider string) []string {
	seen := make(map[string]struct{})
	add := func(value string) {
		value = strings.Trim(strings.TrimSpace(value), "\"'`.,;)")
		if value == "" || len(value) > 1000 {
			return
		}
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}
		}
	}
	if strings.TrimSpace(rule.MatchPattern) != "" {
		if pattern, err := regexp.Compile(rule.MatchPattern); err == nil {
			keyGroup := -1
			for _, name := range []string{"key", "token", "value", "secret"} {
				if index := pattern.SubexpIndex(name); index >= 0 {
					keyGroup = index
					break
				}
			}
			for _, match := range pattern.FindAllStringSubmatch(content, -1) {
				if keyGroup > 0 && keyGroup < len(match) {
					add(match[keyGroup])
				} else if pattern.NumSubexp() == 1 && len(match) > 1 {
					add(match[1])
				} else {
					for _, candidate := range keyPatterns[provider].FindAllString(match[0], -1) {
						add(candidate)
					}
				}
			}
		}
	}
	if len(seen) == 0 {
		for _, candidate := range keyPatterns[provider].FindAllString(content, -1) {
			add(candidate)
		}
	}
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	return values
}

func validateKey(ctx context.Context, provider, key string) (string, error) {
	endpoint, header, err := validationEndpoint(provider, key)
	if err != nil {
		return ValidationDeferred, err
	}
	method := http.MethodGet
	var body io.Reader
	if provider == "sourcegraph" {
		method = http.MethodPost
		var payload []byte
		payload, err = json.Marshal(map[string]string{"query": "query { currentUser { username } }"})
		if err != nil {
			return ValidationDeferred, err
		}
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return ValidationDeferred, err
	}
	req.Header.Set(header[0], header[1]+key)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GShark-key-validator")
	if provider == "sourcegraph" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := keyValidationClient.Do(req)
	if err != nil {
		return ValidationDeferred, err
	}
	defer resp.Body.Close()
	responseBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if readErr != nil {
		return ValidationDeferred, readErr
	}
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		if provider == "sourcegraph" {
			var payload struct {
				Data struct {
					CurrentUser *struct {
						Username string `json:"username"`
					} `json:"currentUser"`
				} `json:"data"`
				Errors []json.RawMessage `json:"errors"`
			}
			if json.Unmarshal(responseBody, &payload) != nil || len(payload.Errors) > 0 || payload.Data.CurrentUser == nil {
				return ValidationInvalid, nil
			}
		}
		return ValidationValid, nil
	case resp.StatusCode == http.StatusUnauthorized:
		return ValidationInvalid, nil
	default:
		return ValidationDeferred, fmt.Errorf("validation API returned HTTP %d", resp.StatusCode)
	}
}

func validationEndpoint(provider, key string) (string, [2]string, error) {
	switch provider {
	case "github":
		if strings.HasPrefix(key, "ghr_") {
			return "", [2]string{}, errors.New("GitHub refresh tokens cannot authenticate REST API requests")
		}
		if strings.HasPrefix(key, "ghs_") {
			return "https://api.github.com/installation/repositories", [2]string{"Authorization", "Bearer "}, nil
		}
		return "https://api.github.com/user", [2]string{"Authorization", "Bearer "}, nil
	case "gitlab":
		base := strings.TrimRight(strings.TrimSpace(global.GVA_CONFIG.System.GitlabBase), "/")
		if base == "" {
			base = "https://gitlab.com"
		}
		return base + "/api/v4/user", [2]string{"PRIVATE-TOKEN", ""}, nil
	case "sourcegraph":
		base := strings.TrimRight(strings.TrimSpace(global.GVA_CONFIG.Search.SourcegraphURL), "/")
		if base == "" {
			base = "https://sourcegraph.com"
		}
		return base + "/.api/graphql", [2]string{"Authorization", "token "}, nil
	case "postman":
		return "https://api.getpostman.com/me", [2]string{"X-Api-Key", ""}, nil
	default:
		return "", [2]string{}, fmt.Errorf("unsupported key validation provider %q", provider)
	}
}

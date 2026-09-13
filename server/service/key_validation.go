package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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

type keyDetector struct {
	provider string
	pattern  *regexp.Regexp
}

type detectedKey struct {
	provider string
	value    string
}

// Built-in detectors make key validation automatic. Search rules only decide
// what to search for; they do not need to know how a provider validates keys.
var keyDetectors = []keyDetector{
	{provider: "github", pattern: regexp.MustCompile(`(?:gh[pousr]_[A-Za-z0-9_]{16,}|github_pat_[A-Za-z0-9_]{20,})`)},
	{provider: "gitlab", pattern: regexp.MustCompile(`glpat-[A-Za-z0-9_-]{20,}`)},
	{provider: "sourcegraph", pattern: regexp.MustCompile(`sgp_[A-Za-z0-9_-]{20,}`)},
	{provider: "postman", pattern: regexp.MustCompile(`PMAK-[A-Za-z0-9_-]{20,}`)},
}

// ValidateSearchResultKeys detects supported keys in a result and validates
// them with the owning provider. Results without a known key keep the existing
// ingest behavior. If keys are detected, at least one must validate before the
// result can be stored.
func ValidateSearchResultKeys(result model.SearchResult) (string, error) {
	keys := detectKeys(SearchResultContent(result))
	if len(keys) == 0 {
		return ValidationSkipped, nil
	}

	var deferredErr error
	for _, key := range keys {
		status, err := validateKey(context.Background(), key.provider, key.value)
		if status == ValidationValid {
			return ValidationValid, nil
		}
		if status == ValidationDeferred || err != nil {
			if err == nil {
				err = errors.New("validation deferred")
			}
			deferredErr = fmt.Errorf("%s key validation failed: %w", key.provider, err)
		}
	}
	if deferredErr != nil {
		return ValidationDeferred, deferredErr
	}
	return ValidationInvalid, nil
}

func detectKeys(content string) []detectedKey {
	seen := make(map[detectedKey]struct{})
	keys := make([]detectedKey, 0)
	for _, detector := range keyDetectors {
		for _, value := range detector.pattern.FindAllString(content, -1) {
			key := detectedKey{provider: detector.provider, value: value}
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}
	return keys
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

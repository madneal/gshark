package githubsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/model"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	global.GVA_LOG = zap.NewNop()
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	global.GVA_VP = initialize.Viper("../../config.yaml") // 初始化Viper
	global.GVA_LOG = initialize.Zap()
	global.GVA_DB = initialize.Gorm()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	rules := make([]model.Rule, 0)
	rules = append(rules, model.Rule{
		Content: "meituan",
	})
	Search(rules)
}

func TestSearchWithNoRules(t *testing.T) {
	Search(nil)
}

func TestGetGithubClientWithoutTokenReturnsError(t *testing.T) {
	original := listGithubTokens
	listGithubTokens = func(string) (error, []model.Token) { return nil, nil }
	t.Cleanup(func() { listGithubTokens = original })

	client, err := GetGithubClient()
	if client != nil {
		t.Fatal("expected no GitHub client without a token")
	}
	if err == nil {
		t.Fatal("expected an error without a GitHub token")
	}
}

func TestNextClientWithoutTokenDoesNotPanic(t *testing.T) {
	original := listGithubTokens
	listGithubTokens = func(string) (error, []model.Token) { return nil, nil }
	t.Cleanup(func() { listGithubTokens = original })

	client := &Client{Token: "removed-token"}
	nextClient, nextToken := client.NextClient()
	if nextClient != nil || nextToken != "" {
		t.Fatalf("expected no replacement client, got client=%v token=%q", nextClient, nextToken)
	}
}

func TestGithubHTTPTransportPreservesDefaultProxy(t *testing.T) {
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok || defaultTransport.Proxy == nil {
		t.Fatal("expected the default HTTP transport to provide an environment proxy")
	}
	if newGithubHTTPTransport().Proxy == nil {
		t.Fatal("expected GitHub transport to preserve the default environment proxy")
	}
}

func TestRunTaskWithoutRulesDoesNotSleep(t *testing.T) {
	originalRules := getGithubRules
	getGithubRules = func(string) (error, []model.Rule) { return nil, nil }
	t.Cleanup(func() { getGithubRules = originalRules })

	slept := stubSleep(t)
	outcome := RunTask()
	if len(*slept) != 0 {
		t.Fatalf("RunTask slept without GitHub rules: %v", *slept)
	}
	if outcome.Status != model.ScanStatusSkipped {
		t.Fatalf("outcome status = %q, want skipped", outcome.Status)
	}
}

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, func() *github.Client) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	newGithubClient := func() *github.Client {
		ghClient := github.NewClient(server.Client())
		ghClient.BaseURL = baseURL
		return ghClient
	}

	return &Client{Client: newGithubClient(), Token: "test-token"}, newGithubClient
}

func stubSleep(t *testing.T) *[]time.Duration {
	t.Helper()
	var slept []time.Duration
	original := sleepFn
	sleepFn = func(d time.Duration) { slept = append(slept, d) }
	t.Cleanup(func() { sleepFn = original })
	return &slept
}

func writeRateLimitResponse(w http.ResponseWriter, resetAt time.Time) {
	w.Header().Set("X-RateLimit-Remaining", "0")
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt.Unix()))
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "API rate limit exceeded"})
}

func writeAbuseRateLimitResponse(w http.ResponseWriter, retryAfterSeconds int) {
	w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfterSeconds))
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message":           "You have exceeded a secondary rate limit",
		"documentation_url": "https://docs.github.com/rest/overview/rate-limits-for-the-rest-api#secondary-rate-limits",
	})
}

func writeSearchSuccess(w http.ResponseWriter, result github.CodeSearchResult) {
	w.Header().Set("X-RateLimit-Remaining", "100")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func TestSearchCodeByOptRetriesAfterRotatingOnPrimaryRateLimit(t *testing.T) {
	stubSleep(t)
	requestCount := 0
	client, newGithubClient := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			writeRateLimitResponse(w, time.Now().Add(time.Hour))
			return
		}
		writeSearchSuccess(w, github.CodeSearchResult{
			Total:       github.Int(1),
			CodeResults: []*github.CodeResult{{Path: github.String("main.go")}},
		})
	})

	rotated := false
	client.rotate = func() bool {
		rotated = true
		// A real token rotation builds a brand-new *github.Client (see
		// InitGithubClient), which also resets go-github's internal
		// per-category rate-limit cache. Mirror that here so the retried
		// request isn't short-circuited by the cached 403 from attempt 1.
		client.Client = newGithubClient()
		return true
	}

	result, nextPage := client.searchCodeByOpt(context.Background(), "test query", github.SearchOptions{})
	if !rotated {
		t.Fatal("expected rotate to be invoked after hitting the primary rate limit")
	}
	if requestCount != 2 {
		t.Fatalf("expected the search request to be retried once, got %d attempts", requestCount)
	}
	if result == nil || len(result.CodeResults) != 1 {
		t.Fatalf("expected the retried request's result to be returned, got %#v", result)
	}
	if nextPage != 0 {
		t.Fatalf("nextPage = %d, want 0", nextPage)
	}
}

func TestSearchCodeByOptSleepsAndRetriesWhenRotationUnavailable(t *testing.T) {
	slept := stubSleep(t)
	requestCount := 0
	// resetAt is already in the past: sleepFn is stubbed out (no real wait),
	// so the reset time must have already elapsed or go-github's own
	// client-side rate-limit cache would short-circuit the retry itself.
	resetAt := time.Now().Add(-time.Second)
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			writeRateLimitResponse(w, resetAt)
			return
		}
		writeSearchSuccess(w, github.CodeSearchResult{Total: github.Int(0)})
	})
	client.rotate = func() bool { return false }

	result, _ := client.searchCodeByOpt(context.Background(), "test query", github.SearchOptions{})
	if len(*slept) != 1 {
		t.Fatalf("expected exactly one sleep call, got %d", len(*slept))
	}
	if requestCount != 2 {
		t.Fatalf("expected the search request to be retried once, got %d attempts", requestCount)
	}
	if result == nil {
		t.Fatal("expected a result after the retried request succeeded")
	}
}

func TestSearchCodeByOptRetriesAfterRotatingOnSecondaryRateLimit(t *testing.T) {
	stubSleep(t)
	requestCount := 0
	client, newGithubClient := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			writeAbuseRateLimitResponse(w, 1)
			return
		}
		writeSearchSuccess(w, github.CodeSearchResult{Total: github.Int(0)})
	})

	rotated := false
	client.rotate = func() bool {
		rotated = true
		client.Client = newGithubClient()
		return true
	}

	_, _ = client.searchCodeByOpt(context.Background(), "test query", github.SearchOptions{})
	if !rotated {
		t.Fatal("expected rotate to be invoked after hitting the secondary rate limit")
	}
	if requestCount != 2 {
		t.Fatalf("expected the search request to be retried once, got %d attempts", requestCount)
	}
}

func TestSearchCodeByOptGivesUpAfterMaxAttempts(t *testing.T) {
	stubSleep(t)
	requestCount := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		writeRateLimitResponse(w, time.Now().Add(time.Hour))
	})
	client.rotate = func() bool { return false }

	result, nextPage := client.searchCodeByOpt(context.Background(), "test query", github.SearchOptions{})
	if result != nil || nextPage != 0 {
		t.Fatalf("expected no result after exhausting retries, got (%#v, %d)", result, nextPage)
	}
	// Without rotation or an elapsed reset, go-github's own client-side rate
	// cache short-circuits attempts 2+ locally, so only the first attempt
	// reaches the server; the important assertion is that the retry loop
	// still terminates (bounded by maxSearchAttempts) instead of looping
	// forever or hammering the server on every attempt.
	if requestCount < 1 || requestCount > maxSearchAttempts {
		t.Fatalf("expected between 1 and %d attempts to reach the server, got %d", maxSearchAttempts, requestCount)
	}
}

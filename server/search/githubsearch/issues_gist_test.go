package githubsearch

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/madneal/gshark/model"
)

func stubFilters(t *testing.T, filters map[string][]model.Filter) {
	t.Helper()
	original := getFiltersByClass
	getFiltersByClass = func(class string) (error, []model.Filter) {
		return nil, filters[class]
	}
	t.Cleanup(func() { getFiltersByClass = original })
}

func TestBuildIssueQueryAddsDefaultInQualifier(t *testing.T) {
	stubFilters(t, nil)
	got, err := BuildIssueQuery("acme")
	if err != nil {
		t.Fatal(err)
	}
	if got != "acme in:title,body,comments" {
		t.Fatalf("BuildIssueQuery() = %q", got)
	}
}

func TestBuildIssueQueryKeepsExplicitInQualifier(t *testing.T) {
	stubFilters(t, nil)
	got, err := BuildIssueQuery("acme in:body is:pr")
	if err != nil {
		t.Fatal(err)
	}
	if got != "acme in:body is:pr" {
		t.Fatalf("BuildIssueQuery() = %q", got)
	}
}

func TestBuildIssueQueryIgnoresInInsideOtherTokens(t *testing.T) {
	stubFilters(t, nil)
	got, err := BuildIssueQuery("login:admin")
	if err != nil {
		t.Fatal(err)
	}
	if got != "login:admin in:title,body,comments" {
		t.Fatalf("BuildIssueQuery() = %q", got)
	}
}

func TestBuildIssueQueryAppliesKeywordFiltersOnly(t *testing.T) {
	stubFilters(t, map[string][]model.Filter{
		"extension": {{FilterType: "blacklist", Content: "md"}},
		"keyword":   {{FilterType: "black", Content: "example"}},
	})
	got, err := BuildIssueQuery("acme")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, " NOT example") {
		t.Fatalf("expected keyword deny, got %q", got)
	}
	if strings.Contains(got, "extension:") {
		t.Fatalf("issue query should not use extension filters, got %q", got)
	}
}

func TestConvertIssuesToSearchResultsMapsIssueAndPull(t *testing.T) {
	results := ConvertIssuesToSearchResults([]*github.IssuesSearchResult{{
		Issues: []*github.Issue{
			{
				Number:        github.Int(12),
				Title:         github.String("leaked token"),
				Body:          github.String("AKIAEXAMPLE"),
				HTMLURL:       github.String("https://github.com/acme/app/issues/12"),
				RepositoryURL: github.String("https://api.github.com/repos/acme/app"),
			},
			{
				Number:           github.Int(8),
				Title:            github.String("add secret"),
				HTMLURL:          github.String("https://github.com/acme/app/pull/8"),
				RepositoryURL:    github.String("https://api.github.com/repos/acme/app"),
				PullRequestLinks: &github.PullRequestLinks{URL: github.String("https://api.github.com/repos/acme/app/pulls/8")},
			},
		},
	}}, "acme")

	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if results[0].Path != "issues/12" || results[0].Repo != "acme/app" {
		t.Fatalf("issue mapping = %#v", results[0])
	}
	if results[1].Path != "pull/8" || results[1].Url != "https://github.com/acme/app/pull/8" {
		t.Fatalf("pull mapping = %#v", results[1])
	}
}

func TestSearchIssuesByOptRetriesAfterRotatingOnPrimaryRateLimit(t *testing.T) {
	stubSleep(t)
	requestCount := 0
	client, newGithubClient := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			writeRateLimitResponse(w, time.Now().Add(time.Hour))
			return
		}
		if !strings.Contains(r.URL.Path, "/search/issues") {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("X-RateLimit-Remaining", "100")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(github.IssuesSearchResult{
			Total: github.Int(1),
			Issues: []*github.Issue{{
				Number:  github.Int(1),
				HTMLURL: github.String("https://github.com/acme/app/issues/1"),
			}},
		})
	})
	client.rotate = func() bool {
		client.Client = newGithubClient()
		return true
	}

	result, nextPage := client.searchIssuesByOpt(context.Background(), "acme", github.SearchOptions{})
	if requestCount != 2 {
		t.Fatalf("requestCount = %d, want 2", requestCount)
	}
	if result == nil || len(result.Issues) != 1 {
		t.Fatalf("result = %#v", result)
	}
	if nextPage != 0 {
		t.Fatalf("nextPage = %d, want 0", nextPage)
	}
}

func TestSearchIssuesStreamInvokesCallbackBeforeNextPage(t *testing.T) {
	requestCount := 0
	callbackCount := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("X-RateLimit-Remaining", "100")
		if requestCount == 1 {
			w.Header().Set("Link", `<https://api.github.com/search/issues?page=2>; rel="next"`)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(github.IssuesSearchResult{
			Issues: []*github.Issue{{Number: github.Int(requestCount)}},
		})
	})

	err := client.SearchIssuesStream("acme", func(page *github.IssuesSearchResult) error {
		callbackCount++
		if requestCount != callbackCount {
			t.Fatalf("callback ran after requesting a later page: requests=%d callbacks=%d", requestCount, callbackCount)
		}
		if len(page.Issues) != 1 || page.Issues[0].GetNumber() != callbackCount {
			t.Fatalf("unexpected page: %#v", page)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if requestCount != 2 || callbackCount != 2 {
		t.Fatalf("requests=%d callbacks=%d, want 2 each", requestCount, callbackCount)
	}
}

func TestParseGistSearchHTMLExtractsUniqueGists(t *testing.T) {
	html := `
<a href="/saumalya75/8ccbb1c6198b862db3ccbfac45efe27f">1 file</a>
<a href="/saumalya75/8ccbb1c6198b862db3ccbfac45efe27f/forks">forks</a>
<a href="/saumalya75/8ccbb1c6198b862db3ccbfac45efe27f/stargazers">stars</a>
<a href="/guichafy/85aa619319b59a436c679cc14fc7954f">other</a>
<a href="/search?q=aws">search</a>
`
	hits := parseGistSearchHTML(html)
	if len(hits) != 2 {
		t.Fatalf("hits = %#v", hits)
	}
	if hits[0].Owner != "saumalya75" || hits[0].ID != "8ccbb1c6198b862db3ccbfac45efe27f" {
		t.Fatalf("first hit = %#v", hits[0])
	}
	if hits[1].ID != "85aa619319b59a436c679cc14fc7954f" {
		t.Fatalf("second hit = %#v", hits[1])
	}
}

func TestGistMatchesKeywordIgnoresQualifiers(t *testing.T) {
	if !gistMatchesKeyword("config.env", "", "company token", "company -extension:md") {
		t.Fatal("expected keyword match after dropping gist qualifiers")
	}
	if gistMatchesKeyword("readme.md", "", "nothing here", "company") {
		t.Fatal("expected no match")
	}
}

func TestGistMatchesKeywordKeepsColonNeedles(t *testing.T) {
	if !gistMatchesKeyword("secrets.env", "", "aws_access_key_id: AKIA", "aws_access_key_id:") {
		t.Fatal("expected a colon-bearing secret token to remain a needle")
	}
	if gistMatchesKeyword("secrets.env", "", "unrelated", "aws_access_key_id:") {
		t.Fatal("expected colon-bearing needle to require a match")
	}
}

func TestGistMatchesKeywordSupportsQuotedPhrases(t *testing.T) {
	if !gistMatchesKeyword("secrets.env", "", "client_secret = token", `"client_secret = token"`) {
		t.Fatal("expected quoted phrase to match as one needle")
	}
	if gistMatchesKeyword("secrets.env", "", "client_secret = old token", `"client_secret = token"`) {
		t.Fatal("expected quoted phrase not to match separated text")
	}
}

func TestConvertGistToSearchResultsKeepsMatchingFiles(t *testing.T) {
	stubFilters(t, nil)
	gist := &github.Gist{
		ID:      github.String("8ccbb1c6198b862db3ccbfac45efe27f"),
		HTMLURL: github.String("https://gist.github.com/acme/8ccbb1c6198b862db3ccbfac45efe27f"),
		Owner:   &github.User{Login: github.String("acme")},
		Files: map[github.GistFilename]github.GistFile{
			"secrets.env": {Filename: github.String("secrets.env"), Content: github.String("password=acme-prod")},
			"readme.md":   {Filename: github.String("readme.md"), Content: github.String("unrelated")},
		},
	}
	results := convertGistToSearchResults(gist, gistHit{Owner: "acme", ID: gist.GetID()}, "acme")
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1, got %#v", len(results), results)
	}
	if results[0].Repo != "acme/8ccbb1c6198b862db3ccbfac45efe27f" {
		t.Fatalf("repo = %q", results[0].Repo)
	}
	if results[0].Path != "secrets.env" {
		t.Fatalf("path = %q", results[0].Path)
	}
}

func TestConvertGistToSearchResultsAppliesExtensionFilters(t *testing.T) {
	stubFilters(t, map[string][]model.Filter{
		"extension": {{FilterType: "blacklist", Content: "md"}},
	})
	gist := &github.Gist{
		ID:      github.String("8ccbb1c6198b862db3ccbfac45efe27f"),
		HTMLURL: github.String("https://gist.github.com/acme/8ccbb1c6198b862db3ccbfac45efe27f"),
		Owner:   &github.User{Login: github.String("acme")},
		Files: map[github.GistFilename]github.GistFile{
			"secrets.env": {Filename: github.String("secrets.env"), Content: github.String("acme token")},
			"notes.md":    {Filename: github.String("notes.md"), Content: github.String("acme token")},
		},
	}
	results := convertGistToSearchResults(gist, gistHit{Owner: "acme", ID: gist.GetID()}, "acme")
	if len(results) != 1 || results[0].Path != "secrets.env" {
		t.Fatalf("results = %#v", results)
	}
}

func TestSearchGistHitsUsesHTMLThenStopsOnRepeatPage(t *testing.T) {
	original := fetchGistSearchPage
	t.Cleanup(func() { fetchGistSearchPage = original })
	pages := 0
	fetchGistSearchPage = func(_ string, page int) (string, error) {
		pages++
		return `<a href="/acme/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa">gist</a>`, nil
	}

	client := &Client{Token: "test-token"}
	hits, err := client.SearchGistHits("acme")
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].ID != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("hits = %#v", hits)
	}
	if pages != 2 {
		t.Fatalf("expected a second page to detect duplicates, got %d fetches", pages)
	}
}

func TestSearchGistHitsEmptyHTMLDoesNotFallback(t *testing.T) {
	original := fetchGistSearchPage
	t.Cleanup(func() { fetchGistSearchPage = original })
	fetchGistSearchPage = func(string, int) (string, error) {
		return `<html><title>Search</title></html>`, nil
	}

	hits, err := (&Client{}).SearchGistHits("password")
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 0 {
		t.Fatalf("hits = %#v", hits)
	}
}

func TestSearchGistHitsHTTPErrorDoesNotFallback(t *testing.T) {
	original := fetchGistSearchPage
	t.Cleanup(func() { fetchGistSearchPage = original })
	fetchGistSearchPage = func(string, int) (string, error) {
		return "", errors.New("gist search returned HTTP 401")
	}

	hits, err := (&Client{}).SearchGistHits("password")
	if err == nil {
		t.Fatal("expected HTML search error")
	}
	if hits != nil {
		t.Fatalf("hits = %#v", hits)
	}
}

func TestGetGistSkipsNotFoundWithoutSleep(t *testing.T) {
	slept := stubSleep(t)
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "100")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
	})

	gist, err := client.getGist(context.Background(), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err == nil || gist != nil {
		t.Fatalf("expected not-found, got gist=%#v err=%v", gist, err)
	}
	if len(*slept) != 0 {
		t.Fatalf("getGist slept on 404: %v", *slept)
	}
}

func TestRunTaskSkipsWhenNoGithubSurfacesEnabled(t *testing.T) {
	original := getGithubRules
	getGithubRules = func(string) (error, []model.Rule) { return nil, nil }
	t.Cleanup(func() { getGithubRules = original })

	slept := stubSleep(t)
	outcome := RunTask()
	if len(*slept) != 0 {
		t.Fatalf("RunTask slept without rules: %v", *slept)
	}
	if outcome.Status != model.ScanStatusSkipped {
		t.Fatalf("status = %q, want skipped", outcome.Status)
	}
}

func TestRunTaskFailsWithoutTokenWhenIssueRulesExist(t *testing.T) {
	originalRules := getGithubRules
	getGithubRules = func(ruleType string) (error, []model.Rule) {
		if ruleType == "github_issue" {
			return nil, []model.Rule{{Content: "acme"}}
		}
		return nil, nil
	}
	originalTokens := listGithubTokens
	listGithubTokens = func(string) (error, []model.Token) { return nil, nil }
	t.Cleanup(func() {
		getGithubRules = originalRules
		listGithubTokens = originalTokens
	})

	outcome := RunTask()
	if outcome.Status != model.ScanStatusFailed {
		t.Fatalf("status = %q, want failed", outcome.Status)
	}
	if !strings.Contains(outcome.Message, "token") && !strings.Contains(outcome.Message, "client") {
		t.Fatalf("message = %q", outcome.Message)
	}
}

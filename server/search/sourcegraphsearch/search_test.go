package sourcegraphsearch

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"go.uber.org/zap"
)

func TestGlobalQueryIncludesExhaustiveRepositoryFilters(t *testing.T) {
	query := globalQuery("ghp_")
	for _, part := range []string{"ghp_", "count:all", "fork:yes", "archived:yes"} {
		if !strings.Contains(query, part) {
			t.Fatalf("query %q does not contain %q", query, part)
		}
	}
}

func TestParseStream(t *testing.T) {
	var gotName, gotData string
	err := parseStream(strings.NewReader("event: matches\ndata: [{\"type\":\"content\"}]\n\nevent: done\ndata: {}\n\n"), func(name, data string) error {
		if name == "matches" {
			gotName, gotData = name, data
		}
		return nil
	})
	if err != nil {
		t.Fatalf("parseStream returned error: %v", err)
	}
	if gotName != "matches" || gotData != "[{\"type\":\"content\"}]" {
		t.Fatalf("unexpected event: %q %q", gotName, gotData)
	}
}

func TestConvertMatch(t *testing.T) {
	global.GVA_LOG = zap.NewNop()
	result := convertMatch(streamMatch{
		Type:       "content",
		Repository: "github.com/acme/example",
		Commit:     "abc123",
		Path:       "config/app.env",
		LineMatches: []lineMatch{{
			Line:       "TOKEN=example",
			LineNumber: 7,
		}},
	})
	if result == nil {
		t.Fatal("expected a converted result")
	}
	if result.RepoUrl != "https://github.com/acme/example" {
		t.Fatalf("unexpected repo URL: %s", result.RepoUrl)
	}
	if !strings.Contains(result.Url, "abc123") || !strings.Contains(result.Url, "L=7") {
		t.Fatalf("unexpected result URL: %s", result.Url)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestSearchForSourcegraphParsesMatches(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("q") != "ghp_ count:all fork:yes archived:yes" {
			t.Fatalf("unexpected query: %s", req.URL.Query().Get("q"))
		}
		body := "event: matches\ndata: [{\"type\":\"content\",\"repository\":\"gitlab.com/acme/app\",\"path\":\".env\",\"commit\":\"deadbeef\",\"lineMatches\":[{\"line\":\"TOKEN=example\",\"lineNumber\":3}]}]\n\nevent: progress\ndata: {\"done\":true,\"matchCount\":1,\"repositoriesCount\":1}\n\nevent: done\ndata: {}\n\n"
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}
	results, warnings, err := SearchForSourcegraph(model.Rule{Content: "ghp_"}, client)
	if err != nil {
		t.Fatalf("SearchForSourcegraph returned error: %v", err)
	}
	if len(warnings) != 0 || len(results) != 1 {
		t.Fatalf("unexpected results: warnings=%v results=%d", warnings, len(results))
	}
	if results[0].Repo != "gitlab.com/acme/app" {
		t.Fatalf("unexpected repository: %s", results[0].Repo)
	}
}

func TestSearchForSourcegraphStreamInvokesCallbackPerMatchEvent(t *testing.T) {
	callbackCount := 0
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body := "event: matches\ndata: [{\"type\":\"content\",\"repository\":\"github.com/acme/one\",\"path\":\".env\",\"lineMatches\":[{\"line\":\"A=1\",\"lineNumber\":1}]}]\n\n" +
			"event: matches\ndata: [{\"type\":\"content\",\"repository\":\"github.com/acme/two\",\"path\":\".env\",\"lineMatches\":[{\"line\":\"B=2\",\"lineNumber\":2}]}]\n\n"
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}

	warnings, err := SearchForSourcegraphStream(model.Rule{Content: "ghp_"}, client, func(results []*model.SearchResult) error {
		callbackCount++
		if len(results) != 1 {
			t.Fatalf("results per event = %d, want 1", len(results))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("SearchForSourcegraphStream returned error: %v", err)
	}
	if len(warnings) != 0 || callbackCount != 2 {
		t.Fatalf("warnings=%v callbacks=%d, want no warnings and 2 callbacks", warnings, callbackCount)
	}
}

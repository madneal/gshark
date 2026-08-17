package service

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/madneal/gshark/config"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
)

func TestSearchResultContentPrefersTextMatches(t *testing.T) {
	previous := global.GVA_CONFIG
	t.Cleanup(func() { global.GVA_CONFIG = previous })

	fragment := "credential=secret-value"
	encoded, err := json.Marshal([]model.TextMatch{{Fragment: &fragment}})
	if err != nil {
		t.Fatal(err)
	}
	content := SearchResultContent(model.SearchResult{
		Matches:         "fallback",
		TextMatchesJson: encoded,
	})
	if content != "credential=secret-value" {
		t.Fatalf("content = %q", content)
	}
}

func TestSearchResultContentUsesBuiltinLimit(t *testing.T) {
	content := SearchResultContent(model.SearchResult{Matches: strings.Repeat("x", defaultAIAnalysisMaxContent+10)})
	if !strings.HasSuffix(content, "\n[truncated]") {
		t.Fatalf("content was not truncated")
	}
	if got := len([]rune(strings.TrimSuffix(content, "\n[truncated]"))); got != defaultAIAnalysisMaxContent {
		t.Fatalf("content length = %d, want %d", got, defaultAIAnalysisMaxContent)
	}
}

func TestParseSearchResultAnalysisAcceptsFencedJSON(t *testing.T) {
	body := []byte("{\"choices\":[{\"message\":{\"content\":\"```json\\n{\\\"real\\\":true,\\\"confidence\\\":0.92,\\\"reason\\\":\\\"usable credential\\\"}\\n```\"}}]}")
	analysis, err := parseSearchResultAnalysis(body)
	if err != nil {
		t.Fatal(err)
	}
	if !analysis.Real || analysis.Confidence != 0.92 {
		t.Fatalf("analysis = %#v", analysis)
	}
}

func TestAnalyzeSearchResultUsesOpenAICompatibleAPI(t *testing.T) {
	previous := global.GVA_CONFIG
	global.GVA_CONFIG.System = config.System{
		AiServer: "",
		AiToken:  "test-token",
		Model:    "test-model",
	}
	requests := 0
	server := newAIHTTPTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		var request ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.Model != "test-model" || len(request.Messages) != 2 {
			t.Fatalf("request = %#v", request)
		}
		if !strings.Contains(request.Messages[1].Content, "password=real") {
			t.Fatalf("missing evidence in request: %q", request.Messages[1].Content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"real\":true,\"confidence\":0.9,\"reason\":\"looks usable\"}"}}]}`))
	}))
	global.GVA_CONFIG.System.AiServer = server.URL
	t.Cleanup(func() { global.GVA_CONFIG = previous })

	result, err := AnalyzeSearchResult(model.SearchResult{Repo: "acme/app", Path: "config.env", Matches: "password=real"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Real || requests != 1 {
		t.Fatalf("result=%#v requests=%d", result, requests)
	}
}

func TestAnalyzeSearchResultFailsClosedOnMalformedResponse(t *testing.T) {
	previous := global.GVA_CONFIG
	global.GVA_CONFIG.System = config.System{AiServer: "", Model: "test-model"}
	server := newAIHTTPTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"yes"}}]}`))
	}))
	global.GVA_CONFIG.System.AiServer = server.URL
	t.Cleanup(func() { global.GVA_CONFIG = previous })

	if _, err := AnalyzeSearchResult(model.SearchResult{Matches: "placeholder"}); err == nil {
		t.Fatal("expected malformed verdict to fail closed")
	}
}

func TestAnalyzeSearchResultFallsBackToNextProvider(t *testing.T) {
	previous := global.GVA_CONFIG
	first := newAIHTTPTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("temporarily unavailable"))
	}))
	second := newAIHTTPTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"real\":true,\"confidence\":0.8,\"reason\":\"usable\"}"}}]}`))
	}))
	global.GVA_CONFIG.System = config.System{AiProviders: []config.AIProvider{
		{Name: "primary", Server: first.URL, Model: "first-model"},
		{Name: "backup", Server: second.URL, Model: "second-model"},
	}}
	t.Cleanup(func() { global.GVA_CONFIG = previous })

	result, err := AnalyzeSearchResult(model.SearchResult{Matches: "password=real"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Real {
		t.Fatalf("result = %#v", result)
	}
}

func TestTestAIConfigUsesSyntheticEvidence(t *testing.T) {
	server := newAIHTTPTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(request.Messages[1].Content, "real-secret") {
			t.Fatal("AI config test must not send a real secret")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"real\":false,\"confidence\":0.99,\"reason\":\"placeholder\"}"}}]}`))
	}))

	results, err := TestAIConfig(config.System{AiServer: server.URL, Model: "test-model", AiToken: "test-token"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || !results[0].Success {
		t.Fatalf("results = %#v", results)
	}
}

func newAIHTTPTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewUnstartedServer(handler)
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)
	return server
}

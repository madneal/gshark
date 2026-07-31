package codesearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	global.GVA_LOG = zap.NewNop()
	os.Exit(m.Run())
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonHTTPResponse(status int, value interface{}) *http.Response {
	var body bytes.Buffer
	_ = json.NewEncoder(&body).Encode(value)
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(&body),
		Header:     make(http.Header),
	}
}

func TestGetResultHandlesTransportError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network unavailable")
		}),
	}

	results, hasResult := GetResult(client, "https://searchcode.test/api/codesearch_I/?q=key&p=0")
	if hasResult {
		t.Fatalf("expected hasResult=false on transport error")
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestGetResultReturnsNoResultsOnNon2xx(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return jsonHTTPResponse(http.StatusNotFound, map[string]string{"error": "not found"}), nil
		}),
	}

	results, hasResult := GetResult(client, "https://searchcode.test/api/codesearch_I/?q=key&p=0")
	if hasResult {
		t.Fatalf("expected hasResult=false on non-2xx response")
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestGetResultStopsOnInvalidJSON(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("not json")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	results, hasResult := GetResult(client, "https://searchcode.test/api/codesearch_I/?q=key&p=0")
	if hasResult {
		t.Fatalf("expected hasResult=false on invalid JSON")
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestGetResultRespectsTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 20 * time.Millisecond}

	done := make(chan struct{})
	var hasResult bool
	go func() {
		_, hasResult = GetResult(client, server.URL)
		close(done)
	}()

	select {
	case <-done:
		if hasResult {
			t.Fatalf("expected hasResult=false when the client times out")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("GetResult did not return within the expected bound after client timeout")
	}
}

func TestSearchForSearchCodePaginatesUpToMax(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			res := model.SearchCodeRes{
				Results: []model.SearchCodeResult{
					{Repo: "example/repo", Filename: "file.go", Url: "https://example.test/file.go"},
				},
			}
			return jsonHTTPResponse(http.StatusOK, res), nil
		}),
	}

	rule := model.Rule{Content: "key"}
	results := SearchForSearchCode(rule, client)
	if len(results) != maxSearchcodePages {
		t.Fatalf("expected pagination to stop at %d pages, got %d results", maxSearchcodePages, len(results))
	}
}

package codesearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
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

func testRequest() codeSearchRequest {
	return codeSearchRequest{Repository: "https://github.com/example/repo", Query: "ghp_", MaxResults: 100}
}

func TestGetResultPostsNewSearchcodeRequest(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", req.Method)
			}
			if req.URL.Query().Get("client") != searchcodeClientName {
				t.Fatalf("client query = %q, want %q", req.URL.Query().Get("client"), searchcodeClientName)
			}
			if req.Header.Get("Content-Type") != "application/json" {
				t.Fatalf("content type = %q, want application/json", req.Header.Get("Content-Type"))
			}
			var got codeSearchRequest
			if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, testRequest()) {
				t.Fatalf("request body = %+v, want %+v", got, testRequest())
			}
			return jsonHTTPResponse(http.StatusOK, codeSearchResponse{Repository: got.Repository, Results: []codeSearchResult{}}), nil
		}),
	}

	result, err := GetResult(client, testRequest())
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || result.Repository != testRequest().Repository {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestGetResultHandlesTransportError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network unavailable")
		}),
	}

	result, err := GetResult(client, testRequest())
	if err == nil {
		t.Fatal("expected a transport error")
	}
	if result != nil {
		t.Fatalf("expected no response, got %+v", result)
	}
}

func TestGetResultReturnsAPIError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return jsonHTTPResponse(http.StatusNotFound, map[string]interface{}{
				"error": map[string]string{"code": "repository_not_found", "message": "repository not found"},
			}), nil
		}),
	}

	_, err := GetResult(client, testRequest())
	if err == nil {
		t.Fatal("expected an API error")
	}
	if !strings.Contains(err.Error(), "repository_not_found") {
		t.Fatalf("error = %v, want API error code", err)
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

	_, err := GetResult(client, testRequest())
	if err == nil {
		t.Fatal("expected an invalid JSON error")
	}
}

func TestGetResultRespectsTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	previousBaseURL := searchcodeBaseURL
	searchcodeBaseURL = server.URL
	defer func() { searchcodeBaseURL = previousBaseURL }()

	client := &http.Client{Timeout: 20 * time.Millisecond}
	done := make(chan struct{})
	var err error
	go func() {
		_, err = GetResult(client, testRequest())
		close(done)
	}()

	select {
	case <-done:
		if err == nil {
			t.Fatal("expected timeout error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("GetResult did not return within the expected bound after client timeout")
	}
}

func TestSearchForSearchCodePaginatesAndConvertsResults(t *testing.T) {
	var offsets []int
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			var request codeSearchRequest
			if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
				t.Fatal(err)
			}
			offsets = append(offsets, request.Offset)
			response := codeSearchResponse{
				Repository: request.Repository,
				CommitSHA:  "abc123",
				HasMore:    request.Offset < 2,
				Results: []codeSearchResult{{
					File: "config.yaml",
					Matches: []codeSearchMatch{{
						Line:          8,
						Content:       "token: ghp_example",
						ContextBefore: []string{"auth:"},
						ContextAfter:  []string{"enabled: true"},
					}},
				}},
			}
			return jsonHTTPResponse(http.StatusOK, response), nil
		}),
	}

	results, err := SearchForSearchCode(model.Rule{Content: "ghp_"}, testRequest().Repository, client)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(offsets, []int{0, 1, 2}) {
		t.Fatalf("offsets = %v, want [0 1 2]", offsets)
	}
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	if results[0].Url != "https://github.com/example/repo/blob/abc123/config.yaml" {
		t.Fatalf("result URL = %q", results[0].Url)
	}
	if results[0].Repo != "example/repo" || results[0].Keyword != "ghp_" {
		t.Fatalf("result metadata = %+v", results[0])
	}
}

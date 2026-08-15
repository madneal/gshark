package postman

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/model"
	"go.uber.org/zap"
)

func TestRunTaskWithoutRulesDoesNotSleepOrSearch(t *testing.T) {
	originalLog := global.GVA_LOG
	originalRules := getPostmanRules
	originalSearch := searchPostman
	global.GVA_LOG = zap.NewNop()
	getPostmanRules = func(string) (error, []model.Rule) { return nil, nil }
	searchPostman = func(*[]model.Rule) error {
		t.Fatal("RunTask searched without Postman rules")
		return nil
	}
	t.Cleanup(func() {
		global.GVA_LOG = originalLog
		getPostmanRules = originalRules
		searchPostman = originalSearch
	})

	outcome := RunTask()
	if outcome.Status != model.ScanStatusSkipped {
		t.Fatalf("outcome status = %q, want skipped", outcome.Status)
	}
}

func TestRunTask(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("RUN_POSTMAN_INTEGRATION") != "1" {
		t.Skip("set RUN_POSTMAN_INTEGRATION=1 to run the Postman integration test")
	}
	InitialDataBase()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	RunTask()
}

func TestSearchAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if os.Getenv("RUN_POSTMAN_INTEGRATION") != "1" {
		t.Skip("set RUN_POSTMAN_INTEGRATION=1 to run the Postman integration test")
	}
	res, err := SearchAPI("mihoyo", "collection")
	if err != nil {
		t.Fatal(err)
	}
	if len(*res) == 0 {
		t.Fatal("expected at least one response page")
	}
}

func TestSearchAPIHandlesTransportError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network unavailable")
		}),
	}

	res, err := searchAPI(client, "https://postman.test/search", "OPENAI_API_KEY", "request")
	if err == nil || !strings.Contains(err.Error(), "network unavailable") {
		t.Fatalf("expected transport error, got %v", err)
	}
	if len(*res) != 0 {
		t.Fatalf("expected no response pages, got %d", len(*res))
	}
}

func TestSearchAPIPaginatesWithOffsetsAndEscapesRule(t *testing.T) {
	var cursors []string
	const rule = "key \"quoted\" \\\\ value"
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var request postmanSearchRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, fmt.Errorf("decode request: %w", err)
		}
		if request.QueryText != rule {
			t.Errorf("query text = %q, want %q", request.QueryText, rule)
		}
		if request.ElementType != "requests" {
			t.Errorf("element type = %q, want requests", request.ElementType)
		}
		if r.URL.Query().Get("limit") != "25" {
			t.Errorf("limit = %q, want 25", r.URL.Query().Get("limit"))
		}
		cursors = append(cursors, r.URL.Query().Get("cursor"))

		count := postmanPageSize
		nextCursor := "second-page"
		start := 0
		if r.URL.Query().Get("cursor") == "second-page" {
			count = 2
			nextCursor = ""
			start = postmanPageSize - 1
		}
		data := make([]map[string]interface{}, count)
		for i := range data {
			data[i] = map[string]interface{}{
				"id":   fmt.Sprintf("request-%d", start+i),
				"name": fmt.Sprintf("Request %d", start+i),
			}
		}
		return jsonHTTPResponse(http.StatusOK, map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"total":      count,
				"nextCursor": nextCursor,
			},
		}), nil
	})}

	res, err := searchAPI(client, "https://postman.test/search", rule, "request")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cursors, []string{"", "second-page"}) {
		t.Fatalf("cursors = %v, want [\"\" \"second-page\"]", cursors)
	}
	if len(*res) != 2 || len((*res)[0].Data)+len((*res)[1].Data) != 26 {
		t.Fatalf("unexpected response pages: %#v", res)
	}
}

func TestSearchAPIStreamInvokesCallbackPerPage(t *testing.T) {
	requestCount := 0
	var pageSizes []int
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		requestCount++
		nextCursor := ""
		data := []map[string]interface{}{{"id": fmt.Sprintf("request-%d", requestCount)}}
		if requestCount == 1 {
			nextCursor = "second-page"
		}
		return jsonHTTPResponse(http.StatusOK, map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{"nextCursor": nextCursor},
		}), nil
	})}

	err := searchAPIStream(client, "https://postman.test/search", "key", "request", func(page PostmanRes) error {
		pageSizes = append(pageSizes, len(page.Data))
		if requestCount != len(pageSizes) {
			t.Fatalf("callback ran after requesting a later page: requests=%d callbacks=%d", requestCount, len(pageSizes))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(pageSizes, []int{1, 1}) {
		t.Fatalf("page sizes = %v, want [1 1]", pageSizes)
	}
}

func TestSearchAPIReturnsPartialPagesOnLaterFailure(t *testing.T) {
	requestCount := 0
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		requestCount++
		if requestCount == 2 {
			return jsonHTTPResponse(http.StatusTooManyRequests, map[string]string{"error": "rate limited"}), nil
		}
		data := make([]map[string]interface{}, postmanPageSize)
		for i := range data {
			data[i] = map[string]interface{}{
				"id": fmt.Sprintf("request-%d", i),
			}
		}
		return jsonHTTPResponse(http.StatusOK, map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"total":      postmanPageSize,
				"nextCursor": "second-page",
			},
		}), nil
	})}

	res, err := searchAPI(client, "https://postman.test/search", "key", "request")
	if err == nil || !strings.Contains(err.Error(), "429") {
		t.Fatalf("expected HTTP 429 error, got %v", err)
	}
	if len(*res) != 1 {
		t.Fatalf("expected one successful page, got %d", len(*res))
	}
}

func TestConvertToSearchResultUsesCurrentResponseShape(t *testing.T) {
	var response PostmanRes
	err := json.Unmarshal([]byte(`{
		"data": [{
			"id": "request-id",
			"method": "POST",
			"name": "Chat Completions",
			"url": "https://api.openai.com/v1/chat/completions",
			"description": "Create a chat completion",
			"collection": {
				"id": "collection-id",
				"name": "OpenAI API"
			},
			"organization": {
				"id": "organization-id",
				"name": "Postman DevRel"
			},
			"links": {
				"web": {
					"href": "https://go.postman.co/request/request-id"
				}
			}
		}]
	}`), &response)
	if err != nil {
		t.Fatal(err)
	}

	results := response.ConvertToSearchResult("OPENAI_API_KEY")
	if len(*results) != 1 {
		t.Fatalf("expected one result, got %d", len(*results))
	}
	result := (*results)[0]
	if result.Url != "https://go.postman.co/request/request-id" {
		t.Fatalf("URL = %q", result.Url)
	}
	if result.RepoUrl != result.Url {
		t.Fatalf("repository URL = %q, want %q", result.RepoUrl, result.Url)
	}
	if result.Path != "request-id" {
		t.Fatalf("path = %q, want request-id", result.Path)
	}
	if result.Matches != "POST | Chat Completions | https://api.openai.com/v1/chat/completions | Create a chat completion" {
		t.Fatalf("matches = %q", result.Matches)
	}
	if result.Repo != "Postman DevRel | OpenAI API | Chat Completions [request-id]" {
		t.Fatalf("repo = %q", result.Repo)
	}
}

func TestBuildPostmanRepoPreservesLongNameAndUniqueIDSuffix(t *testing.T) {
	resource := PostmanResource{
		ID:   "request-id",
		Name: strings.Repeat("长", 250),
	}
	resource.Organization.Name = "Postman"
	resource.Collection.Name = "Collection"

	repo := buildPostmanRepo(resource)
	want := "Postman | Collection | " + resource.Name + " [request-id]"
	if repo != want {
		t.Fatalf("repo = %q, want %q", repo, want)
	}
	if len([]rune(repo)) <= 200 {
		t.Fatalf("test repo length = %d, want more than the old column limit", len([]rune(repo)))
	}
	if !strings.HasSuffix(repo, " [request-id]") {
		t.Fatalf("repo does not preserve ID suffix: %q", repo)
	}
}

func TestBuildPostmanRepoFallsBackToWorkspaceNameWhenOrganizationIsBlank(t *testing.T) {
	resource := PostmanResource{
		ID:   "request-id",
		Name: "Launcher HTML Widget Content",
	}
	resource.Collection.Name = "Genshin Launcher"
	resource.Workspace.Name = "MiHoYo"

	repo := buildPostmanRepo(resource)
	want := "MiHoYo | Genshin Launcher | Launcher HTML Widget Content [request-id]"
	if repo != want {
		t.Fatalf("repo = %q, want %q", repo, want)
	}
}

func TestSearchAPIRejectsRepeatedCursor(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return jsonHTTPResponse(http.StatusOK, map[string]interface{}{
			"data": []map[string]interface{}{{"id": "request-id"}},
			"meta": map[string]interface{}{"nextCursor": "same-cursor"},
		}), nil
	})}

	res, err := searchAPI(client, "https://postman.test/search", "key", "request")
	if err == nil || !strings.Contains(err.Error(), "repeated cursor") {
		t.Fatalf("expected repeated cursor error, got %v", err)
	}
	if len(*res) != 1 {
		t.Fatalf("expected one unique page, got %d", len(*res))
	}
}

func InitialDataBase() {
	global.GVA_VP = initialize.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = initialize.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()                                                 // gorm连接数据库
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

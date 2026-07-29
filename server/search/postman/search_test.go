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
)

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
	var offsets []int
	const rule = "key \"quoted\" \\\\ value"
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var request postmanSearchRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, fmt.Errorf("decode request: %w", err)
		}
		if request.Body.QueryText != rule {
			t.Errorf("query text = %q, want %q", request.Body.QueryText, rule)
		}
		offsets = append(offsets, request.Body.From)

		count := postmanPageSize
		if request.Body.From == postmanPageSize {
			count = 5
		}
		data := make([]map[string]interface{}, count)
		for i := range data {
			data[i] = map[string]interface{}{
				"document": map[string]interface{}{
					"id":           fmt.Sprintf("request-%d", request.Body.From+i),
					"documentType": "request",
				},
			}
		}
		return jsonHTTPResponse(http.StatusOK, map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"total": map[string]interface{}{"request": 30},
			},
		}), nil
	})}

	res, err := searchAPI(client, "https://postman.test/search", rule, "request")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(offsets, []int{0, 25}) {
		t.Fatalf("offsets = %v, want [0 25]", offsets)
	}
	if len(*res) != 2 || len((*res)[0].Data)+len((*res)[1].Data) != 30 {
		t.Fatalf("unexpected response pages: %#v", res)
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
				"document": map[string]interface{}{"id": fmt.Sprintf("request-%d", i)},
			}
		}
		return jsonHTTPResponse(http.StatusOK, map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"total": map[string]interface{}{"request": 30},
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
			"document": {
				"id": "request-id",
				"documentType": "request",
				"method": "POST",
				"name": "Chat Completions",
				"url": "{{baseUrl}}/chat/completions",
				"publisherName": "Postman DevRel",
				"publisherHandle": "devrel",
				"workspaceSlug": "openai"
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
	if result.Url != "https://www.postman.com/devrel/workspace/openai/request/request-id" {
		t.Fatalf("URL = %q", result.Url)
	}
	if result.Matches != "POST | Chat Completions | {{baseUrl}}/chat/completions" {
		t.Fatalf("matches = %q", result.Matches)
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

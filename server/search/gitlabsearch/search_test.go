package gitlabsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/model"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	global.GVA_LOG = zap.NewNop()
	os.Exit(m.Run())
}

func TestGetClient(t *testing.T) {
	InitialDataBase()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	client := GetClient()
	if client == nil {
		t.Skip("gitlab client not available (no token)")
	}
	assert.Equal(t, true, client != nil, "the client is not nil")
}

func InitialDataBase() {
	global.GVA_VP = initialize.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	global.GVA_LOG = initialize.Zap()                                                 // 初始化zap日志库
	global.GVA_DB = initialize.Gorm()                                                 // gorm连接数据库
}

func TestGetProjects(t *testing.T) {
	InitialDataBase()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	client := GetClient()
	if client == nil {
		t.Skip("gitlab client not available (no token)")
	}
	GetProjects(client, 1)
}

func TestListValidProjects(t *testing.T) {
	InitialDataBase()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	client := GetClient()
	if client == nil {
		t.Skip("gitlab client not available (no token)")
	}
	projects := ListValidProjects()
	assert.Equal(t, true, len(*projects) > 0, "there is should one more project")
}

func TestGetProjectById(t *testing.T) {
	InitialDataBase()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	client := GetClient()
	if client == nil {
		t.Skip("gitlab client not available (no token)")
	}
	GetProjectById(client, 32123952)
}

func TestSearchBlobs(t *testing.T) {
	InitialDataBase()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	client := GetClient()
	if client == nil {
		t.Skip("gitlab client not available (no token)")
	}
	blobs, _, _ := SearchBlobs(client, "mihoyo")
	fmt.Println(blobs)
}

func TestNextProjectBatch(t *testing.T) {
	InitialDataBase()
	if global.GVA_DB == nil {
		t.Skip("database not available")
	}
	batch := NextProjectBatch(1)
	assert.Equal(t, true, len(batch) <= 1, "batch should be truncated to the requested size")
}

func newTestGitlabClient(t *testing.T, handler http.HandlerFunc) *gitlab.Client {
	t.Helper()
	// gitlab.Client makes one throwaway GET to the bare base URL the first
	// time it's used, to read rate-limit headers and configure its limiter
	// (see (*gitlab.Client).configureLimiter). Absorb that here so callers'
	// handlers only see the API calls they actually care about.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v4/" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)

	client, err := gitlab.NewClient("test-token", gitlab.WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func writeGitlabJSON(w http.ResponseWriter, status int, value interface{}) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func TestIsGlobalSearchUnsupported(t *testing.T) {
	if !isGlobalSearchUnsupported(&gitlab.Response{Response: &http.Response{StatusCode: http.StatusBadRequest}}) {
		t.Fatal("expected 400 response to be reported as unsupported")
	}
	if isGlobalSearchUnsupported(&gitlab.Response{Response: &http.Response{StatusCode: http.StatusInternalServerError}}) {
		t.Fatal("expected 500 response to not be reported as unsupported")
	}
	if isGlobalSearchUnsupported(nil) {
		t.Fatal("expected a nil response to not be reported as unsupported")
	}
}

func TestSearchBlobsReturnsNotOkOnBadRequest(t *testing.T) {
	client := newTestGitlabClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeGitlabJSON(w, http.StatusBadRequest, map[string]string{"message": "Scope not supported without Elasticsearch!"})
	})

	blobs, resp, ok := SearchBlobs(client, "test-query")
	if ok {
		t.Fatal("expected ok to be false on a 400 response")
	}
	if len(blobs) != 0 {
		t.Fatalf("expected no blobs, got %d", len(blobs))
	}
	if resp == nil || resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected the 400 response to be returned, got %#v", resp)
	}
}

func TestSearchBlobsPaginatesUntilNextPageIsZero(t *testing.T) {
	requestCount := 0
	client := newTestGitlabClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			w.Header().Set("X-Next-Page", "2")
			writeGitlabJSON(w, http.StatusOK, []gitlab.Blob{{Basename: "first"}})
			return
		}
		writeGitlabJSON(w, http.StatusOK, []gitlab.Blob{{Basename: "second"}})
	})

	blobs, resp, ok := SearchBlobs(client, "test-query")
	if !ok {
		t.Fatal("expected ok to be true")
	}
	if resp != nil {
		t.Fatalf("expected a nil response after a fully successful paginated search, got %#v", resp)
	}
	if len(blobs) != 2 {
		t.Fatalf("expected blobs from both pages, got %d", len(blobs))
	}
}

func TestRunGlobalSearchTaskFallsBackWhenUnsupported(t *testing.T) {
	requestCount := 0
	client := newTestGitlabClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		writeGitlabJSON(w, http.StatusBadRequest, map[string]string{"message": "Scope not supported without Elasticsearch!"})
	})

	rules := []model.Rule{{Content: "first"}, {Content: "second"}}
	if RunGlobalSearchTask(client, rules) {
		t.Fatal("expected RunGlobalSearchTask to report unavailable global search")
	}
	if requestCount != 1 {
		t.Fatalf("expected to stop after the first rule's 400 response, got %d requests", requestCount)
	}
}

func TestSearchCodePaginatesAcrossPages(t *testing.T) {
	requestCount := 0
	client := newTestGitlabClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount == 1 {
			w.Header().Set("X-Next-Page", "2")
			writeGitlabJSON(w, http.StatusOK, []gitlab.Blob{{Filename: "a.go", Basename: "a.go"}})
			return
		}
		writeGitlabJSON(w, http.StatusOK, []gitlab.Blob{{Filename: "b.go", Basename: "b.go"}})
	})

	project := model.Repo{ProjectId: 1, Url: "https://gitlab.test/group/project"}
	results := SearchCode("test-query", project, client)
	if requestCount != 2 {
		t.Fatalf("expected two paginated requests, got %d", requestCount)
	}
	if len(results) != 2 {
		t.Fatalf("expected results from both pages, got %d", len(results))
	}
}

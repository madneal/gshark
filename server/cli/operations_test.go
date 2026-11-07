package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOperations(t *testing.T) {
	for _, tc := range []struct {
		args         []string
		method, path string
		check        func(*testing.T, map[string]interface{})
	}{
		{[]string{"results", "set-status", "1", "2", "1", "--status", "ignored"}, "POST", "/searchResult/updateSearchResultStatusByIds", func(t *testing.T, p map[string]interface{}) {
			if p["status"] != float64(2) || len(p["ids"].([]interface{})) != 2 {
				t.Fatalf("%v", p)
			}
		}},
		{[]string{"results", "delete", "1", "2"}, "DELETE", "/searchResult/deleteSearchResultByIds", nil},
		{[]string{"rules", "create", "--type", "github", "--name", "test", "--content", "secret", "--enabled=false"}, "POST", "/rule/createRule", func(t *testing.T, p map[string]interface{}) {
			if p["status"] != false || p["ruleType"] != "github" {
				t.Fatalf("%v", p)
			}
		}},
		{[]string{"rules", "update", "7", "--description", ""}, "PUT", "/rule/updateRule", func(t *testing.T, p map[string]interface{}) {
			if p["desc"] != "" || p["content"] != "keep-me" || p["status"] != false || p["ID"] != float64(7) {
				t.Fatalf("update lost fields: %v", p)
			}
		}},
		{[]string{"rules", "enable", "7"}, "POST", "/rule/switchRuleStatus", func(t *testing.T, p map[string]interface{}) {
			if p["status"] != float64(1) {
				t.Fatal(p)
			}
		}},
		{[]string{"rules", "disable", "7"}, "POST", "/rule/switchRuleStatus", func(t *testing.T, p map[string]interface{}) {
			if p["status"] != float64(0) {
				t.Fatal(p)
			}
		}},
		{[]string{"rules", "delete", "7"}, "DELETE", "/rule/deleteRuleByIds", nil},
		{[]string{"repos", "create", "--type", "gitlab", "--url", "https://example.com/repo", "--project-id", "42"}, "POST", "/repo/createRepo", func(t *testing.T, p map[string]interface{}) {
			if p["projectId"] != float64(42) {
				t.Fatal(p)
			}
		}},
		{[]string{"repos", "update", "7", "--status", "0"}, "PUT", "/repo/updateRepo", func(t *testing.T, p map[string]interface{}) {
			if p["status"] != float64(0) || p["url"] != "https://keep.example" {
				t.Fatal(p)
			}
		}},
		{[]string{"repos", "delete", "7"}, "DELETE", "/repo/deleteRepoByIds", nil},
		{[]string{"filters", "create", "--type", "blacklist", "--class", "extension", "--content", ".md"}, "POST", "/filter/createFilter", nil},
		{[]string{"filters", "update", "7", "--content", ".txt"}, "PUT", "/filter/updateFilter", func(t *testing.T, p map[string]interface{}) {
			if p["filter_class"] != "extension" || p["content"] != ".txt" {
				t.Fatal(p)
			}
		}},
		{[]string{"filters", "delete", "7"}, "DELETE", "/filter/deleteFilterByIds", nil},
	} {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			writes := 0
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("x-token") != "test-token" {
					t.Error("missing credential")
				}
				if r.Method == "GET" {
					if r.URL.Query().Get("ID") != "7" {
						t.Error("wrong read ID")
					}
					var name string
					switch {
					case strings.HasPrefix(r.URL.Path, "/rule/"):
						name = "rule"
					case strings.HasPrefix(r.URL.Path, "/repo/"):
						name = "repo"
					default:
						name = "filter"
					}
					fmt.Fprintf(w, `{"code":0,"data":{"%s":{"ID":7,"content":"keep-me","status":false,"desc":"previous","url":"https://keep.example","filter_class":"extension","filter_type":"blacklist"}}}`, name)
					return
				}
				writes++
				if r.Method != tc.method || r.URL.Path != tc.path {
					t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
				}
				var p map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					t.Fatal(err)
				}
				if tc.check != nil {
					tc.check(t, p)
				}
				w.Write([]byte(`{"code":0,"data":{},"msg":"ok"}`))
			}))
			defer s.Close()
			t.Setenv("GSHARK_TOKEN", "test-token")
			cmd := NewCommand()
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetArgs(append([]string{"--server", s.URL}, tc.args...))
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if writes != 1 || !json.Valid(out.Bytes()) {
				t.Fatalf("writes=%d output=%s", writes, out.String())
			}
		})
	}
}

func TestOperationValidationBeforeRequests(t *testing.T) {
	called := false
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	defer s.Close()
	t.Setenv("GSHARK_TOKEN", "test-token")
	for _, args := range [][]string{
		{"results", "delete", "0"}, {"results", "set-status", "1", "--status", "bad"},
		{"rules", "create", "--name", "only-name"}, {"rules", "update", "1"},
		{"filters", "create", "--type", "bad", "--class", "extension", "--content", "x"},
		{"repos", "update", "1", "--project-id", "-1"}, {"rules", "enable", "-5"},
		{"--output", "bad", "rules", "disable", "1"},
	} {
		cmd := NewCommand()
		cmd.SetArgs(append([]string{"--server", s.URL}, args...))
		if err := cmd.Execute(); err == nil {
			t.Fatalf("expected failure: %v", args)
		}
	}
	if called {
		t.Fatal("invalid input reached server")
	}
}

func TestExportPaginationAndCleanup(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprint(fail), func(t *testing.T) {
			pages := 0
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				pages++
				if r.URL.Query().Get("keyword") != "test" || r.URL.Query().Get("status") != "0" {
					t.Error("lost filters")
				}
				if pages == 2 && fail {
					w.Write([]byte(`{"code":7,"data":{},"msg":"expired"}`))
					return
				}
				n := 100
				if pages == 2 {
					n = 1
				}
				rows := make([]map[string]interface{}, n)
				for i := range rows {
					rows[i] = map[string]interface{}{"ID": (pages-1)*100 + i + 1, "matches": "super-secret"}
				}
				json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "data": map[string]interface{}{"list": rows, "total": 101}})
			}))
			defer s.Close()
			path := filepath.Join(t.TempDir(), "results.jsonl")
			t.Setenv("GSHARK_TOKEN", "test-token")
			cmd := NewCommand()
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetArgs([]string{"--server", s.URL, "results", "export", "--file", path, "--keyword", "test", "--status", "0"})
			err := cmd.Execute()
			if fail {
				if err == nil {
					t.Fatal("expected error")
				}
				if _, err = os.Stat(path); !os.IsNotExist(err) {
					t.Fatal("partial export left behind")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if pages != 2 || bytes.Count(data, []byte("\n")) != 101 || bytes.Contains(data, []byte("super-secret")) {
				t.Fatal("pagination or redaction failed")
			}
			if !strings.Contains(out.String(), `"exported":101`) {
				t.Fatal(out.String())
			}
			overwrite := NewCommand()
			overwrite.SetArgs([]string{"--server", s.URL, "results", "export", "--file", path})
			if err = overwrite.Execute(); err == nil {
				t.Fatal("overwrote existing file")
			}
		})
	}
}

func TestCSVExportLimitAndSensitiveOptIn(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"code":0,"data":{"list":[{"ID":1,"repo":"=CMD()","matches":"secret"},{"ID":2,"matches":"other"}],"total":2}}`))
	}))
	defer s.Close()
	t.Setenv("GSHARK_TOKEN", "test-token")
	path := filepath.Join(t.TempDir(), "results.csv")
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--server", s.URL, "results", "export", "--file", path, "--format", "csv", "--include-sensitive", "--max-results", "1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "'=CMD()") || !strings.Contains(string(data), "secret") || strings.Contains(string(data), "other") || !strings.Contains(out.String(), `"truncated":true`) {
		t.Fatalf("csv=%s summary=%s", data, out.String())
	}
}

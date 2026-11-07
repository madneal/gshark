package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResultCommands(t *testing.T) {
	for _, args := range [][]string{{"results", "list", "--keyword", "a & b"}, {"results", "get", "123"}, {"results", "get", "123", "--include-sensitive"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if args[1] == "list" && (r.URL.Query().Get("status") != "-1" || r.URL.Query().Get("keyword") != "a & b" || r.URL.Query().Get("pageSize") != "20") {
					t.Errorf("query: %v", r.URL.Query())
				}
				if args[1] == "get" && r.URL.Query().Get("ID") != "123" {
					t.Error("missing ID")
				}
				w.Write([]byte(`{"code":0,"data":{"list":[{"ID":123,"matches":"secret-credential","text_matches":[{"fragment":"secret-credential"}],"keyword":"secret-credential","repo":"secret-credential","url":"secret-credential"}],"total":1}}`))
			}))
			defer s.Close()
			t.Setenv("GSHARK_TOKEN", "test-token")
			cmd := NewCommand()
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetArgs(append([]string{"--server", s.URL}, args...))
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if !json.Valid(out.Bytes()) {
				t.Fatalf("not JSON: %s", out.String())
			}
			raw := args[len(args)-1] == "--include-sensitive"
			if strings.Contains(out.String(), "secret-credential") != raw {
				t.Fatalf("redaction mismatch: %s", out.String())
			}
		})
	}
}

func TestValidationWithoutServer(t *testing.T) {
	for _, args := range [][]string{{"results", "list", "--limit", "101"}, {"results", "get", "0"}, {"rules", "list", "--page", "0"}, {"scans", "overview", "--output", "xml"}, {"results", "list", "--status", "5"}} {
		cmd := NewCommand()
		cmd.SetArgs(args)
		if err := cmd.Execute(); err == nil {
			t.Fatalf("expected error for %v", args)
		}
	}
	cmd := NewCommand()
	cmd.SetArgs([]string{"--help"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil || !strings.Contains(out.String(), "results") {
		t.Fatalf("help failed: %v", err)
	}
}

func TestRuleAndScanRoutes(t *testing.T) {
	for _, tc := range []struct {
		args             []string
		path, key, value string
	}{
		{[]string{"rules", "list", "--type", "github"}, "/rule/getRuleList", "ruleType", "github"},
		{[]string{"scans", "logs", "--cycle-id", "cycle-1"}, "/scanLog/getScanLogList", "cycleId", "cycle-1"},
		{[]string{"scans", "overview"}, "/scanLog/getScanLogOverview", "", ""},
	} {
		t.Run(tc.path, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path || (tc.key != "" && r.URL.Query().Get(tc.key) != tc.value) {
					t.Errorf("unexpected request %s", r.URL)
				}
				w.Write([]byte(`{"code":0,"data":{"total":0,"list":[]}}`))
			}))
			defer s.Close()
			t.Setenv("GSHARK_TOKEN", "test-token")
			cmd := NewCommand()
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetArgs(append([]string{"--server", s.URL, "--output", "text"}, tc.args...))
			if err := cmd.Execute(); err != nil || !strings.Contains(out.String(), "Total: 0") {
				t.Fatalf("%s %v", out.String(), err)
			}
		})
	}
}

func TestAPIErrorDoesNotWriteSuccessOutput(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"code":7,"data":{},"msg":"permission denied"}`))
	}))
	defer s.Close()
	t.Setenv("GSHARK_TOKEN", "test-token")
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--server", s.URL, "scans", "overview"})
	if err := cmd.Execute(); err == nil || out.Len() != 0 {
		t.Fatalf("error=%v output=%s", err, out.String())
	}
}

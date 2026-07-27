package search

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// repoRoot walks up from this test file to the monorepo root (contains server/ and web/).
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file) // server/search
	root := filepath.Clean(filepath.Join(dir, "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "server", "main.go")); err != nil {
		t.Fatalf("expected monorepo root at %s: %v", root, err)
	}
	return root
}

func readRepoFile(t *testing.T, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(b)
}

// TestArchitectureServeScanCLISplit asserts the shipped CLI keeps HTTP management
// and multi-platform scanning as separate entrypoints (not merged into serve).
func TestArchitectureServeScanCLISplit(t *testing.T) {
	mainGo := readRepoFile(t, "server/main.go")
	if !strings.Contains(mainGo, `Use:   "serve"`) {
		t.Fatal("main.go missing serve subcommand")
	}
	if !strings.Contains(mainGo, `Use:   "scan"`) {
		t.Fatal("main.go missing scan subcommand")
	}
	if !strings.Contains(mainGo, "core.RunServer()") {
		t.Fatal("serve must call core.RunServer")
	}
	if !strings.Contains(mainGo, "search.ScanTask()") {
		t.Fatal("scan must call search.ScanTask")
	}
	// serve must not start ScanTask
	serveBlockStart := strings.Index(mainGo, `Use:   "serve"`)
	scanBlockStart := strings.Index(mainGo, `Use:   "scan"`)
	if serveBlockStart < 0 || scanBlockStart < 0 || scanBlockStart <= serveBlockStart {
		t.Fatal("cannot locate serve/scan command blocks")
	}
	serveBlock := mainGo[serveBlockStart:scanBlockStart]
	if strings.Contains(serveBlock, "ScanTask") {
		t.Fatal("serve command must not invoke ScanTask")
	}
}

// TestArchitectureScanTaskPlatforms documents the multi-platform loop body.
func TestArchitectureScanTaskPlatforms(t *testing.T) {
	scanGo := readRepoFile(t, "server/search/scan.go")
	for _, name := range []string{
		"gitlabsearch.RunTask",
		"codesearch.RunTask",
		"githubsearch.RunTask",
		"gobuster.RunTask",
		"postman.RunTask",
	} {
		if !strings.Contains(scanGo, name) {
			t.Fatalf("ScanTask loop missing %s", name)
		}
	}
	if !strings.Contains(scanGo, "GVA_DB == nil") {
		t.Fatal("ScanTask should bail when DB is not initialized")
	}
}

// TestArchitectureDockerEntryPoints asserts dual-image ENTRYPOINTs.
func TestArchitectureDockerEntryPoints(t *testing.T) {
	serveDF := readRepoFile(t, "server/deploy/serve/Dockerfile")
	scanDF := readRepoFile(t, "server/deploy/scan/Dockerfile")
	if !strings.Contains(serveDF, "./gshark serve") {
		t.Fatal("serve Dockerfile must ENTRYPOINT gshark serve")
	}
	if !strings.Contains(scanDF, "./gshark scan") {
		t.Fatal("scan Dockerfile must ENTRYPOINT gshark scan")
	}
	compose := readRepoFile(t, "docker-compose.yaml")
	if !strings.Contains(compose, "container_name: gshark-server") {
		t.Fatal("compose missing gshark-server")
	}
	if !strings.Contains(compose, "container_name: gshark-scanner") {
		t.Fatal("compose missing gshark-scanner")
	}
	if !strings.Contains(compose, "deploy/scan/Dockerfile") {
		t.Fatal("compose scan service should build scan Dockerfile")
	}
	quick := readRepoFile(t, "scripts/quick-docker.sh")
	if !strings.Contains(quick, "mysql server web") {
		t.Fatal("quick-docker default must start mysql server web only")
	}
	if !strings.Contains(quick, "--with-scan") {
		t.Fatal("quick-docker must document/with-scan for scanner")
	}
}

// TestArchitectureWebTasksAreOnServeAPI asserts result-page tasks hit searchResult
// routes on the API server, not the scan process.
func TestArchitectureWebTasksAreOnServeAPI(t *testing.T) {
	router := readRepoFile(t, "server/router/search_result.go")
	for _, route := range []string{
		"startSecFilterTask",
		"getTaskStatus",
		"startAITask",
	} {
		if !strings.Contains(router, route) {
			t.Fatalf("search_result router missing %s", route)
		}
	}
	api := readRepoFile(t, "server/api/search_result.go")
	if !strings.Contains(api, `var taskStatus = "stop"`) {
		t.Fatal("getTaskStatus must use in-process taskStatus (serve memory), not scan loop")
	}
	if !strings.Contains(api, "func StartSecFilterTask") || !strings.Contains(api, "func StartAITask") {
		t.Fatal("serve API must define StartSecFilterTask and StartAITask")
	}
	// These handlers must not call ScanTask
	for _, fn := range []string{"StartSecFilterTask", "StartAITask", "GetTaskStatus"} {
		// crude: whole file should not import search.ScanTask path
		_ = fn
	}
	if strings.Contains(api, "search.ScanTask") || strings.Contains(api, "ScanTask()") {
		t.Fatal("search_result API must not start ScanTask")
	}

	webAPI := readRepoFile(t, "web/src/api/searchResult.js")
	if !strings.Contains(webAPI, "/searchResult/startSecFilterTask") {
		t.Fatal("web client must call startSecFilterTask on searchResult API")
	}
	if !strings.Contains(webAPI, "/searchResult/startAITask") {
		t.Fatal("web client must call startAITask on searchResult API")
	}
	vue := readRepoFile(t, "web/src/view/searchResult/searchResult.vue")
	if !strings.Contains(vue, "startAITask") || !strings.Contains(vue, "startFilterTask") {
		t.Fatal("searchResult.vue must expose AI and secondary-filter actions")
	}
}

// TestArchitectureScanReadsRulesFromDB documents the config→execution handoff:
// Web writes rules; scan loads enabled rules from MySQL.
func TestArchitectureScanReadsRulesFromDB(t *testing.T) {
	gh := readRepoFile(t, "server/search/githubsearch/search.go")
	if !strings.Contains(gh, `GetValidRulesByType("github")`) {
		t.Fatal("github RunTask must load rules via GetValidRulesByType")
	}
	ruleSvc := readRepoFile(t, "server/service/rule.go")
	if !strings.Contains(ruleSvc, "status = 1") {
		t.Fatal("GetValidRulesByType should filter enabled rules (status=1)")
	}
	// task table is seeded but must not be required for RunTask (current design)
	if strings.Contains(gh, "model.Task") {
		t.Fatal("unexpected: github search should not depend on model.Task for RunTask")
	}
}

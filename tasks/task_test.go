package tasks_test

import (
	"gshark/models"
	"gshark/tasks"
	"gshark/util/index"
	"gshark/vars"

	"encoding/json"
	"testing"
)

func TestGenerateSearcher(t *testing.T) {
	reposCfg := make([]models.RepoConfig, 0)
	repoConfig := models.RepoConfig{Name: "xsec-traffic", Url: "https://github.com/netxfly/xsec-traffic",
		PollInterval: 30, Vcs: "git", UrlPattern: models.UrlPattern{BaseUrl: vars.DefaultBaseUrl, Anchor: vars.DefaultAnchor},
		AutoPullUpdate: true, ExcludeDotFiles: true,
	}

	reposCfg = append(reposCfg, repoConfig)

	t.Log(tasks.GenerateSearcher(reposCfg))
}

func TestSearchRepos(t *testing.T) {
	reposCfg := make([]models.RepoConfig, 0)
	repoConfig := models.RepoConfig{Name: "xsec-traffic", Url: "https://github.com/netxfly/xsec-traffic",
		PollInterval: 30, Vcs: "git", UrlPattern: models.UrlPattern{BaseUrl: vars.DefaultBaseUrl, Anchor: vars.DefaultAnchor},
		AutoPullUpdate: true, ExcludeDotFiles: true,
	}

	var filesOpened int
	var durationMs int
	reposCfg = append(reposCfg, repoConfig)
	searchers, errors, hasError, err := tasks.GenerateSearcher(reposCfg)
	t.Log(searchers, errors, hasError, err)

	repos := make([]string, 0)
	repos = append(repos, "xsec-traffic")
	rule := models.Rule{Part: "keyword", Type: "regex", Pattern: "password", Caption: "Contains word: password",
		Description: "Contains word: password"}
	opts := index.SearchOptions{IgnoreCase: true, LinesOfContext: tasks.DefaultLinesOfContext, FileRegexp: ""}
	respSearch, err := tasks.SearchRepos(rule, &opts, repos, searchers, &filesOpened, &durationMs)
	respJson, err := json.Marshal(respSearch)
	t.Log(string(respJson), err)
}

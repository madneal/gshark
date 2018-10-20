package searcher_test

import (
	"gshark/models"
	"gshark/vars"

	"testing"
	"gshark/util/searcher"
)

func TestMakeAll(t *testing.T) {
	reposCfg := make([]models.RepoConfig, 0)
	repoConfig := models.RepoConfig{Name: "netxfly", Url: "https://github.com/netxfly/xsec-traffic",
		PollInterval: 30, Vcs: "git", UrlPattern: models.UrlPattern{BaseUrl: vars.DefaultBaseUrl, Anchor: vars.DefaultAnchor},
		AutoPullUpdate: true, ExcludeDotFiles: true,
	}

	reposCfg = append(reposCfg, repoConfig)
	t.Log(searcher.MakeAll(reposCfg))
}

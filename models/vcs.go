package models

import (
	"x-patrol/vars"

	"time"
)

type UrlPattern struct {
	Id      int64
	BaseUrl string `json:"base-url"`
	Anchor  string `json:"anchor"`
	Vcs     string `json:"vcs"`
}

func NewUrlParttern(baseUrl, anchor, vcs string) *UrlPattern {
	return &UrlPattern{BaseUrl: baseUrl, Anchor: anchor, Vcs: vcs}
}

func (u *UrlPattern) Exist() (bool, error) {
	return Engine.Get(u)
}

func (u *UrlPattern) Insert() (int64, error) {
	return Engine.Insert(u)
}

func InitUrlPattern() {
	u := NewUrlParttern(vars.DefaultBaseUrl, vars.DefaultAnchor, "git")
	has, err := u.Exist()
	if err == nil && !has {
		u.Insert()
	}
}

func GetUrlPattern(vcs string) (bool, error, *UrlPattern) {
	u := new(UrlPattern)
	has, err := Engine.Table("url_pattern").Where("vcs=?", vcs).Get(u)
	return has, err, u
}

type RepoConfig struct {
	Id              int64
	Name            string        `json:"name"`
	Url             string        `json:"url"`
	PollInterval    time.Duration `json:"poll_interval"`
	Vcs             string        `json:"vcs"`
	UrlPattern      UrlPattern    `json:"url_pattern"`
	AutoPullUpdate  bool          `json:"auto_pull_update"`
	ExcludeDotFiles bool          `json:"exclude_dot_files"`
}

func NewRepoConfig(name string,
	url string,
	interval time.Duration,
	vcs string,
	urlPat UrlPattern,
	isPull bool,
	isExclude bool) *RepoConfig {
	return &RepoConfig{Name: name, Url: url, PollInterval: interval, Vcs: vcs, UrlPattern: urlPat,
		AutoPullUpdate:  isPull,
		ExcludeDotFiles: isExclude}
}

func (r *RepoConfig) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *RepoConfig) Exist() (bool, error) {
	has := false
	rs := make([]RepoConfig, 0)
	err := Engine.Where("url=?", r.Url).Find(&rs)
	if err == nil && len(rs) > 0 {
		has = true
	}
	return has, err
}

func ListRepoConfig() ([]RepoConfig, error) {
	reposConfig := make([]RepoConfig, 0)
	err := Engine.Find(&reposConfig)
	return reposConfig, err
}

func InsertReposConfig() {
	// first delete all repos config
	ClearReposConfig()
	_, _, urlPat := GetUrlPattern("git")
	repos, err := ListEnableRepos()
	if err == nil {
		for _, repo := range repos {
			repoCnf := NewRepoConfig(repo.Name, repo.Url, vars.DefaultPollInterval, "git",
				*urlPat, true, false)
			has, err := repoCnf.Exist()
			if err == nil && !has {
				repoCnf.Insert()
			}
		}
	}
}

func ClearReposConfig() error {
	_, err := Engine.Exec("delete from repo_config")
	return err
}

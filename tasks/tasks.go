package tasks

import (
	"gshark/logger"
	"gshark/models"
	"gshark/util/index"
	"gshark/util/searcher"
	"gshark/vars"

	"os"
	"strings"
	"sync"
	"time"
	"gshark/util/githubsearch"
)

const (
	DefaultLinesOfContext uint = 4
	MaxLinesOfContext     uint = 20
)

type Stats struct {
	FilesOpened int
	Duration    int
}

type SearchResponse struct {
	repo string
	res  *index.SearchResponse
	err  error
}

func GenerateSearcher(reposConfig []models.RepoConfig) (map[string]*searcher.Searcher, map[string]bool, bool, error) {
	errRepos := make(map[string]bool)
	hasError := false
	// Ensure we have a repos path
	if _, err := os.Stat(vars.REPO_PATH); err != nil {
		if err := os.MkdirAll(vars.REPO_PATH, os.ModePerm); err != nil {
			return nil, errRepos, hasError, err
		}
	}

	searchers, errs, err := searcher.MakeAll(reposConfig)
	if err == nil {
		if len(errs) > 0 {
			// NOTE: This mutates the original config so the repos
			// are not even seen by other code paths.

			for name := range errs {
				errRepos[name] = true
			}
			hasError = true
		}
	}

	return searchers, errRepos, hasError, nil
}

/* search repos, return a  map[string]*index.SearchResponse map */
func SearchRepos(
	rule models.Rule,
	opts *index.SearchOptions,
	repos []string,
	idx map[string]*searcher.Searcher,
	filesOpened *int,
	duration *int,
) (map[string]*index.SearchResponse, error) {
	query := rule.Pattern
	startedAt := time.Now()
	num := len(repos)

	// use a buffered channel to avoid routine leaks on errs.
	ch := make(chan *SearchResponse, num)
	for _, repo := range repos {
		go func(repo string) {
			fms, err := idx[repo].Search(query, opts)
			ch <- &SearchResponse{repo, fms, err}
		}(repo)
	}

	res := map[string]*index.SearchResponse{}
	for i := 0; i < num; i++ {

		r := <-ch
		r.res.RuleId = rule.Id
		r.res.RuleCaption = rule.Caption
		r.res.RulePattern = rule.Pattern

		if r.err != nil {
			return nil, r.err
		}

		if r.res.Matches == nil {
			continue
		}
		res[r.repo] = r.res
		*filesOpened += r.res.FilesOpened
	}

	*duration = int(time.Now().Sub(startedAt).Seconds() * 1000)

	return res, nil
}

func DoSearch(reposConfig []models.RepoConfig, rules models.Rule) (map[string]*index.SearchResponse, models.Rule, error) {
	searchers, errors, _, err := GenerateSearcher(reposConfig)
	respSearch := make(map[string]*index.SearchResponse)
	if err == nil {
		repos := make([]string, 0)
		for _, repoCfg := range reposConfig {
			repo := repoCfg.Name
			if !errors[repo] {
				repos = append(repos, repoCfg.Name)
			}
		}

		opts := index.SearchOptions{IgnoreCase: true, LinesOfContext: DefaultLinesOfContext}
		if strings.ToLower(rules.Part) == "keyword" {
			// search keyword from all files
			opts.FileRegexp = ""
		} else {
			// when rules.Part in ("filename", "path", "extension"), only search filename, and set rules.Pattern = "\\."
			opts.FileRegexp = rules.Pattern
			rules.Pattern = "\\."
		}

		var filesOpened int
		var durationMs int

		respSearch, err = SearchRepos(rules, &opts, repos, searchers, &filesOpened, &durationMs)
	}
	return respSearch, rules, err
}

// 分割任务为map形式，key为批次，value为一批models.RepoConfig
func SegmentationTask(reposConfig []models.RepoConfig) map[int][]models.RepoConfig {
	tasks := make(map[int][]models.RepoConfig)
	totalRepos := len(reposConfig)
	scanBatch := totalRepos / vars.MAX_Concurrency_REPOS

	for i := 0; i < scanBatch; i++ {
		curTask := reposConfig[vars.MAX_Concurrency_REPOS*i : vars.MAX_Concurrency_REPOS*(i+1)]
		tasks[i] = curTask
	}

	if totalRepos%vars.MAX_Concurrency_REPOS > 0 {
		n := len(tasks)
		tasks[n] = reposConfig[vars.MAX_Concurrency_REPOS*scanBatch : totalRepos]
	}
	return tasks
}

// 按批次分发、执行任务
func DistributionTask(tasksMap map[int][]models.RepoConfig, rules []models.Rule) {
	for _, rule := range rules {
		for _, reposConf := range tasksMap {
			Run(reposConf, rule)
		}
	}
}

func Run(reposConfig []models.RepoConfig, rule models.Rule) {
	var wg sync.WaitGroup
	wg.Add(len(reposConfig))
	for _, rConfig := range reposConfig {
		// wg.Add(1)
		reposCfg := make([]models.RepoConfig, 0)
		reposCfg = append(reposCfg, rConfig)

		go func(config []models.RepoConfig, rule models.Rule) {
			defer wg.Done()
			SaveSearchResult(DoSearch(reposCfg, rule))
		}(reposCfg, rule)
		// wg.Wait()
	}
	wg.Wait()
}

func SaveSearchResult(responses map[string]*index.SearchResponse, rule models.Rule, err error) {
	if err == nil {
		for repo, resp := range responses {
			result := models.NewSearchResult(resp.Matches,
				repo,
				resp.FilesWithMatch,
				resp.FilesOpened, resp.Duration,
				resp.Revision, rule)

			has, _ := result.Exist()
			if !has {
				result.Insert()
			}
		}
	}
}

func ScheduleTasks(duration time.Duration) {
	for {
		// insert repos from inputInfo
		githubsearch.InsertAllRepos()

		// insert all enable repos to repos config table
		models.InsertReposConfig()

		rules, err := models.GetRules()
		if err == nil {
			reposConfig, err := models.ListRepoConfig()
			if err == nil {
				mapTasks := SegmentationTask(reposConfig)
				DistributionTask(mapTasks, rules)
			}
		}

		logger.Log.Infof("Complete the scan local repos, start to sleep %v seconds", duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}

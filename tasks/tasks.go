package tasks

import (
	"gshark/models"
	"gshark/util/index"
	"gshark/util/searcher"
	"gshark/vars"

	"os"
	"time"
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


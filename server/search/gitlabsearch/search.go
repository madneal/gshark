package gitlabsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// defaultGitlabDiscoverPages/defaultGitlabBatchSize bound the per-cycle crawl
// fallback when global search is unavailable (see RunGlobalSearchTask). They
// are used whenever the matching config value is unset (<= 0).
const (
	defaultGitlabDiscoverPages   = 5
	defaultGitlabBatchSize       = 50
	maxConcurrentProjectSearches = 5
	gitlabRequestTimeout         = 2 * time.Minute
)

var (
	getGitlabClient = GetClient
	getGitlabRules  = service.GetValidRulesByType
)

func RunTask() model.ScanOutcome {
	client := getGitlabClient()
	if client == nil {
		message := "There is no client for GitLab; configure a GitLab token to enable this provider"
		color.Warnln(message)
		return model.ScanSkipped(message)
	}

	err, rules := getGitlabRules("gitlab")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType gitlab err", zap.Error(err))
		return model.ScanFailed("Failed to load GitLab rules: " + err.Error())
	}
	if len(rules) == 0 {
		message := "No enabled GitLab rules; provider skipped"
		global.GVA_LOG.Info(message)
		return model.ScanSkipped(message)
	}

	if !RunGlobalSearchTask(client, rules) {
		global.GVA_LOG.Info("GitLab global blob search unavailable (no Advanced Search), falling back to per-project crawl")
		GetProjects(client, discoverPages())
		RunSearchTaskByProject(NextProjectBatch(batchSize()), rules, client)
		return model.ScanSuccess(fmt.Sprintf("Completed %d GitLab rules with per-project fallback", len(rules)))
	}
	global.GVA_LOG.Info("Complete the scan of GitLab")
	return model.ScanSuccess(fmt.Sprintf("Completed %d GitLab rules", len(rules)))
}

func discoverPages() int {
	if global.GVA_CONFIG.Search.GitlabDiscoverPages > 0 {
		return global.GVA_CONFIG.Search.GitlabDiscoverPages
	}
	return defaultGitlabDiscoverPages
}

func batchSize() int {
	if global.GVA_CONFIG.Search.GitlabBatchSize > 0 {
		return global.GVA_CONFIG.Search.GitlabBatchSize
	}
	return defaultGitlabBatchSize
}

// isGlobalSearchUnsupported reports whether resp is GitLab's specific response
// for scope=blobs search without Advanced Search/Elasticsearch enabled (HTTP
// 400 "Scope not supported without Elasticsearch!"). GitLab.com disables this
// by default even on paid tiers; self-hosted instances may not have it
// configured either.
func isGlobalSearchUnsupported(resp *gitlab.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusBadRequest
}

// RunGlobalSearchTask attempts GitLab's global blobs search for every rule.
// It returns false as soon as the first rule reports the instance doesn't
// support it, so RunTask can fall back to the per-project crawl for the whole
// cycle instead of eating one guaranteed-failing request per rule. Instances
// with Advanced Search enabled (self-hosted with Elasticsearch, or GitLab.com
// Ultimate with it turned on) get full global coverage here.
func RunGlobalSearchTask(client *gitlab.Client, rules []model.Rule) bool {
	for i, rule := range rules {
		resp, ok := SearchBlobsStream(client, rule.Content, func(blobs []*gitlab.Blob) error {
			results := ConvertBlobsToResults(client, blobs, rule.Content)
			SaveResult(results, &rule.Content)
			return nil
		})
		if !ok {
			if i == 0 && isGlobalSearchUnsupported(resp) {
				return false
			}
			continue
		}
	}
	return true
}

// NextProjectBatch returns up to batchSize not-yet-searched gitlab repos. If
// every known repo has already been searched, it recycles the whole pool
// (resets status back to unsearched) so the crawl never permanently stalls.
func NextProjectBatch(batchSize int) []model.Repo {
	valid := ListValidProjects()
	if len(*valid) == 0 {
		if err := service.ResetRepoStatusByType("gitlab"); err != nil {
			global.GVA_LOG.Error("reset gitlab repo status error", zap.Error(err))
			return nil
		}
		valid = ListValidProjects()
	}
	if len(*valid) > batchSize {
		return (*valid)[:batchSize]
	}
	return *valid
}

// RunSearchTaskByProject runs every rule against every project in a bounded
// worker pool, then marks each project searched so subsequent cycles pick up
// where this one left off.
func RunSearchTaskByProject(projects []model.Repo, rules []model.Rule, client *gitlab.Client) {
	if len(projects) == 0 || len(rules) == 0 {
		return
	}

	sem := make(chan struct{}, maxConcurrentProjectSearches)
	var wg sync.WaitGroup
	for _, project := range projects {
		wg.Add(1)
		sem <- struct{}{}
		go func(project model.Repo) {
			defer wg.Done()
			defer func() { <-sem }()
			searchProjectForRules(project, rules, client)
		}(project)
	}
	wg.Wait()
}

func searchProjectForRules(project model.Repo, rules []model.Rule, client *gitlab.Client) {
	for _, rule := range rules {
		if err := SearchCodeStream(rule.Content, project, client, func(results []*model.SearchResult) error {
			SaveResult(results, &rule.Content)
			return nil
		}); err != nil {
			global.GVA_LOG.Error("search project stream error", zap.Error(err), zap.String("project", project.Path))
		}
	}
	project.Status = 1
	if err := service.UpdateRepo(project); err != nil {
		global.GVA_LOG.Error("mark gitlab project searched error", zap.Error(err))
	}
}

func SaveResult(results []*model.SearchResult, keyword *string) {
	if len(results) == 0 {
		return
	}
	stats := service.SaveSearchResultPointersWithStats(results, *keyword)
	global.GVA_LOG.Info(stats.Summary(*keyword, "GitLab"))
}

// SearchCode searches for keyword inside a single project, paginating through
// every page of results.
func SearchCode(keyword string, project model.Repo, client *gitlab.Client) []*model.SearchResult {
	codeResults := make([]*model.SearchResult, 0)
	_ = SearchCodeStream(keyword, project, client, func(results []*model.SearchResult) error {
		codeResults = append(codeResults, results...)
		return nil
	})
	return codeResults
}

// SearchCodeStream searches one project page by page and lets the caller
// persist each page before the next page is fetched.
func SearchCodeStream(keyword string, project model.Repo, client *gitlab.Client, onPage func([]*model.SearchResult) error) error {
	global.GVA_LOG.Info(fmt.Sprintf("Search inside project %s", project.Path))
	opt := &gitlab.SearchOptions{Page: 1, PerPage: 100}
	for {
		results, resp, err := client.Search.BlobsByProject(project.ProjectId, keyword, opt)
		if err != nil {
			global.GVA_LOG.Error("search inside project error", zap.Error(err))
			return err
		}
		if resp != nil && resp.StatusCode != http.StatusOK {
			global.GVA_LOG.Info(fmt.Sprintf("Request error for project statuscode %d", resp.StatusCode))
			return fmt.Errorf("search project returned status %d", resp.StatusCode)
		}
		pageResults := make([]*model.SearchResult, 0, len(results))
		for _, result := range results {
			url := project.Url + "/blob/master/" + result.Filename
			textMatches := make([]model.TextMatch, 0)
			textMatch := model.TextMatch{
				Fragment: &result.Data,
			}
			textMatches = append(textMatches, textMatch)
			b, err := json.Marshal(textMatches)
			if err != nil {
				global.GVA_LOG.Error("json marshal error", zap.Error(err))
			}
			codeResult := model.SearchResult{
				Path:            result.Filename,
				Repo:            result.Basename,
				Url:             url,
				TextMatchesJson: b,
				Status:          0,
				Keyword:         keyword,
			}
			pageResults = append(pageResults, &codeResult)
		}
		if err := onPage(pageResults); err != nil {
			return err
		}
		if resp == nil || resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil
}

// ListValidProjects returns every known gitlab repo that hasn't been searched
// yet (Status != 1).
func ListValidProjects() *[]model.Repo {
	validProjects := make([]model.Repo, 0)
	err, projects := service.GetRepoByType("gitlab")
	if err != nil {
		global.GVA_LOG.Error("list projects error", zap.Any("err", err))
	}
	for _, p := range projects {
		// if the project has been searched
		if p.Status == 1 {
			continue
		}
		validProjects = append(validProjects, p)
	}
	return &validProjects
}

func GetClient() *gitlab.Client {
	var baseURL string
	if global.GVA_CONFIG.System.GitlabBase != "" {
		baseURL = global.GVA_CONFIG.System.GitlabBase
	} else {
		baseURL = "https://gitlab.com"
	}
	err, tokens := service.ListTokenByType("gitlab")
	if len(tokens) == 0 {
		return nil
	}
	client, err := gitlab.NewClient(tokens[0].Content,
		gitlab.WithBaseURL(baseURL),
		gitlab.WithHTTPClient(&http.Client{Timeout: gitlabRequestTimeout}),
	)
	if err != nil {
		global.GVA_LOG.Error("getClient error", zap.Error(err))
	}
	return client
}

// SearchBlobBySearchOptions is utilized to search inside blob by keyword
func SearchBlobBySearchOptions(client *gitlab.Client, keyword string, searchOptions *gitlab.SearchOptions) ([]*gitlab.Blob, *gitlab.Response, error) {
	blobs, resp, err := client.Search.Blobs(keyword, searchOptions)
	return blobs, resp, err
}

// SearchBlobs is utilized to search all the results by keyword via GitLab's
// global scope=blobs search. This only succeeds when the instance has
// Advanced Search/Elasticsearch enabled; the returned bool is false otherwise
// (or on any other request failure), with resp carrying the failing response
// so callers can tell "unsupported" apart from a transient error.
func SearchBlobs(client *gitlab.Client, keyword string) ([]*gitlab.Blob, *gitlab.Response, bool) {
	blobs := make([]*gitlab.Blob, 0)
	resp, ok := SearchBlobsStream(client, keyword, func(page []*gitlab.Blob) error {
		blobs = append(blobs, page...)
		return nil
	})
	return blobs, resp, ok
}

// SearchBlobsStream fetches GitLab global-search pages one at a time. The
// callback runs before the next page is requested to bound retained results.
func SearchBlobsStream(client *gitlab.Client, keyword string, onPage func([]*gitlab.Blob) error) (*gitlab.Response, bool) {
	opt := &gitlab.SearchOptions{Page: 1, PerPage: 100}
	for {
		page, resp, err := SearchBlobBySearchOptions(client, keyword, opt)
		if err != nil {
			global.GVA_LOG.Error("SearchBlobs error", zap.Error(err))
			return resp, false
		}
		if err := onPage(page); err != nil {
			return resp, false
		}
		if resp == nil || resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return nil, true
}

// GetProjectById is utilized to get the project by id
func GetProjectById(client *gitlab.Client, id int) *gitlab.Project {
	project, _, err := client.Projects.GetProject(id, &gitlab.GetProjectOptions{})
	if err != nil {
		global.GVA_LOG.Error("GetProjectById err", zap.Error(err))
	}
	return project
}

// ConvertBlobsToResults is utilized to convert blobs to results
func ConvertBlobsToResults(client *gitlab.Client, blobs []*gitlab.Blob, keyword string) []*model.SearchResult {
	results := make([]*model.SearchResult, 0)
	for _, blob := range blobs {
		projectId := blob.ProjectID
		project := GetProjectById(client, projectId)
		textMatches := make([]model.TextMatch, 0)
		textMatches = append(textMatches, model.TextMatch{
			Fragment: &blob.Data,
		})
		dataJson, err := json.Marshal(textMatches)
		if err != nil {
			global.GVA_LOG.Error("blob.Data marshal error", zap.Error(err))
		}
		result := model.SearchResult{
			Url:             fmt.Sprintf("%s/blob/master/%s", project.WebURL, blob.Filename),
			Path:            blob.Filename,
			Repo:            blob.Basename,
			TextMatchesJson: dataJson,
			Status:          0,
			Keyword:         keyword,
		}
		results = append(results, &result)
	}
	return results
}

// GetProjects discovers recently-active public projects from gitlab, upserting
// them into the repo table for the crawl fallback. maxPages bounds how many
// pages are fetched in a single call (<= 0 means unbounded).
func GetProjects(client *gitlab.Client, maxPages int) {
	isSimple := true
	date := time.Now().AddDate(0, -1, 0)
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Simple:            &isSimple,
		LastActivityAfter: &date,
	}
	projectNum := 0
	pagesFetched := 0
	for {
		// Get the first page with projects.
		ps, resp, err := client.Projects.ListProjects(opt)
		if err != nil {
			fmt.Println(err)
			break
		}

		// List all the projects we've found so far.
		for _, p := range ps {
			if strings.HasPrefix(p.PathWithNamespace, "gitlab") {
				continue
			}
			repo := model.Repo{
				Url:            p.WebURL,
				Path:           p.PathWithNamespace,
				Type:           "gitlab",
				ProjectId:      p.ID,
				Status:         0,
				LastActivityAt: *(p.LastActivityAt),
			}
			fmt.Println(repo.Path)
			err, has := service.CheckRepoExist(&repo)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Error("CheckRepoExist error", zap.Error(err))
			}
			if !has {
				err := service.CreateRepo(repo)
				if err != nil {
					global.GVA_LOG.Error("CreateRepo error", zap.Error(err))
				}
				projectNum++
			}
		}

		pagesFetched++
		if maxPages > 0 && pagesFetched >= maxPages {
			break
		}

		if resp.NextPage == 0 {
			fmt.Println("next page is 0")
			break
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("request error: %d", resp.StatusCode)
			break
		}

		opt.Page = resp.NextPage
	}
	global.GVA_LOG.Info(fmt.Sprintf("Has inserted %d projects", projectNum))
}

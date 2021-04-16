package gitlabsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

func RunTask(duration time.Duration) {
	RunSearchTask(GenerateSearchCodeTask())
	global.GVA_LOG.Info(fmt.Sprintf("Complete the scan of Gitlab, start to sleep %d seconds", duration))
	time.Sleep(duration * time.Second)
}

func GenerateSearchCodeTask() (map[int][]model.Rule, error) {
	result := make(map[int][]model.Rule)
	// get rules with the type of github
	err, rules := service.GetValidRulesByType("gitlab")
	ruleNum := len(rules)
	batch := ruleNum / global.GVA_CONFIG.Search.SearchNum

	for i := 0; i < batch; i++ {
		result[i] = rules[global.GVA_CONFIG.Search.SearchNum*i : global.GVA_CONFIG.Search.SearchNum*(i+1)]
	}

	if ruleNum%global.GVA_CONFIG.Search.SearchNum != 0 {
		result[batch] = rules[global.GVA_CONFIG.Search.SearchNum*batch : ruleNum]
	}
	return result, err
}

func RunSearchTask(mapRules map[int][]model.Rule, err error) {
	client := GetClient()
	if client == nil {
		return
	}
	// get all public projects
	GetProjects(client)
	if err == nil {
		for _, rules := range mapRules {
			startTime := time.Now()
			Search(rules, client)
			usedTime := time.Since(startTime).Seconds()
			if usedTime < 60 {
				time.Sleep(time.Duration(60 - usedTime))
			}
		}
	}
}

func Search(rules []model.Rule, client *gitlab.Client) {
	var wg sync.WaitGroup
	wg.Add(len(rules))

	for _, rule := range rules {
		go func(rule model.Rule) {
			defer wg.Done()
			SearchInsideProjects(rule.Content, client)
		}(rule)
	}
	wg.Wait()
}

func SearchInsideProjects(keyword string, client *gitlab.Client) {
	err, projects := service.GetRepoByType("gitlab")
	if err != nil {
		global.GVA_LOG.Error("list projects error", zap.Any("err", err))
	}
	for _, project := range projects {
		results := SearchCode(keyword, project, client)
		SaveResult(results, &keyword)
	}
}

func SaveResult(results []*model.SearchResult, keyword *string) {
	insertCount := 0
	if len(results) > 0 {
		for _, resultItem := range results {
			err, has := service.CheckExistOfSearchResult(resultItem)
			if err != nil {
			}
			if !has {
				resultItem.Keyword = *keyword
				err := service.CreateSearchResult(*resultItem)
				if err != nil {
				}
				insertCount++
			}

		}
		global.GVA_LOG.Info(fmt.Sprintf("Has inserted %d results", insertCount))
	}
}

func SearchCode(keyword string, project model.Repo, client *gitlab.Client) []*model.SearchResult {
	codeResults := make([]*model.SearchResult, 0)
	//queryString := BuildQueryString(keyword, "ext")
	global.GVA_LOG.Info(fmt.Sprintf("Search inside project %s", project.Path))
	results, resp, err := client.Search.BlobsByProject(project.ProjectId, keyword, &gitlab.SearchOptions{})
	if err != nil {
		global.GVA_LOG.Error("search inside project error", zap.Error(err))
	}
	if resp != nil && resp.StatusCode != 200 {
		global.GVA_LOG.Info(fmt.Sprintf("Request error for project statuscode %d", resp.StatusCode))
		return codeResults
	}
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
		//if !mergeTextMatches(codeResults, result.Filename, textMatch) {
		codeResults = append(codeResults, &codeResult)
		//}

	}
	return codeResults
}

// mergeTextMatches is utilized to merge multi textMatches in the same file
// return: if has merged
//func mergeTextMatches(codeResults []*model.SearchResult, filename string, textMatch models.TextMatch) bool {
//	flag := false
//	for index, result := range codeResults {
//		if *result.Name == filename {
//			flag = true
//			codeResults[index].TextMatches = append(codeResults[index].TextMatches, textMatch)
//			return flag
//		}
//	}
//	return flag
//}

func ListValidProjects() *[]model.Repo {
	validProjects := make([]model.Repo, 0)
	err, projects := service.GetRepoByType("gitlab")
	if err != nil {
		global.GVA_LOG.Error("list projects error", zap.Error(err))
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
	err, tokens := service.ListTokenByType("gitlab")
	if len(tokens) == 0 {
		return nil
	}
	client, err := gitlab.NewClient(tokens[0].Content)
	if err != nil {
		global.GVA_LOG.Error("getClient error", zap.Error(err))
	}
	return client
}

// GetProjects is utilized to obtain public projects from gitlab
func GetProjects(client *gitlab.Client) {
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
					global.GVA_LOG.Error("creareRepo error", zap.Error(err))
				}
				projectNum++
			}
		}

		if resp.NextPage == 0 {
			fmt.Println("next page is 0")
			break
		}

		if resp.StatusCode != 200 {
			fmt.Printf("request error: %d", resp.StatusCode)
			break
		}

		opt.Page = resp.NextPage
	}
	global.GVA_LOG.Info(fmt.Sprintf("Has inserted %d projects", projectNum))
}

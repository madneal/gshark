package gitlabsearch

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/models"
	"github.com/madneal/gshark/vars"
	"github.com/xanzy/go-gitlab"
	"strings"
	//"sync"
	"time"
)

func RunTask(duration time.Duration) {
	//RunSearchTask(GenerateSearchCodeTask())

	logger.Log.Infof("Complete the scan of Gitlab, start to sleep %v seconds", duration*time.Second)
	time.Sleep(duration * time.Second)
}

//func GenerateSearchCodeTask() (map[int][]models.Rule, error) {
//	result := make(map[int][]models.Rule)
//	// get rules with the type of github
//	rules, err := models.GetValidRulesByType(vars.GITLAB)
//	ruleNum := len(rules)
//	batch := ruleNum / vars.SearchNum
//
//	for i := 0; i < batch; i++ {
//		result[i] = rules[vars.SearchNum*i : vars.SearchNum*(i+1)]
//	}
//
//	if ruleNum%vars.SearchNum != 0 {
//		result[batch] = rules[vars.SearchNum*batch : ruleNum]
//	}
//	return result, err
//}
//
//func RunSearchTask(mapRules map[int][]models.Rule, err error) {
//	client := GetClient()
//	if client == nil {
//		return
//	}
//	// get all public projects
//	GetProjects(client)
//	if err == nil {
//		for _, rules := range mapRules {
//			startTime := time.Now()
//			//Search(rules, client)
//			usedTime := time.Since(startTime).Seconds()
//			if usedTime < 60 {
//				time.Sleep(time.Duration(60 - usedTime))
//			}
//		}
//	}
//}

//func Search(rules []models.Rule, client *gitlab.Client) {
//	var wg sync.WaitGroup
//	wg.Add(len(rules))
//
//	for _, rule := range rules {
//		go func(rule models.Rule) {
//			defer wg.Done()
//			SearchInsideProjects(rule.Pattern, client)
//		}(rule)
//	}
//	wg.Wait()
//}

func SearchCode(keyword string, client *resty.Client) *resty.Response {
	res := client.R()
	res.SetHeader("PRIVATE-TOKEN", vars.GITLAB_TOKEN)
	url := "https://gitlab.com/api/v4/search?scope=blobs&search=" + keyword
	response, err := res.Get(url)
	if err != nil {
		logger.Log.Error(err)
	}
	fmt.Println(response)
	return response
}

//func SearchInsideProjects(keyword string, client *gitlab.Client) {
//	projects := ListValidProjects()
//	for _, project := range projects {
//		results := SearchCode(keyword, project, client)
//		SaveResult(results, &keyword)
//	}
//}

func SaveResult(results []*models.CodeResult, keyword *string) {
	insertCount := 0
	if len(results) > 0 {
		for _, resultItem := range results {
			has, err := resultItem.Exist()
			if err != nil {
				logger.Log.Error(err)
			}
			if !has {
				resultItem.Keyword = keyword
				_, err := resultItem.Insert()
				if err != nil {
					logger.Log.Error(err)
				}
				insertCount++
			}

		}
		logger.Log.Infof("Has inserted %d results into code_result", insertCount)
	}
}

// buildQueryString is utilize to build query string
// add extension filters
func BuildQueryString(keyword, key string) string {
	filterRules, err := models.GetFilterRules()
	queryString := keyword
	if err != nil {
		logger.Log.Error(err)
	}
	for _, filterRule := range filterRules {
		ruleValue := filterRule.RuleValue
		ruleType := filterRule.RuleType
		ruleKey := filterRule.RuleKey
		ruleValueList := strings.Split(ruleValue, ",")
		if ruleKey != key {
			continue
		}
		for _, value := range ruleValueList {
			if ruleType == 0 {
				queryString += " -"
			} else {
				queryString += " +"
			}

			if ruleKey == "ext" {
				queryString += "extension:"
			}

			value = strings.TrimSpace(value)
			queryString += value
		}
	}
	return queryString
}

//func SearchCode(keyword string, project models.InputInfo, client *gitlab.Client) []*models.CodeResult {
//	codeResults := make([]*models.CodeResult, 0)
//	queryString := BuildQueryString(keyword, "ext")
//	//logger.Log.Infof("Search inside project %s", project.Url)
//	results, resp, err := client.Search.BlobsByProject(project.ProjectId, queryString, &gitlab.SearchOptions{})
//	if err != nil {
//		logger.Log.Error(err)
//	}
//	if resp.StatusCode != 200 {
//		fmt.Printf("request error for projectId-%d: %d\n", project.ProjectId, resp.StatusCode)
//		if resp.StatusCode == 404 {
//			err = project.DeleteByProjectId()
//			if err != nil {
//				logger.Log.Error(err)
//			}
//		}
//		return codeResults
//	}
//	for _, result := range results {
//		url := project.Url + "/blob/master/" + result.Filename
//		textMatches := make([]models.TextMatch, 0)
//		textMatch := models.TextMatch{
//			Fragment: &result.Data,
//		}
//		textMatches = append(textMatches, textMatch)
//		codeResult := models.CodeResult{
//			Id:          0,
//			Name:        &result.Filename,
//			Path:        &result.Basename,
//			RepoName:    result.Basename,
//			HTMLURL:     &url,
//			TextMatches: textMatches,
//			Status:      0,
//			Keyword:     &keyword,
//		}
//		if !mergeTextMatches(codeResults, result.Filename, textMatch) {
//			codeResults = append(codeResults, &codeResult)
//		}
//	}
//	err = models.UpdateStatusById(1, project.ProjectId)
//	if err != nil {
//		logger.Log.Error(err)
//	}
//	return codeResults
//}

// mergeTextMatches is utilized to merge multi textMatches in the same file
// return: if has merged
func mergeTextMatches(codeResults []*models.CodeResult, filename string, textMatch models.TextMatch) bool {
	flag := false
	for index, result := range codeResults {
		if *result.Name == filename {
			flag = true
			codeResults[index].TextMatches = append(codeResults[index].TextMatches, textMatch)
			return flag
		}
	}
	return flag
}

func ListValidProjects() []models.InputInfo {
	validProjects := make([]models.InputInfo, 0)
	projects, err := models.ListInputInfoByType(vars.GITLAB)
	if err != nil {
		logger.Log.Error(err)
	}
	for _, p := range projects {
		// if the project has been searched
		//if p.Status == 1 {
		//	continue
		//}
		validProjects = append(validProjects, p)
	}
	return validProjects
}

func GetClient() *gitlab.Client {
	tokens, err := models.ListValidTokens(vars.GITLAB)
	if err != nil {
		logger.Log.Error(err)
	}
	if len(tokens) == 0 {
		logger.Log.Warn("There is no valid gitlab token")
		return nil
	}
	return gitlab.NewClient(nil, tokens[0].Token)
}

// GetProjects is utilized to obtain public projects from gitlab
func GetProjects(client *gitlab.Client) {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
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
			inputInfo := models.InputInfo{
				Url:         p.WebURL,
				Path:        p.PathWithNamespace,
				Type:        vars.GITLAB,
				ProjectId:   p.ID,
				Status:      2,
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
			}
			has, err := inputInfo.Exist()
			if err != nil {
				fmt.Println(err)
			}
			if !has {
				//logger.Log.Infof("Insert project %s", p.WebURL)
				_, err := inputInfo.Insert()
				if err != nil {
					logger.Log.Error(err)
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
	logger.Log.Infof("Has found %d projects", projectNum)
}

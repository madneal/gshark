package gitlabsearch

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"github.com/xanzy/go-gitlab"
	"time"
)

func RunTask(duration time.Duration) {
	//RunSearchTask(GenerateSearchCodeTask())

	//logger.Log.Infof("Complete the scan of Gitlab, start to sleep %v seconds", duration*time.Second)
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

//func RunSearchTask(mapRules map[int][]model.Rule, err error) {
//	client := GetClient()
//	if client == nil {
//		return
//	}
//	// get all public projects
//	GetProjects(client)
//	if err == nil {
//		for _, rules := range mapRules {
//			startTime := time.Now()
//			Search(rules, client)
//			usedTime := time.Since(startTime).Seconds()
//			if usedTime < 60 {
//				time.Sleep(time.Duration(60 - usedTime))
//			}
//		}
//	}
//}

//func Search(rules []model.Rule, client *gitlab.Client) {
//	var wg sync.WaitGroup
//	wg.Add(len(rules))
//
//	for _, rule := range rules {
//		go func(rule model.Rule) {
//			defer wg.Done()
//			SearchInsideProjects(rule.Content, client)
//		}(rule)
//	}
//	wg.Wait()
//}

//func SearchInsideProjects(keyword string, client *gitlab.Client) {
//	projects := ListValidProjects()
//	for _, project := range projects {
//		results := SearchCode(keyword, project, client)
//		SaveResult(results, &keyword)
//	}
//}

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
		//logger.Log.Infof("Has inserted %d results into code_result", insertCount)
	}
}

//func SearchCode(keyword string, project models.InputInfo, client *gitlab.Client) []*model.SearchResult {
//	codeResults := make([]*model.SearchResult, 0)
//	//queryString := BuildQueryString(keyword, "ext")
//	//logger.Log.Infof("Search inside project %s", project.Url)
//	results, resp, err := client.Search.BlobsByProject(project.ProjectId, queryString, &gitlab.SearchOptions{})
//	if err != nil {
//		//logger.Log.Error(err)
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
//		//logger.Log.Error(err)
//	}
//	return codeResults
//}
//
//// mergeTextMatches is utilized to merge multi textMatches in the same file
//// return: if has merged
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
//
//func ListValidProjects() []model.Repo {
//	validProjects := make([]model.Repo, 0)
//	projects, err := service.GetRepoInfoList()
//	if err != nil {
//		logger.Log.Error(err)
//	}
//	for _, p := range projects {
//		// if the project has been searched
//		//if p.Status == 1 {
//		//	continue
//		//}
//		validProjects = append(validProjects, p)
//	}
//	return validProjects
//}

func GetClient() *gitlab.Client {
	err, tokens := service.ListTokenByType("")
	if len(tokens) == 0 {
		return nil
	}
	client, err := gitlab.NewClient(tokens[0].Content)
	if err != nil {
		//logger.Log.Error(err)
	}
	return client
}

//// GetProjects is utilized to obtain public projects from gitlab
//func GetProjects(client *gitlab.Client) {
//	opt := &gitlab.ListProjectsOptions{
//		ListOptions: gitlab.ListOptions{
//			PerPage: 100,
//			Page:    1,
//		},
//	}
//	projectNum := 0
//	for {
//		// Get the first page with projects.
//		ps, resp, err := client.Projects.ListProjects(opt)
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//
//		// List all the projects we've found so far.
//		for _, p := range ps {
//			inputInfo := models.InputInfo{
//				Url:         p.WebURL,
//				Path:        p.PathWithNamespace,
//				Type:        vars.GITLAB,
//				ProjectId:   p.ID,
//				Status:      2,
//				CreatedTime: time.Now(),
//				UpdatedTime: time.Now(),
//			}
//			has, err := inputInfo.Exist()
//			if err != nil {
//				fmt.Println(err)
//			}
//			if !has {
//				//logger.Log.Infof("Insert project %s", p.WebURL)
//				_, err := inputInfo.Insert()
//				if err != nil {
//					logger.Log.Error(err)
//				}
//				projectNum++
//			}
//		}
//
//		if resp.NextPage == 0 {
//			fmt.Println("next page is 0")
//			break
//		}
//
//		if resp.StatusCode != 200 {
//			fmt.Printf("request error: %d", resp.StatusCode)
//			break
//		}
//
//		opt.Page = resp.NextPage
//	}
//	logger.Log.Infof("Has found %d projects", projectNum)
//}

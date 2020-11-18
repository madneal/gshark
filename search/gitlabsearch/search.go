package gitlabsearch

import (
	"github.com/go-resty/resty/v2"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/misc"
	"github.com/madneal/gshark/models"
	"github.com/madneal/gshark/vars"
	"strconv"
	"strings"
	"sync"

	//"sync"
	"time"
)

func RunTask(duration time.Duration) {
	//RunSearchTask(GenerateSearchCodeTask())

	logger.Log.Infof("Complete the scan of Gitlab, start to sleep %v seconds", duration*time.Second)
	time.Sleep(duration * time.Second)
}

func GenerateSearchCodeTask() (map[int][]models.Rule, error) {
	result := make(map[int][]models.Rule)
	// get rules with the type of github
	rules, err := models.GetValidRulesByType(vars.GITLAB)
	ruleNum := len(rules)
	batch := ruleNum / vars.SearchNum

	for i := 0; i < batch; i++ {
		result[i] = rules[vars.SearchNum*i : vars.SearchNum*(i+1)]
	}

	if ruleNum%vars.SearchNum != 0 {
		result[batch] = rules[vars.SearchNum*batch : ruleNum]
	}
	return result, err
}

func RunSearchTask(mapRules map[int][]models.Rule, err error) {
	if err == nil {
		for _, rules := range mapRules {
			startTime := time.Now()
			//Search(rules, client)
			usedTime := time.Since(startTime).Seconds()
			if usedTime < 60 {
				time.Sleep(time.Duration(60 - usedTime))
			}
		}
	}
}

func Search(rules []models.Rule) {
	var wg sync.WaitGroup
	wg.Add(len(rules))
	client := resty.New()

	for _, rule := range rules {
		go func(rule models.Rule) {
			defer wg.Done()
			SearchCode(rule.Pattern, client)
		}(rule)
	}
	wg.Wait()
}

func SearchCode(keyword string, client *resty.Client) []byte {
	if vars.GITLAB_TOKEN == "" {
		logger.Log.Error("The gitlab token is empty, please specify the gitlab token!")
		return nil
	}
	var result []byte
	var nextPage string
	req := client.R()
	req.SetHeader("PRIVATE-TOKEN", vars.GITLAB_TOKEN)
	query, err := BuildQuery(keyword)
	if err != nil {
		logger.Log.Error(err)
	}
	url := vars.GitlabSearchUrl + query + "&page=1"
	result, nextPage = DoGet(url, req)
	SaveResult(result, &keyword, client)
	for nextPage != "" {
		url = vars.GitlabSearchUrl + query + "&page=" + nextPage
		result, nextPage = DoGet(url, req)
		SaveResult(result, &keyword, client)
	}
	return nil
}

func DoGet(url string, req *resty.Request) ([]byte, string) {
	var result []byte
	var nextPage string
	response, err := req.Get(url)
	if response != nil {
		result = response.Body()
		nextPage = response.Header().Get("X-Next-Page")
	}
	if err != nil {
		logger.Log.Error(err)
	}
	return result, nextPage
}

func ConvertStringToCodeResults(result []byte, keyword string, client *resty.Client) []*models.CodeResult {
	codeResults := make([]*models.CodeResult, 0)
	gitlabRes := Parse(result)
	for _, res := range gitlabRes {
		projectId := res.Project_id
		r := client.R()
		projectUrl := "https://gitlab.com/api/v4/projects/" + strconv.Itoa(int(projectId))
		response, err := r.Get(projectUrl)
		if err != nil {
			logger.Log.Error(err)
		}
		project := ParseProject(response.Body())
		url := "https://gitlab.com/" + project.Name + "/~/blob/master/" + res.Filename
		textMatchmd5 := misc.GenMd5WithSpecificLen(res.Data, 50)
		codeResult := models.CodeResult{
			Name:         &res.Filename,
			Path:         &res.Basename,
			RepoName:     res.Basename,
			HTMLURL:      &url,
			Status:       0,
			Keyword:      &keyword,
			Textmatchmd5: &textMatchmd5,
		}
		codeResults = append(codeResults, &codeResult)
	}
	return codeResults
}

func BuildQuery(query string) (string, error) {
	filterRules, err := models.GetFilterRules()
	str := ""
	for _, filterRule := range filterRules {
		ruleValue := filterRule.RuleValue
		ruleType := filterRule.RuleType
		ruleKey := filterRule.RuleKey
		ruleValueList := strings.Split(ruleValue, ",")
		for _, value := range ruleValueList {
			if ruleType == 0 {
				str += " -"
			} else {
				str += " +"
			}

			if ruleKey == "ext" {
				str += "extension:"
			}

			value = strings.TrimSpace(value)
			str += value
		}
	}
	builtQuery := query + str
	return builtQuery, err
}

func SaveResult(data []byte, keyword *string, client *resty.Client) {
	results := ConvertStringToCodeResults(data, *keyword, client)
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

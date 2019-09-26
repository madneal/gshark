package codesearch

import (
	"encoding/json"
	"fmt"
	"github.com/neal1991/gshark/logger"
	"github.com/neal1991/gshark/models"
	"github.com/neal1991/gshark/vars"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"strings"
	"time"
)

func ScheduleTasks(duration time.Duration) {
	for {
		RunSearchTask(GenerateSearchCodeTask())
		logger.Log.Infof("Complete the scan of searchcode, start to sleep %v seconds", duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}

func GenerateSearchCodeTask() (map[int][]models.Rule, error) {
	result := make(map[int][]models.Rule)
	rules, err := models.GetValidRulesByType("searchcode")
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
	request := gorequest.New()
	if err == nil {
		for _, rules := range mapRules {
			for _, rule := range rules {
				codeResults := SearchForSearchCode(rule, request)
				SaveResults(codeResults, &rule.Pattern)
			}
		}
	}
}

func SaveResults(results []*models.CodeResult, keyword *string) {
	insertCount := 0
	for _, result := range results {
		if result != nil {
			exist, err := result.Exist()
			result.Keyword = keyword
			if err != nil {
				fmt.Println(err)
			}
			if !exist {
				result.Insert()
				insertCount++
			}
		}
		logger.Log.Infof("Has inserted %d results into code_result", insertCount)
	}
}

func SearchForSearchCode(rule models.Rule, request *gorequest.SuperAgent) []*models.CodeResult {
	keyword := rule.Pattern
	totalCodeResults := make([]*models.CodeResult, 0)
	page := 0
	for {
		url := "https://searchcode.com/api/codesearch_I/?q=" + keyword + "&p=" + strconv.Itoa(page) + "&per_page=100"
		codeResults, total := GetResult(request, url)
		totalCodeResults = append(totalCodeResults, codeResults...)
		page++
		if page == total {
			break
		}
	}
	return totalCodeResults
}

func GetResult(request *gorequest.SuperAgent, url string) ([]*models.CodeResult, int) {
	codeResults := make([]*models.CodeResult, 0)
	resp, body, err := request.Get(url).End()
	if err != nil {
		logger.Log.Error(err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Request to %s error, status code: %d", url, resp.StatusCode)
	}
	var result models.SearchCodeRes
	//fmt.Println(body)
	json.Unmarshal([]byte(body), &result)
	total := result.Total
	//fmt.Println(total)
	results := result.Results
	//fmt.Println(results)
	for _, val := range results {
		if strings.Contains(val.Repo, "github") {
			continue
		}
		//fmt.Println(val.Filename)
		var lines string
		for _, line := range val.Lines {
			lines += fmt.Sprint(line) + "\n"
		}
		repoPath := val.Repo
		textMatch := new(models.TextMatch)
		textMatch.Fragment = &lines
		textMatchs := make([]models.TextMatch, 0)
		textMatchs = append(textMatchs, *textMatch)
		codeResult := models.CodeResult{
			Name:        &val.Filename,
			RepoName:    val.Location,
			Status:      0,
			HTMLURL:     &val.Url,
			RepoPath:    &repoPath,
			TextMatches: textMatchs,
		}
		codeResults = append(codeResults, &codeResult)
	}
	return codeResults, total
}

package codesearch

import (
	"encoding/json"
	"fmt"
	"github.com/neal1991/gshark/logger"
	"github.com/neal1991/gshark/models"
	"github.com/neal1991/gshark/vars"
	"github.com/parnurzeal/gorequest"
	"time"
)

func ScheduleTasks(duration time.Duration) {
	for {
		RunSearchTask(GenerateSearchCodeTask())
		logger.Log.Infof("Complete the scan of APP, start to sleep %v seconds", duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}

func GenerateSearchCodeTask() (map[int][]models.Rule, error) {
	result := make(map[int][]models.Rule)
	rules, err := models.GetValidRulesByType("app")
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
				SearchForSearchCode(rule, request)
				//SaveResults(results)
			}
		}
	}
}

func SearchForSearchCode(rule models.Rule, request *gorequest.SuperAgent)  {
	keyword := rule.Pattern
	url := "https://searchcode.com/api/codesearch_I/?q=" + keyword + "&p=0&per_page=100"
	GetResult(request, url)
}

func GetResult(request *gorequest.SuperAgent, url string)  []*models.CodeResult{
	codeResults := make([]*models.CodeResult, 0)
	resp, body, err := request.Get(url).End()
	if err != nil {
		logger.Log.Error(err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Request to %s error, status code: %d", url, resp.StatusCode)
	}
	var result models.SearchCodeRes
	fmt.Println(body)
	json.Unmarshal([]byte(body), &result)
	total := result.Total
	fmt.Println(total)
	results := result.Results
	fmt.Println(results)
	for _, val := range results {
		fmt.Println(val.Filename)
		var lines string
		for _, line := range val.Lines {
			lines += fmt.Sprint(line) + "\n"
		}
		fmt.Println(lines)
		codeResult := models.CodeResult{
			Name: &val.Name,
			Status: 0,

		}
	}
	return codeResults
}

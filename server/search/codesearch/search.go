package codesearch

import (
	"encoding/json"
	"fmt"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"github.com/parnurzeal/gorequest"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func RunTask(duration time.Duration) {
	RunSearchTask(GenerateSearchCodeTask())
	global.GVA_LOG.Info(fmt.Sprintf("Complete the scan of searchcode, start to sleep %d seconds", duration))
	time.Sleep(duration * time.Second)
}

func GenerateSearchCodeTask() (map[int][]model.Rule, error) {
	result := make(map[int][]model.Rule)
	err, rules := service.GetValidRulesByType("searchcode")
	ruleNum := len(rules)
	searchNum := global.GVA_CONFIG.Search.SearchNum
	batch := ruleNum / searchNum

	for i := 0; i < batch; i++ {
		result[i] = rules[searchNum*i : searchNum*(i+1)]
	}

	if ruleNum%searchNum != 0 {
		result[batch] = rules[searchNum*batch : ruleNum]
	}
	return result, err
}

func RunSearchTask(mapRules map[int][]model.Rule, err error) {
	request := gorequest.New()
	if err == nil {
		for _, rules := range mapRules {
			for _, rule := range rules {
				global.GVA_LOG.Info(fmt.Sprintf("Search for %s in searchcode", rule.Content))
				codeResults := SearchForSearchCode(rule, request)
				SaveResults(codeResults, &rule.Content)
			}
		}
	}
}

func SaveResults(results []*model.SearchResult, keyword *string) {
	insertCount := 0
	for _, result := range results {
		if result != nil {
			err, exist := service.CheckExistOfSearchResult(result)
			result.Keyword = *keyword
			if err != nil {
				fmt.Println(err)
			}
			if !exist {
				err = service.CreateSearchResult(*result)
				insertCount++
			}
			if err != nil {
				global.GVA_LOG.Error("search code save result error", zap.Any("err", err))
			}
		}
		global.GVA_LOG.Info(fmt.Sprintf("Has inserted %d results into code_result", insertCount))
	}
}

func SearchForSearchCode(rule model.Rule, request *gorequest.SuperAgent) []*model.SearchResult {
	keyword := rule.Content
	totalCodeResults := make([]*model.SearchResult, 0)
	page := 0
	for {
		url := "https://searchcode.com/api/codesearch_I/?q=" + keyword + "&p=" + strconv.Itoa(page)
		codeResults, hasResult := GetResult(request, url)
		totalCodeResults = append(totalCodeResults, codeResults...)
		page++
		if !hasResult {
			break
		}
	}
	return totalCodeResults
}

func GetResult(request *gorequest.SuperAgent, url string) ([]*model.SearchResult, bool) {
	hasResult := true
	codeResults := make([]*model.SearchResult, 0)
	resp, body, err := request.Get(url).End()
	if err != nil {
		global.GVA_LOG.Error("search result of search code error", zap.Any("err", err))
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Request to %s error, status code: %d", url, resp.StatusCode)
	}
	var result model.SearchCodeRes
	json.Unmarshal([]byte(body), &result)
	results := result.Results
	if len(results) == 0 {
		hasResult = false
	}
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
		textMatch := new(model.TextMatch)
		textMatch.Fragment = &lines
		textMatchs := make([]model.TextMatch, 0)
		textMatchs = append(textMatchs, *textMatch)
		b, _ := json.Marshal(textMatchs)
		codeResult := model.SearchResult{
			Path:            val.Filename,
			RepoUrl:         val.Location,
			Status:          0,
			Url:             val.Url,
			Repo:            repoPath,
			TextMatchesJson: b,
		}
		codeResults = append(codeResults, &codeResult)
	}
	return codeResults, hasResult
}

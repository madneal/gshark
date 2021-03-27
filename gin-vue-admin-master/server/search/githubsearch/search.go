package githubsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/models"
	"github.com/madneal/gshark/util"
	"github.com/madneal/gshark/vars"
	"regexp"
	"strings"
	"sync"
	"time"
)

func GenerateSearchCodeTask() (map[int][]models.Rule, error) {
	result := make(map[int][]models.Rule)
	// get rules with the type of github
	rules, err := models.GetValidRulesByType("github")
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

func Search(rules []models.Rule) {
	var wg sync.WaitGroup
	wg.Add(len(rules))
	client, token, err := GetGithubClient()
	var content string
	if err == nil && token != "" {
		for _, rule := range rules {
			go func(rule models.Rule) {
				defer wg.Done()
				results, err := client.SearchCode(rule.Pattern)
				if err != nil {
					logger.Log.Error(err)
					return
				}
				counts := SaveResult(results, &rule.Pattern)
				if counts > 0 {
					content += fmt.Sprintf("%s: %d条\n", rule.Pattern, counts)
				}
			}(rule)
		}
		wg.Wait()
	}
	if vars.SCKEY != "" && content != "" {
		util.SendMessage(vars.SCKEY, "扫描结果", content)
	}
}

func RunSearchTask(mapRules map[int][]models.Rule, err error) {
	if err == nil {
		for _, rules := range mapRules {
			startTime := time.Now()
			Search(rules)
			usedTime := time.Since(startTime).Seconds()
			if usedTime < 60 {
				time.Sleep(time.Duration(60 - usedTime))
			}
		}
	}
}

// The filters are utilized to filter the codeResult
func PassFilters(codeResult *models.CodeResult, fullName string) bool {
	// detect if the Repository url exist in input_info
	repoUrl := codeResult.Repository.GetHTMLURL()

	inputInfo := models.NewInputInfo(CONST_REPO, repoUrl, fullName)
	has, err := inputInfo.Exist()
	if err != nil {
		fmt.Print(err)
		return false
	}
	if !has {
		_, err = inputInfo.Insert()
		if err != nil {
			logger.Log.Error(err)
		}
	}
	// detect if the codeResult exist
	exist, err := codeResult.Exist()
	// detect if there are any random characters in text matches
	textMatches := codeResult.TextMatches[0].Fragment
	reg := regexp.MustCompile(`[A-Za-z0-9_+]{50,}`)
	return !reg.MatchString(*textMatches) && !has && !exist
}

func SaveResult(results []*github.CodeSearchResult, keyword *string) int {
	insertCount := 0
	for _, result := range results {
		if result != nil && len(result.CodeResults) > 0 {
			for _, resultItem := range result.CodeResults {
				ret, err := json.Marshal(resultItem)
				if err == nil {
					var codeResult *models.CodeResult
					err = json.Unmarshal(ret, &codeResult)
					codeResult.Keyword = keyword
					fullName := codeResult.Repository.GetFullName()
					codeResult.RepoName = fullName
					if len(codeResult.TextMatches) > 0 {
						hash := util.GenMd5WithSpecificLen(*(codeResult.TextMatches[0].Fragment), 50)
						codeResult.Textmatchmd5 = &hash
					}

					if err == nil && PassFilters(codeResult, fullName) {
						insertCount++
						logger.Log.Infoln(codeResult.Insert())
					}
				}
			}
		}
		logger.Log.Infof("Has inserted %d results into code_result", insertCount)
	}
	return insertCount
}

func RunTask(duration time.Duration) {
	RunSearchTask(GenerateSearchCodeTask())

	// insert repos from inputInfo
	InsertAllRepos()

	logger.Log.Infof("Complete the scan of Github, start to sleep %v seconds", duration*time.Second)
	time.Sleep(duration * time.Second)
}

func (c *Client) SearchCode(keyword string) ([]*github.CodeSearchResult, error) {
	var allSearchResult []*github.CodeSearchResult
	var err error
	ctx := context.Background()
	listOpt := github.ListOptions{PerPage: 100}
	opt := &github.SearchOptions{Sort: "indexed", Order: "desc", TextMatch: true, ListOptions: listOpt}
	query := keyword + " +in:file"
	query, err = BuildQuery(query)
	fmt.Println("search with the query:" + query)
	for {
		result, nextPage := searchCodeByOpt(c, ctx, query, *opt)
		time.Sleep(time.Second * 3)
		if result != nil {
			allSearchResult = append(allSearchResult, result)
		}
		if nextPage <= 0 {
			break
		}
		opt.Page = nextPage
	}
	return allSearchResult, err
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
			} else if ruleKey == "lang" {
				str += "language:"
			}

			value = strings.TrimSpace(value)
			str += value
		}
	}
	builtQuery := query + str
	return builtQuery, err
}

func searchCodeByOpt(c *Client, ctx context.Context, query string, opt github.SearchOptions) (*github.CodeSearchResult, int) {
	result, res, err := c.Client.Search.Code(ctx, query, &opt)

	if res != nil && res.Rate.Remaining < 10 {
		time.Sleep(45 * time.Second)
	}

	if err == nil {
		logger.Log.Infof("remaining: %d, nextPage: %d, lastPage: %d", res.Rate.Remaining, res.NextPage, res.LastPage)
	} else {
		logger.Log.Error(err)
		return nil, 0
	}
	return result, res.NextPage
}

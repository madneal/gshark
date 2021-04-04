package githubsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/service"
	"github.com/madneal/gshark/utils"
	"go.uber.org/zap"

	//"github.com/madneal/gshark/logger"
	//"github.com/madneal/gshark/model"
	//"github.com/madneal/gshark/util"
	//"github.com/madneal/gshark/vars"
	"github.com/madneal/gshark/model"
	//"github.com/madneal/gshark/service"
	"regexp"
	"sync"
	"time"
)

const SearchNum = 10

func GenerateSearchCodeTask() (map[int][]model.Rule, error) {
	result := make(map[int][]model.Rule)
	// get rules with the type of github
	err, rules := service.GetValidRulesByType("github")
	if len(rules) == 0 {
		global.GVA_LOG.Info("Rules of github is empty, please specify one rule for scan at least")
	}
	ruleNum := len(rules)
	batch := ruleNum / SearchNum

	for i := 0; i < batch; i++ {
		result[i] = rules[SearchNum*i : SearchNum*(i+1)]
	}

	if ruleNum%SearchNum != 0 {
		result[batch] = rules[SearchNum*batch : ruleNum]
	}
	return result, err
}

func Search(rules []model.Rule) {
	var wg sync.WaitGroup
	wg.Add(len(rules))
	client, token, err := GetGithubClient()
	var content string
	if err == nil && token != "" {
		for _, rule := range rules {
			go func(rule model.Rule) {
				defer wg.Done()
				results, err := client.SearchCode(rule.Content)
				if err != nil {
					//logger.Log.Error(err)
					return
				}
				counts := SaveResult(results, &rule.Content)
				if counts > 0 {
					content += fmt.Sprintf("%s: %d条\n", rule.Content, counts)
				}
			}(rule)
		}
		wg.Wait()
	}
	//if global.GVA_CONFIG.Serverj.SCKEY != "" && content != "" {
	//	util.SendMessage(global.GVA_CONFIG.Serverj.SCKEY, "扫描结果", content)
	//}
}

func RunSearchTask(mapRules map[int][]model.Rule, err error) {
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
func PassFilters(codeResult *model.SearchResult, fullName string) bool {
	// detect if the Repository url exist in input_info
	//repoUrl := codeResult.Repository.GetHTMLURL()
	//
	//inputInfo := model.NewInputInfo(CONST_REPO, repoUrl, fullName)
	//has, err := inputInfo.Exist()
	//if err != nil {
	//	return false
	//}
	//if !has {
	//	_, err = inputInfo.Insert()
	//	if err != nil {
	//		logger.Log.Error(err)
	//	}
	//}
	// detect if the codeResult exist
	_, exist := service.CheckExistOfSearchResult(codeResult)
	// detect if there are any random characters in text matches
	textMatches := codeResult.TextMatches[0].Fragment
	reg := regexp.MustCompile(`[A-Za-z0-9_+]{50,}`)
	return !reg.MatchString(*textMatches) && !exist
}

func SaveResult(results []*github.CodeSearchResult, keyword *string) int {
	insertCount := 0
	for _, result := range results {
		if result != nil && len(result.CodeResults) > 0 {
			for _, resultItem := range result.CodeResults {
				ret, err := json.Marshal(resultItem)
				if err == nil {
					var codeResult *model.SearchResult
					err = json.Unmarshal(ret, &codeResult)
					codeResult.Keyword = *keyword
					fullName := codeResult.Repository.GetFullName()
					codeResult.Repo = fullName
					if len(codeResult.TextMatches) > 0 {
						hash := utils.GenMd5WithSpecificLen(*(codeResult.TextMatches[0].Fragment), 50)
						codeResult.TextmatchMd5 = hash
					}

					if err == nil && PassFilters(codeResult, fullName) {
						insertCount++
						//logger.Log.Infoln(codeResult.Insert())
					}
				}
			}
		}
		//logger.Log.Infof("Has inserted %d results into code_result", insertCount)
	}
	return insertCount
}

func RunTask(duration time.Duration) {
	RunSearchTask(GenerateSearchCodeTask())

	// insert repos from inputInfo
	//InsertAllRepos()

	//logger.Log.Infof("Complete the scan of Github, start to sleep %v seconds", duration*time.Second)
	time.Sleep(duration * time.Second)
}

func (c *Client) SearchCode(keyword string) ([]*github.CodeSearchResult, error) {
	var allSearchResult []*github.CodeSearchResult
	var err error
	ctx := context.Background()
	listOpt := github.ListOptions{PerPage: 100}
	opt := &github.SearchOptions{Sort: "indexed", Order: "desc", TextMatch: true, ListOptions: listOpt}
	query := keyword + " +in:file"
	//query, err = BuildQuery(query)
	global.GVA_LOG.Info("Github scan with the query:", zap.Any("github", query))
	//fmt.Println("search with the query:" + query)
	for {
		result, nextPage := searchCodeByOpt(c, ctx, query, *opt)
		//time.Sleep(time.Second * 3)
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

//func BuildQuery(query string) (string, error) {
//	filterRules, err := model.GetFilterRules()
//	str := ""
//	for _, filterRule := range filterRules {
//		ruleValue := filterRule.RuleValue
//		ruleType := filterRule.RuleType
//		ruleKey := filterRule.RuleKey
//		ruleValueList := strings.Split(ruleValue, ",")
//		for _, value := range ruleValueList {
//			if ruleType == 0 {
//				str += " -"
//			} else {
//				str += " +"
//			}
//
//			if ruleKey == "ext" {
//				str += "extension:"
//			} else if ruleKey == "lang" {
//				str += "language:"
//			}
//
//			value = strings.TrimSpace(value)
//			str += value
//		}
//	}
//	builtQuery := query + str
//	return builtQuery, err
//}

func searchCodeByOpt(c *Client, ctx context.Context, query string, opt github.SearchOptions) (*github.CodeSearchResult, int) {
	result, res, err := c.Client.Search.Code(ctx, query, &opt)

	// for best guidelines, wait one second
	// https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits

	time.Sleep(1 * time.Second)
	if res != nil && res.Rate.Remaining < 10 {
		time.Sleep(45 * time.Second)
	}

	if err == nil {
		global.GVA_LOG.Info("Search for "+query, zap.Any("remaining", res.Rate.Remaining), zap.Any("nextPage",
			res.NextPage), zap.Any("lastPage", res.LastPage))
	} else {
		global.GVA_LOG.Error("Search error", zap.Any("github search error", err))
		//if errors.Is(err, github.AbuseRateLimitError)
		time.Sleep(30 * time.Second)
		return nil, 0
	}
	return result, res.NextPage
}

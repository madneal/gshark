package appsearch

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/models"
	"github.com/madneal/gshark/vars"
	"strconv"
	"strings"
	"time"
)

func RunTask(duration time.Duration) {
	RunSearchTask(GenerateSearchCodeTask())
	logger.Log.Infof("Complete the scan of APP, start to sleep %v seconds", duration*time.Second)
	time.Sleep(duration * time.Second)
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
	if err == nil {
		for _, rules := range mapRules {
			for _, rule := range rules {
				results := SearchForApp(rule)
				SaveResults(results)
			}
		}
	}
}

func SaveResults(results []*models.AppSearchResult) {
	for _, result := range results {
		isExist, err := result.Exist()
		if err != nil {
			fmt.Println(err)
		}
		if !isExist {
			_, err := result.Insert()
			if err != nil {
				fmt.Println(err)
			}
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SearchForApp(rule models.Rule) []*models.AppSearchResult {
	appSearchResults := make([]*models.AppSearchResult, 0)
	var hasNext bool
	if rule.Caption == "HUAWEI" {
		baseUrl := "http://appstore.huawei.com"
		url := baseUrl + "/search/" + rule.Pattern
		for i := 1; i < 99; i++ {
			c := colly.NewCollector()
			linkUrl := url + "/" + strconv.Itoa(i)
			c.OnHTML("body", func(e *colly.HTMLElement) {
				var results []*models.AppSearchResult
				hasNext, results = saveAppSearchResult(baseUrl, e)
				if hasNext {
					appSearchResults = append(appSearchResults, results...)
				} else {
					return
				}
				fmt.Println(linkUrl)
			})
			c.Visit(linkUrl)
			fmt.Println("the length of appSearchResults")
			fmt.Println(len(appSearchResults))
		}
		// todo
		// other app market
	} else {

	}
	return appSearchResults
}

func saveAppSearchResult(baseUrl string, e *colly.HTMLElement) (bool, []*models.AppSearchResult) {
	var hasNext bool
	appSearchResults := make([]*models.AppSearchResult, 0)
	e.ForEach(".list-game-app.dotline-btn.nofloat", func(i int, element *colly.HTMLElement) {
		appSearchResult := new(models.AppSearchResult)
		var title = element.ChildText(".title")
		var content = element.ChildText(".content")
		var deployDate = strings.Replace(element.ChildText(".date"),
			"发布时间： ", "", -1)
		var appUrl = baseUrl + element.ChildAttr(".title a", "href")
		var market = "HUAWEI"
		appSearchResult.Name = &title
		appSearchResult.Description = &content
		appSearchResult.DeployDate = &deployDate
		appSearchResult.AppUrl = &appUrl
		appSearchResult.Market = &market
		appSearchResult.Status = 0
		fmt.Println(*appSearchResult.Name)
		fmt.Println(*appSearchResult.Description)
		fmt.Println(*appSearchResult.DeployDate)
		fmt.Println(*appSearchResult.AppUrl)
		appSearchResults = append(appSearchResults, appSearchResult)
	})
	fmt.Println(e.DOM.Find(".list-game-app.dotline-btn.nofloat").Length())
	if e.DOM.Find(".list-game-app.dotline-btn.nofloat").Length() > 1 {
		hasNext = true
	} else {
		hasNext = false
	}
	return hasNext, appSearchResults
}

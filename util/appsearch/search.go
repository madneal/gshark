package appsearch

import (
	"github.com/neal1991/gshark/logger"
	"github.com/neal1991/gshark/models"
	"github.com/neal1991/gshark/vars"
	"time"
	"github.com/gocolly/colly"
)

func ScheduleTasks(duration time.Duration) {
	for {
		logger.Log.Infof("Complete the scan of APP, start to sleep %v seconds", duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}

func GenerateSearchCodeTask() (map[int][]models.Rule, error) {
	result := make(map[int][]models.Rule)
	rules, err := models.GetValidRules()
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

func SearchForApp(rule models.Rule)  {
	appSearchResult := new(models.APPSearchResult)
	if rule.Caption == "HUAWEI" {
		url := "http://appstore.huawei.com/search/"
		url = url + rule.Pattern
		c := colly.NewCollector()
		c.OnHTML("body", func(e *colly.HTMLElement) {
			e.ForEach(".list-game-app", func(i int, element *colly.HTMLElement) {
				*appSearchResult.Name = element.ChildText(".title")
				*appSearchResult.Description = element.ChildText(".content")
				*appSearchResult.DeployDate = element.ChildText("")
			})
		})
		c.Visit(url)
	} else {

	}
}

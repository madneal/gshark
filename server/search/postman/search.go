package postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

var postmanUrl = "https://www.postman.com/_api/ws/proxy"

type Document struct {
	Summary            string        `json:"summary"`
	RequestCount       int           `json:"requestCount"`
	PublisherType      string        `json:"publisherType"`
	Imports            int           `json:"imports,omitempty"`
	WatcherCount       int           `json:"watcherCount"`
	EntityType         string        `json:"entityType"`
	ForkCount          int           `json:"forkCount"`
	Tags               []interface{} `json:"tags"`
	Quality            int           `json:"quality,omitempty"`
	PublisherId        string        `json:"publisherId"`
	ForkLabel          string        `json:"forkLabel"`
	Apis               []interface{} `json:"apis,omitempty"`
	PublisherHandle    string        `json:"publisherHandle"`
	PublisherName      string        `json:"publisherName"`
	PublisherLogo      string        `json:"publisherLogo"`
	IsDomainNonTrivial bool          `json:"isDomainNonTrivial"`
	Name               string        `json:"name"`
	IsPublic           bool          `json:"isPublic"`
	Workspaces         []struct {
		VisibilityStatus string `json:"visibilityStatus"`
		Name             string `json:"name"`
		Id               string `json:"id"`
		Slug             string `json:"slug"`
	} `json:"workspaces"`
	Id           string        `json:"id"`
	Categories   []interface{} `json:"categories"`
	Views        int           `json:"views"`
	DocumentType string        `json:"documentType"`
}

type Requests struct {
	Score    float64 `json:"score"`
	Document struct {
		Method string `json:"method"`
		Name   string `json:"name"`
		Id     string `json:"id"`
		Url    string `json:"url"`
	} `json:"document"`
}

type PostmanRes struct {
	Data []struct {
		Score           float64   `json:"score"`
		NormalizedScore float64   `json:"normalizedScore"`
		Document        Document  `json:"document"`
		Requests        *Requests `json:"requests"`
	} `json:"data"`
	Meta struct {
		QueryText string `json:"queryText"`
		Total     struct {
			Collection int `json:"collection"`
			Workspace  int `json:"workspace"`
			Api        int `json:"api"`
			Team       int `json:"team"`
			User       int `json:"user"`
			Request    int `json:"request"`
		} `json:"total"`
		State              string `json:"state"`
		CorrectedQueryText string `json:"correctedQueryText"`
	} `json:"meta"`
}

func RunTask() {
	err, rules := service.GetValidRulesByType("postman")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType postman err", zap.Error(err))
		return
	}
	color.Infoln("begin the postman search task")
	Search(&rules)
	color.Infof("finish the postman search task, ready to sleep\n")
	time.Sleep(900 * time.Second)
}

func Search(rules *[]model.Rule) {
	postmanClient := GetPostmanClient()
	for _, rule := range *rules {
		//postmanClient.SearchByType(rule.Content, "collection")
		postmanClient.SearchByType(rule.Content, "request")
	}
}

func (postmanClient *Client) SearchByType(keyword, searchType string) {
	resList, err := postmanClient.SearchAPI(keyword, searchType)
	if err != nil {
		global.GVA_LOG.Error("postman SearchAPI err", zap.Error(err))
		return
	}
	for _, res := range *resList {
		results := res.CovertToSearchResult(keyword)
		for _, result := range *results {
			err = service.CreateSearchResult(result)
			if err != nil {
				global.GVA_LOG.Error("CreateSearchResult err", zap.Error(err))
			}
		}
	}
}

type Client struct {
	client *http.Client
	sid    string
}

func (res *PostmanRes) CovertToSearchResult(keyword string) *[]model.SearchResult {
	results := make([]model.SearchResult, 0)
	for _, data := range res.Data {
		document := data.Document
		var requestURL string
		if document.DocumentType == "collection" {
			requestURL = fmt.Sprintf("https://www.postman.com/workspace/collection/%s", document.Id)
		}
		requests := data.Requests
		matches := document.PublisherName + " | " + document.Summary
		if requests != nil {
			matches = requests.Document.Name + " | " + requests.Document.Url
		}
		fmt.Println(matches)
		result := model.SearchResult{
			Path:    document.PublisherName,
			Url:     requestURL,
			Matches: matches,
			Keyword: keyword,
			Repo:    document.PublisherName + "/" + document.Name,
		}
		results = append(results, result)
	}
	return &results
}

func (client *Client) SearchAPI(rule, searchType string) (*[]PostmanRes, error) {
	page := 0
	resList := make([]PostmanRes, 0)
	var err error
	for {
		color.Infof("search for the rule %s of page %d\n", rule, page)
		body := fmt.Sprintf(`{"service":"search","method":"POST","path":"/search-all","body":{"queryIndices":["runtime.%s"],"queryText":"%s","size":100,"from": %d, "mergeEntities":true}}`,
			searchType, rule, page)
		req, err := http.NewRequest("POST", postmanUrl, bytes.NewBufferString(body))
		req.Header.Set("Cookie", "postman.sid="+client.sid)
		req.Header.Set("Host", "www.postman.com")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
		if err != nil {
			return &resList, err
		}
		httpClient := http.Client{}
		res, err := httpClient.Do(req)

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			global.GVA_LOG.Error("postman ReadAll err", zap.Error(err))
			return &resList, err
		}
		var postRes PostmanRes
		if err = json.Unmarshal(resBody, &postRes); err != nil {
			return &resList, err
		}
		resList = append(resList, postRes)
		page = page + 1
		var total float64
		if searchType == "collection" {
			total = float64(postRes.Meta.Total.Collection)
		} else if searchType == "request" {
			total = float64(postRes.Meta.Total.Request)
		}
		if float64(page) > math.Ceil(total/100) {
			break
		}
	}
	return &resList, err
}

func GetPostmanClient() *Client {
	var sid string
	err, tokens := service.ListTokenByType("postman")
	if err != nil {
		global.GVA_LOG.Error("ListTokenByType postman err", zap.Error(err))
		return nil
	}
	if len(tokens) > 0 {
		sid = tokens[0].Content
	}

	client := Client{
		sid: sid,
	}
	return &client
}

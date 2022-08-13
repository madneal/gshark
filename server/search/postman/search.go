package postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var postmanUrl = "https://bifrost-web-https-v4.gw.postman.com/ws/proxy"

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
		Score           float64  `json:"score"`
		NormalizedScore int      `json:"normalizedScore"`
		Document        Document `json:"document"`
		Requests        Requests `json:"requests"`
	} `json:"data"`
	Meta struct {
		QueryText string `json:"queryText"`
		Total     struct {
			Collection int `json:"collection"`
			Workspace  int `json:"workspace"`
			Api        int `json:"api"`
			Team       int `json:"team"`
			User       int `json:"user"`
			Flow       int `json:"flow"`
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
	Search(&rules)
}

func Search(rules *[]model.Rule) {
	postmanClient := GetPostmanClient()
	for _, rule := range *rules {
		resList, err := postmanClient.SearchAPI(rule.Content)
		if err != nil {
			global.GVA_LOG.Error("postman SearchAPI err", zap.Error(err))
			return
		}
		for _, res := range *resList {
			results := res.CovertToSearchResult()
			for _, result := range *results {
				err = service.CreateSearchResult(result)
				if err != nil {
					global.GVA_LOG.Error("CreateSearchResult err", zap.Error(err))
				}
			}
		}
	}
}

type Client struct {
	client *http.Client
	sid    string
}

func (res *PostmanRes) CovertToSearchResult() *[]model.SearchResult {
	results := make([]model.SearchResult, 0)
	for _, data := range res.Data {
		document := data.Document
		requestURL := fmt.Sprintf("https://www.postman.com/%s/workspace/%s/request/%s", document.PublisherHandle,
			document.Workspaces[0].Slug, data.Requests.Document.Id)
		result := model.SearchResult{
			Path:    data.Requests.Document.Name + "/" + data.Requests.Document.Name,
			Url:     requestURL,
			Matches: data.Requests.Document.Url,
		}
		results = append(results, result)
	}
	return &results
}

func (client *Client) SearchAPI(rule string) (*[]PostmanRes, error) {
	page := 0
	resList := make([]PostmanRes, 0)
	var err error
	for {
		body := fmt.Sprintf(`{"service":"search","method":"POST","path":"/search-all",
"body":{"queryIndices":["runtime.collection"],"queryText":"%s","size":100,"from": %d, "mergeEntities":true}}`,
			rule, page)
		res, err := client.client.Post(postmanUrl, "application/json", bytes.NewBufferString(body))
		if err != nil {
			return &resList, err
		}

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
		total := float64(postRes.Meta.Total.Collection)
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
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse("https://bifrost-web-https-v4.gw.postman.com")
	cookies := []*http.Cookie{{
		Name:  "postman.sid",
		Value: sid,
	}}
	jar.SetCookies(u, cookies)
	httpClient := http.Client{
		Jar: jar,
	}

	client := Client{
		sid:    sid,
		client: &httpClient,
	}
	return &client
}

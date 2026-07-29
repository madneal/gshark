package postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

const (
	postmanURL             = "https://www.postman.com/_api/ws/proxy"
	postmanPageSize        = 25
	postmanRequestTimeout  = 30 * time.Second
	maxPostmanResponseSize = 10 << 20
)

var postmanHTTPClient = &http.Client{Timeout: postmanRequestTimeout}

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
	Method             string        `json:"method"`
	URL                string        `json:"url"`
	WorkspaceSlug      string        `json:"workspaceSlug"`
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
	for _, rule := range *rules {
		SearchByType(rule.Content, "request")
	}
}

func SearchByType(keyword, searchType string) {
	resList, err := SearchAPI(keyword, searchType)
	for _, res := range *resList {
		results := res.ConvertToSearchResult(keyword)
		stats := service.SaveSearchResultsWithStats(*results)
		global.GVA_LOG.Info(stats.Summary(keyword, "Postman"))
	}
	if err != nil {
		global.GVA_LOG.Error("postman SearchAPI err", zap.Error(err))
	}
}

func (res *PostmanRes) ConvertToSearchResult(keyword string) *[]model.SearchResult {
	results := make([]model.SearchResult, 0)
	for _, data := range res.Data {
		document := data.Document
		name := document.Name
		method := document.Method
		requestValue := document.URL
		requests := data.Requests
		if requests != nil {
			if name == "" {
				name = requests.Document.Name
			}
			if method == "" {
				method = requests.Document.Method
			}
			if requestValue == "" {
				requestValue = requests.Document.Url
			}
		}
		matches := joinNonEmpty(method, name, requestValue, document.Summary)
		result := model.SearchResult{
			Path:    document.PublisherName,
			Url:     buildPostmanURL(document),
			Matches: matches,
			Keyword: keyword,
			Repo:    document.PublisherName + "/" + name,
		}
		results = append(results, result)
	}
	return &results
}

func SearchAPI(rule, searchType string) (*[]PostmanRes, error) {
	return searchAPI(postmanHTTPClient, postmanURL, rule, searchType)
}

type postmanSearchRequest struct {
	Service string                   `json:"service"`
	Method  string                   `json:"method"`
	Path    string                   `json:"path"`
	Body    postmanSearchRequestBody `json:"body"`
}

type postmanSearchRequestBody struct {
	QueryIndices  []string `json:"queryIndices"`
	QueryText     string   `json:"queryText"`
	Size          int      `json:"size"`
	From          int      `json:"from"`
	MergeEntities bool     `json:"mergeEntities"`
}

func searchAPI(client *http.Client, endpoint, rule, searchType string) (*[]PostmanRes, error) {
	if searchType != "collection" && searchType != "request" {
		return &[]PostmanRes{}, fmt.Errorf("unsupported Postman search type %q", searchType)
	}

	resList := make([]PostmanRes, 0)
	for page, offset := 0, 0; ; page, offset = page+1, offset+postmanPageSize {
		color.Infof("search for the rule %s of page %d\n", rule, page)
		payload := postmanSearchRequest{
			Service: "search",
			Method:  http.MethodPost,
			Path:    "/search-all",
			Body: postmanSearchRequestBody{
				QueryIndices:  []string{"runtime." + searchType},
				QueryText:     rule,
				Size:          postmanPageSize,
				From:          offset,
				MergeEntities: true,
			},
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return &resList, fmt.Errorf("marshal Postman search page %d: %w", page, err)
		}

		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return &resList, fmt.Errorf("create Postman search request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")

		res, err := client.Do(req)
		if err != nil {
			return &resList, fmt.Errorf("request Postman search page %d: %w", page, err)
		}

		resBody, err := readPostmanResponse(res)
		if err != nil {
			return &resList, fmt.Errorf("read Postman search page %d: %w", page, err)
		}
		var postRes PostmanRes
		if err = json.Unmarshal(resBody, &postRes); err != nil {
			return &resList, fmt.Errorf("decode Postman search page %d: %w", page, err)
		}
		resList = append(resList, postRes)

		total := postRes.Meta.Total.Request
		if searchType == "collection" {
			total = postRes.Meta.Total.Collection
		}
		if len(postRes.Data) == 0 || offset+postmanPageSize >= total {
			break
		}
	}
	return &resList, nil
}

func readPostmanResponse(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, maxPostmanResponseSize+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxPostmanResponseSize {
		return nil, fmt.Errorf("response exceeded %d bytes", maxPostmanResponseSize)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		message := strings.TrimSpace(string(body))
		if len(message) > 1024 {
			message = message[:1024]
		}
		return nil, fmt.Errorf("unexpected HTTP status %d: %s", res.StatusCode, message)
	}
	return body, nil
}

func buildPostmanURL(document Document) string {
	workspaceSlug := document.WorkspaceSlug
	if workspaceSlug == "" && len(document.Workspaces) > 0 {
		workspaceSlug = document.Workspaces[0].Slug
	}
	if document.PublisherHandle != "" && workspaceSlug != "" && document.Id != "" {
		return fmt.Sprintf(
			"https://www.postman.com/%s/workspace/%s/%s/%s",
			url.PathEscape(document.PublisherHandle),
			url.PathEscape(workspaceSlug),
			url.PathEscape(document.DocumentType),
			url.PathEscape(document.Id),
		)
	}
	if document.DocumentType == "collection" && document.Id != "" {
		return fmt.Sprintf("https://www.postman.com/workspace/collection/%s", url.PathEscape(document.Id))
	}
	return ""
}

func joinNonEmpty(values ...string) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			parts = append(parts, value)
		}
	}
	return strings.Join(parts, " | ")
}

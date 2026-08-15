package postman

import (
	"bytes"
	"encoding/json"
	"errors"
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
	postmanURL             = "https://api.postman.com/search"
	postmanPageSize        = 25
	postmanRequestTimeout  = 30 * time.Second
	maxPostmanResponseSize = 10 << 20
)

var postmanHTTPClient = &http.Client{Timeout: postmanRequestTimeout}

var (
	getPostmanRules = service.GetValidRulesByType
	searchPostman   = Search
)

type PostmanResourceRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PostmanOrganization struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsVerified bool   `json:"isVerified"`
}

type PostmanLinks struct {
	Web struct {
		Href string `json:"href"`
	} `json:"web"`
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

type PostmanResource struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Method       string              `json:"method"`
	Description  string              `json:"description"`
	URL          string              `json:"url"`
	Collection   PostmanResourceRef  `json:"collection"`
	Workspace    PostmanResourceRef  `json:"workspace"`
	Organization PostmanOrganization `json:"organization"`
	Team         PostmanResourceRef  `json:"team"`
	Links        PostmanLinks        `json:"links"`
}

type PostmanRes struct {
	Data []PostmanResource `json:"data"`
	Meta struct {
		QueryText  string `json:"q"`
		Total      int    `json:"total"`
		NextCursor string `json:"nextCursor"`
	} `json:"meta"`
}

func RunTask() model.ScanOutcome {
	err, rules := getPostmanRules("postman")
	if err != nil {
		global.GVA_LOG.Error("GetValidRulesByType postman err", zap.Error(err))
		return model.ScanFailed("Failed to load Postman rules: " + err.Error())
	}
	if len(rules) == 0 {
		message := "No enabled Postman rules; provider skipped"
		global.GVA_LOG.Info(message)
		return model.ScanSkipped(message)
	}
	color.Infoln("begin the postman search task")
	if err := searchPostman(&rules); err != nil {
		return model.ScanFailed("Postman scan completed with errors: " + err.Error())
	}
	color.Infof("finish the postman search task\n")
	return model.ScanSuccess(fmt.Sprintf("Completed %d Postman rules", len(rules)))
}

func Search(rules *[]model.Rule) error {
	var searchErrors []error
	for _, rule := range *rules {
		if err := SearchByType(rule.Content, "request"); err != nil {
			searchErrors = append(searchErrors, err)
		}
	}
	return errors.Join(searchErrors...)
}

func SearchByType(keyword, searchType string) error {
	err := SearchAPIStream(keyword, searchType, func(res PostmanRes) error {
		results := res.ConvertToSearchResult(keyword)
		stats := service.SaveSearchResultsWithStats(*results)
		global.GVA_LOG.Info(stats.Summary(keyword, "Postman"))
		return nil
	})
	if err != nil {
		global.GVA_LOG.Error("postman SearchAPI err", zap.Error(err))
	}
	return err
}

func (res *PostmanRes) ConvertToSearchResult(keyword string) *[]model.SearchResult {
	results := make([]model.SearchResult, 0)
	for _, resource := range res.Data {
		result := model.SearchResult{
			Path:    resource.ID,
			RepoUrl: resource.Links.Web.Href,
			Url:     resource.Links.Web.Href,
			Matches: joinNonEmpty(resource.Method, resource.Name, resource.URL, resource.Description),
			Keyword: keyword,
			Repo:    buildPostmanRepo(resource),
		}
		results = append(results, result)
	}
	return &results
}

func SearchAPI(rule, searchType string) (*[]PostmanRes, error) {
	return searchAPI(postmanHTTPClient, postmanURL, rule, searchType)
}

func SearchAPIStream(rule, searchType string, onPage func(PostmanRes) error) error {
	return searchAPIStream(postmanHTTPClient, postmanURL, rule, searchType, onPage)
}

type postmanSearchRequest struct {
	ElementType string               `json:"elementType"`
	QueryText   string               `json:"q"`
	Ownership   string               `json:"ownership"`
	Filters     postmanSearchFilters `json:"filters"`
}

type postmanSearchFilters struct {
	And []postmanSearchFilter `json:"$and"`
}

type postmanSearchFilter struct {
	Visibility postmanEqualsFilter `json:"visibility"`
}

type postmanEqualsFilter struct {
	Equal string `json:"$eq"`
}

func searchAPI(client *http.Client, endpoint, rule, searchType string) (*[]PostmanRes, error) {
	resList := make([]PostmanRes, 0)
	err := searchAPIStream(client, endpoint, rule, searchType, func(page PostmanRes) error {
		resList = append(resList, page)
		return nil
	})
	return &resList, err
}

func searchAPIStream(client *http.Client, endpoint, rule, searchType string, onPage func(PostmanRes) error) error {
	if searchType != "collection" && searchType != "request" {
		return fmt.Errorf("unsupported Postman search type %q", searchType)
	}

	body, err := json.Marshal(postmanSearchRequest{
		ElementType: searchType + "s",
		QueryText:   rule,
		Ownership:   "all",
		Filters: postmanSearchFilters{
			And: []postmanSearchFilter{{
				Visibility: postmanEqualsFilter{Equal: "public"},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("marshal Postman search request: %w", err)
	}

	seenCursors := make(map[string]struct{})
	seenResources := make(map[string]struct{})
	cursor := ""
	for page := 0; ; page++ {
		color.Infof("search for the rule %s of page %d\n", rule, page)
		requestURL, err := buildPostmanSearchURL(endpoint, cursor)
		if err != nil {
			return fmt.Errorf("build Postman search page %d URL: %w", page, err)
		}

		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("create Postman search request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")

		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("request Postman search page %d: %w", page, err)
		}

		resBody, err := readPostmanResponse(res)
		if err != nil {
			return fmt.Errorf("read Postman search page %d: %w", page, err)
		}
		var postRes PostmanRes
		if err = json.Unmarshal(resBody, &postRes); err != nil {
			return fmt.Errorf("decode Postman search page %d: %w", page, err)
		}

		uniqueResources := make([]PostmanResource, 0, len(postRes.Data))
		for _, resource := range postRes.Data {
			key := postmanResourceKey(resource)
			if key != "" {
				if _, exists := seenResources[key]; exists {
					continue
				}
				seenResources[key] = struct{}{}
			}
			uniqueResources = append(uniqueResources, resource)
		}
		postRes.Data = uniqueResources
		if len(postRes.Data) > 0 {
			if err := onPage(postRes); err != nil {
				return err
			}
		}

		nextCursor := postRes.Meta.NextCursor
		if nextCursor == "" {
			break
		}
		if _, exists := seenCursors[nextCursor]; exists {
			return fmt.Errorf("Postman search repeated cursor on page %d", page)
		}
		seenCursors[nextCursor] = struct{}{}
		cursor = nextCursor
	}
	return nil
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

func buildPostmanSearchURL(endpoint, cursor string) (string, error) {
	searchURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	query := searchURL.Query()
	query.Set("limit", fmt.Sprintf("%d", postmanPageSize))
	if cursor != "" {
		query.Set("cursor", cursor)
	}
	searchURL.RawQuery = query.Encode()
	return searchURL.String(), nil
}

func postmanResourceKey(resource PostmanResource) string {
	if resource.ID != "" {
		return resource.ID
	}
	if resource.Links.Web.Href != "" {
		return resource.Links.Web.Href
	}
	return ""
}

func buildPostmanRepo(resource PostmanResource) string {
	owner := resource.Organization.Name
	if owner == "" {
		owner = resource.Workspace.Name
	}
	label := joinNonEmpty(owner, resource.Collection.Name, resource.Name)
	if resource.ID == "" {
		return label
	}
	return label + " [" + resource.ID + "]"
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

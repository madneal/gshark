package codesearch

import (
	"github.com/madneal/gshark/models"
	"github.com/parnurzeal/gorequest"
	"testing"
)

func TestSearchForSearchCode(t *testing.T) {
	rule := new(models.Rule)
	rule.Pattern = "SPDB"
	request := gorequest.New()
	SearchForSearchCode(*rule, request)
}

func TestGetResult(t *testing.T) {
	url := "https://searchcode.com/api/codesearch_I/?q=spdb&p=0&per_page=100"
	request := gorequest.New()
	GetResult(request, url)
}

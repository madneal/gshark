package githubsearch_test

import (
	"fmt"
	"testing"
	"x-patrol/util/githubsearch"
)

func TestSearchCode(t *testing.T) {
	gitClient, _, _ := githubsearch.GetGithubClient()
	codeSearchResults, _ := gitClient.SearchCode("proxy\\.spdb\\.com")
	for _, codeSearchResult := range codeSearchResults {
		for _, codeResult := range codeSearchResult.CodeResults {
			fmt.Println(codeResult.TextMatches)
			fmt.Println(codeResult.HTMLURL)
		}
	}
}

func TestBuildQuery(t *testing.T) {
	query := "shang"
	buildedQuery, err := githubsearch.BuildQuery(query)
	if err == nil {
		fmt.Println(buildedQuery)
	}
}

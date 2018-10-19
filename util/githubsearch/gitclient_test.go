package githubsearch

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSearchCode(t *testing.T) {
	gitClient, _, _ := GetGithubClient()
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
	buildedQuery, err := BuildQuery(query)
	if err == nil {
		fmt.Println(buildedQuery)
	}
}

func TestClient_GetUserInfo(t *testing.T) {
	gitClient, _, _ := GetGithubClient()
	user, resp, _ := gitClient.GetUserInfo("neal1991")
	assert.Equal(t, "https://madneal.com", *user.Blog)
	assert.True(t, resp.StatusCode == 200)
}

package githubsearch

import "testing"
import (
	"github.com/stretchr/testify/assert"
)

func TestGetUserInfo(t *testing.T) {
	gitClient, _, _ := GetGithubClient()
	ownerName := "neal1991"
	user := GetGithubUserInfo(gitClient, &ownerName)
	assert.Equal(t, *user.Blog, "https://madneal.com")
}

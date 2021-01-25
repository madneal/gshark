package githubsearch

import (
	"fmt"
	"github.com/madneal/gshark/models"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	gitClient, _, _ := GetGithubClient()
	ownerName := "madneal"
	user := GetGithubUserInfo(gitClient, &ownerName)
	assert.Equal(t, "https://madneal.com", *user.Blog)
}

func TestInsertAllRepos(t *testing.T) {
	assets, err := models.ListInputInfo()
	if err == nil {
		for index, asset := range assets {
			fmt.Println("index:" + strconv.Itoa(index))
			fmt.Println(asset.Url)
			length := len(strings.Split(asset.Url, "/"))
			fmt.Println(strings.Split(asset.Url, "/")[length-1])
		}
	}
}

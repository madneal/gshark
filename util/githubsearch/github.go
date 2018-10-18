package githubsearch

import (
	"github.com/google/go-github/github"
	"strings"
	"x-patrol/logger"
	"x-patrol/models"
)

const (
	CONST_REPO  = "repo"
	CONST_REPOS = "repos"
	CONST_ORGS  = "organizations"
	CONST_USER  = "user"
)

func InsertAllRepos() {
	gitClient, _, _ := GetGithubClient()

	assets, err := models.ListInputInfo()
	if err == nil {
		for _, asset := range assets {
			assetType := strings.ToLower(asset.Type)
			name := asset.Content
			switch assetType {
			case CONST_REPO, CONST_REPOS:
				repos := strings.Split(name, ",")
				for _, item := range repos {
					r := models.NewRepo(item, item, &asset)
					has, err := r.Exist()
					if err == nil && !has {
						r.Insert()
					}
				}

			case CONST_ORGS:
				orgs := strings.Split(name, ",")
				var orgsRepos []*github.Repository
				var usersAll []*github.User
				for _, org := range orgs {
					users, resp, err := gitClient.GetOrgsMembers(org)
					usersAll = append(usersAll, users...)
					logger.Log.Println(users, resp, err)
					repos, resp, err := gitClient.GetOrgsRepos(org)
					orgsRepos = append(orgsRepos, repos...)
					models.UpdateRate(gitClient.Token, resp)
				}
				mapRepos := gitClient.GetUsersRepos(usersAll)
				for _, rs := range mapRepos {
					orgsRepos = append(orgsRepos, rs...)
				}

				for _, repo := range orgsRepos {
					r := models.NewRepo(*repo.Name, *repo.HTMLURL, &asset)
					has, err := r.Exist()
					if err == nil && !has {
						r.Insert()
					}
				}

			case CONST_USER:
				var usersRepos []*github.Repository
				users := strings.Split(name, ",")
				mapRepos := gitClient.GetStrUsersRepos(users)
				for _, rs := range mapRepos {
					usersRepos = append(usersRepos, rs...)
				}
				for _, repo := range usersRepos {
					r := models.NewRepo(*repo.Name, *repo.HTMLURL, &asset)
					has, err := r.Exist()
					if err == nil && !has {
						r.Insert()
					}
				}
			}
		}
	}
}

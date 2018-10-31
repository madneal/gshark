package models

import (
	"github.com/google/go-github/github"
	"time"
)

type GithubToken struct {
	Id    int64
	Token string
	Desc  string
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`
	// The number of remaining requests the client can make this hour.
	Remaining int `xorm:"default 5000 notnull" json:"remaining"`
	// The time at which the current rate limit will reset.
	Reset time.Time `json:"reset"`
}

// create a GithubToken with limit and remain
func NewGithubToken(token, desc string) *GithubToken {
	return &GithubToken{Token: token, Desc: desc, Limit: 5000, Remaining: 5000}
}

// insert a GithubToken into database
func (g *GithubToken) Insert() (int64, error) {
	return Engine.Insert(g)
}

// detect if the GithubToken exists
func (g *GithubToken) Exist() (bool, error) {
	return Engine.Get(g)
}

func ListTokens() ([]GithubToken, error) {
	tokens := make([]GithubToken, 0)
	err := Engine.Find(&tokens)
	return tokens, err
}

func ListValidTokens() ([]GithubToken, error) {
	tokens := make([]GithubToken, 0)
	err := Engine.Table("github_token").Where("remaining>50").Find(&tokens)
	return tokens, err
}

func GetTokenById(id int64) (*GithubToken, bool, error) {
	token := new(GithubToken)
	has, err := Engine.ID(id).Get(token)
	return token, has, err
}

func EditTokenById(id int64, token, desc string) error {
	githubToken := new(GithubToken)
	has, err := Engine.ID(id).Get(githubToken)
	if err == nil && has {
		githubToken.Token = token
		githubToken.Desc = desc
		Engine.ID(id).Update(githubToken)
	}
	return err
}

func DeleteTokenById(id int64) error {
	token := new(GithubToken)
	_, err := Engine.ID(id).Delete(token)
	return err
}

func UpdateRate(token string, response *github.Response) error {
	githubToken := new(GithubToken)
	has, err := Engine.Table("github_token").Where("token=?", token).Get(githubToken)
	if err == nil && has  {
		id := githubToken.Id
		githubToken.Remaining = response.Rate.Remaining
		githubToken.Reset = response.Rate.Reset.Time
		githubToken.Limit = response.Rate.Limit
		Engine.ID(id).Update(githubToken)
	}
	return err
}

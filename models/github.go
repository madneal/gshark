package models

import (
	"github.com/google/go-github/github"
	"time"
)

type GitToken struct {
	Id    int64
	Token string
	Desc  string
	Type  string
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`
	// The number of remaining requests the client can make this hour.
	Remaining int `xorm:"default 5000 notnull" json:"remaining"`
	// The time at which the current rate limit will reset.
	Reset time.Time `json:"reset"`
}

// create a GitToken with limit and remain
func NewGithubToken(token, desc, tokenType string) *GitToken {
	return &GitToken{Token: token, Desc: desc, Limit: 5000, Remaining: 5000, Type: tokenType}
}

// insert a GitToken into database
func (g *GitToken) Insert() (int64, error) {
	return Engine.Insert(g)
}

// detect if the GitToken exists
func (g *GitToken) Exist() (bool, error) {
	return Engine.Get(g)
}

func ListTokens() ([]GitToken, error) {
	tokens := make([]GitToken, 0)
	err := Engine.Find(&tokens)
	return tokens, err
}

func ListValidTokens(tokenType string) ([]GitToken, error) {
	tokens := make([]GitToken, 0)
	err := Engine.Table("git_token").Where("remaining>50 and type = ?", tokenType).Find(&tokens)
	return tokens, err
}

func GetTokenById(id int64) (*GitToken, bool, error) {
	token := new(GitToken)
	has, err := Engine.ID(id).Get(token)
	return token, has, err
}

func EditTokenById(id int64, token, desc, tokenType string) error {
	githubToken := new(GitToken)
	has, err := Engine.ID(id).Get(githubToken)
	if err == nil && has {
		githubToken.Token = token
		githubToken.Desc = desc
		githubToken.Type = tokenType
		Engine.ID(id).Update(githubToken)
	}
	return err
}

func DeleteTokenById(id int64) error {
	token := new(GitToken)
	_, err := Engine.ID(id).Delete(token)
	return err
}

func UpdateRate(token string, response *github.Response) error {
	githubToken := new(GitToken)
	has, err := Engine.Table("git_token").Where("token=? and type = 'github'", token).Get(githubToken)
	if err == nil && has {
		id := githubToken.Id
		githubToken.Remaining = response.Rate.Remaining
		githubToken.Reset = response.Rate.Reset.Time
		githubToken.Limit = response.Rate.Limit
		Engine.ID(id).Update(githubToken)
	}
	return err
}

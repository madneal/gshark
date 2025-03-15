package githubsearch

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/google/go-github/v57/github"
	"github.com/madneal/gshark/global"
	"go.uber.org/zap"
	"net/http"

	"github.com/madneal/gshark/service"
)

var (
	GithubClients map[string]*Client
	GithubClient  *Client
)

type Client struct {
	Client *github.Client
	Token  string
}

func InitGithubClient(token string) *Client {
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: httpTransport}
	gitClient := github.NewClient(httpClient).WithAuthToken(token)
	client := NewGitClient(gitClient, token)
	return client
}

func GetGithubClient() (*Client, error) {
	err, tokens := service.ListTokenByType("github")
	if err != nil {
		return nil, err
	}
	client := InitGithubClient(tokens[0].Content)
	if client == nil {
		err = errors.New("github Client initial failed, please add token")
	}
	return client, err
}

func NewGitClient(GithubClient *github.Client, token string) *Client {
	return &Client{Client: GithubClient, Token: token}
}

func (c *Client) NextClient() *Client {
	currentToken := c.Token
	err, tokens := service.ListTokenByType("github")
	if err != nil {
		global.GVA_LOG.Error("github Client initial failed, please add token", zap.Error(err))
		return nil
	}
	var currentIndex int
	for index, token := range tokens {
		if token.Content == currentToken {
			currentIndex = index
		}
	}
	nextIndex := (currentIndex + 1) % len(tokens)
	nextToken := tokens[nextIndex]
	return InitGithubClient(nextToken.Content)
}

func (c *Client) GetUserInfo(username string) (*github.User, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Users.Get(ctx, username)
}

func (c *Client) GetOrgsMembers(org string) ([]*github.User, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Organizations.ListMembers(ctx, org, nil)
}

func (c *Client) GetOrgsRepos(org string) ([]*github.Repository, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Repositories.ListByOrg(ctx, org, nil)
}

func (c *Client) GetUserRepos(username string) ([]*github.Repository, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Repositories.List(ctx, username, nil)
}

func (c *Client) GetUserOrgs(username string) ([]*github.Organization, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Organizations.List(ctx, username, nil)
}

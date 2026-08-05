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

type Client struct {
	Client *github.Client
	Token  string

	// rotate overrides token rotation in tests; nil means use rotateToken,
	// which calls NextClient (and therefore hits the database).
	rotate func() bool
}

var listGithubTokens = service.ListTokenByType

func InitClient(token string) *Client {
	githubClient := InitGithubClient(token)
	return &Client{
		Client: githubClient,
		Token:  token,
	}
}

func InitGithubClient(token string) *github.Client {
	httpClient := &http.Client{Transport: newGithubHTTPTransport()}
	gitClient := github.NewClient(httpClient).WithAuthToken(token)
	return gitClient
}

func newGithubHTTPTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return transport
}

func GetGithubClient() (*Client, error) {
	err, tokens := listGithubTokens("github")
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, errors.New("github client initialization failed: no token configured")
	}
	client := InitClient(tokens[0].Content)
	if client == nil {
		err = errors.New("github Client initial failed, please add token")
	}
	return client, err
}

func (c *Client) NextClient() (*github.Client, string) {
	currentToken := c.Token
	err, tokens := listGithubTokens("github")
	if err != nil {
		global.GVA_LOG.Error("github Client initial failed, please add token", zap.Error(err))
		return nil, ""
	}
	if len(tokens) == 0 {
		return nil, ""
	}
	var currentIndex int
	for index, token := range tokens {
		if token.Content == currentToken {
			currentIndex = index
		}
	}
	nextIndex := (currentIndex + 1) % len(tokens)
	nextToken := tokens[nextIndex]
	return InitGithubClient(nextToken.Content), nextToken.Content
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

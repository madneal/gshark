package githubsearch

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/models"
	"golang.org/x/oauth2"
	"os"
)

var (
	GithubClients map[string]*Client
	GithubClient  *Client
)

type Client struct {
	Client *github.Client
	Token  string
}

func init() {
	GithubClients = make(map[string]*Client)
	GithubClients, _ = InitGithubClients()
}

func InitGithubClients() (map[string]*Client, error) {
	githubClients := make(map[string]*Client)
	tokens, err := models.ListValidTokens("github")
	if err == nil {
		for _, token := range tokens {
			githubToken := token.Token
			gitClient := &github.Client{}
			if githubToken != "" {
				ctx := context.Background()
				ts := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: githubToken},
				)
				tc := oauth2.NewClient(ctx, ts)
				gitClient = github.NewClient(tc)
				githubClients[token.Token] = NewGitClient(gitClient, githubToken)
			}
		}
	}
	return githubClients, err
}

func GetGithubClient() (*Client, string, error) {
	var c *Client
	clients, err := InitGithubClients()
	for _, client := range clients {
		c = client
		break
	}
	if err != nil {
		logger.Log.Error(err)
	}
	if c == nil {
		fmt.Println("Github Client initial failed, please add token")
		os.Exit(3)
	}
	return c, c.Token, err
}

func NewGitClient(GithubClient *github.Client, token string) *Client {
	return &Client{Client: GithubClient, Token: token}
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

func (c *Client) GetUsersRepos(users []*github.User) map[string][]*github.Repository {
	result := make(map[string][]*github.Repository)
	for _, u := range users {
		repos, resp, _ := c.GetUserRepos(*u.Login)
		models.UpdateRate(c.Token, resp)
		result[*u.Login] = repos
	}
	return result
}

func (c *Client) GetStrUsersRepos(users []string) map[string][]*github.Repository {
	result := make(map[string][]*github.Repository)
	for _, u := range users {
		repos, resp, _ := c.GetUserRepos(u)
		models.UpdateRate(c.Token, resp)
		result[u] = repos
	}
	return result
}

func (c *Client) GetUserOrgs(username string) ([]*github.Organization, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Organizations.List(ctx, username, nil)
}

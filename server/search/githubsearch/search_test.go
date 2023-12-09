package githubsearch

import (
	"context"
	"fmt"
	"github.com/google/go-github/v57/github"
	"github.com/madneal/gshark/model"
	"os"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	token := os.Getenv("TOKEN")
	tokens := make([]model.Token, 0)
	tokens = append(tokens, model.Token{
		Content: token,
	})
	githubClients := InitGithubClients(tokens)
	ctx := context.Background()
	for _, client := range githubClients {
		i := 10
		for i > 0 {
			_, resp, err := client.Client.Search.Code(ctx, "repo:madneal/gshark security", &github.SearchOptions{
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 100,
				},
			})
			fmt.Println(resp.StatusCode)
			time.Sleep(100 * time.Second)
			if err != nil {
				fmt.Println(err)
			}
			i--
		}
	}
	fmt.Println(token)
}

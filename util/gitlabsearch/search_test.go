package gitlabsearch

import (
	"fmt"
	"github.com/neal1991/gshark/models"
	"github.com/xanzy/go-gitlab"
	"log"
	"testing"
)

func TestGetProjects(t *testing.T) {
	client := GetClient()
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		// Get the first page with projects.
		ps, resp, err := client.Projects.ListProjects(opt)
		if err != nil {
			log.Fatal(err)
		}

		// List all the projects we've found so far.
		for _, p := range ps {
			fmt.Printf("Found project: %s", p.Name)
		}

		// Exit the loop when we've seen all pages.
		//if resp.CurrentPage >= resp.TotalPages {
		//	break
		//}
		fmt.Println(resp.NextPage)
		if resp.NextPage == 0 {
			break
		}
		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}
}

func TestSearchCode(t *testing.T) {
	client := gitlab.NewClient(nil, "J2atZLikgHs1CWhEE_sV")
	inputInfo := models.InputInfo{
		ProjectId: 14625899,
	}
	codeResults := SearchCode("pkg", inputInfo, client)
	for _, result := range codeResults {
		fmt.Println(*result.Name)
		for index, text := range result.TextMatches {
			fmt.Println("==================")
			fmt.Println(index)
			fmt.Println(*text.Fragment)
		}
	}
	fmt.Println(len(codeResults))
}

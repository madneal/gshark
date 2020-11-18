package gitlabsearch

import (
	"fmt"
	"github.com/go-resty/resty/v2"
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
	client := resty.New()
	data := SearchCode("installation", client).Body()
	results := Parse(data)
	fmt.Println(results[0])
}

func TestBuildQueryString(t *testing.T) {
	fmt.Println(BuildQueryString("baidu", "ext"))
}

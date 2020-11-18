package gitlabsearch

import (
	"encoding/json"
	"github.com/madneal/gshark/logger"
)

type GitlabResponse struct {
	Basename   string
	Data       string
	Path       string
	Filename   string
	Ref        string
	Startline  int
	Project_id int64
}

type Project struct {
	Id                  int64
	Name                string
	Path_with_namespace string
}

func Parse(data []byte) []GitlabResponse {
	var gitlabResults []GitlabResponse
	err := json.Unmarshal(data, &gitlabResults)
	if err != nil {
		logger.Log.Error(err)
	}
	return gitlabResults
}

func ParseProject(data []byte) Project {
	var project Project
	err := json.Unmarshal(data, &project)
	if err != nil {
		logger.Log.Error(err)
	}
	return project
}

package gitlabsearch

import (
	"encoding/json"
	"github.com/madneal/gshark/logger"
)

type GitlabResponse struct {
	basename   string
	data       string
	path       string
	filename   string
	ref        string
	startline  int
	project_id int64
}

func Parse(data []byte) []GitlabResponse {
	var gitlabResults []GitlabResponse
	err := json.Unmarshal(data, &gitlabResults)
	if err != nil {
		logger.Log.Error(err)
	}
	return gitlabResults
}

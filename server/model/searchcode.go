package model

type SearchCodeRes struct {
	Previouspage string             `json:"previouspage"`
	Query        string             `json:"query"`
	Total        int                `json:"total"`
	Page         int                `json:"page"`
	Nextpage     int                `json:"nextpage"`
	Results      []SearchCodeResult `json:"results"`
}

type SearchCodeResult struct {
	Repo     string
	Name     string
	Language string
	Url      string
	Location string
	Lines    map[string]interface{} `json:"lines"`
	Filename string
	Id       int
}

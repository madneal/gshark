package request

type SearchResultSearch struct {
	//model.SearchResult
	PageInfo
	SearchInfo
	//Status int
	//Keyword string
	//Query   string
}

type SearchInfo struct {
	Status  int    `json:"status" form:"status"`
	Keyword string `json:"keyword" form:"keyword"`
	Query   string `json:"query" form:"query"`
}

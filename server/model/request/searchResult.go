package request

type SearchResultSearch struct {
	PageInfo
	SearchInfo
}

type SearchInfo struct {
	Status  int    `json:"status" form:"status"`
	Keyword string `json:"keyword" form:"keyword"`
	Query   string `json:"query" form:"query"`
}

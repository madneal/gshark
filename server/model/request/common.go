package request

// Paging common input parameter structure
type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// Find by id structure
type GetById struct {
	Id float64 `json:"id" form:"id"`
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

type BatchUpdateReq struct {
	Ids    []int `json:"ids" form:"ids"`
	Status int   `json:"status" form:"status"`
}

type UpdateReq struct {
	Status int    `json:"status" form:"status"`
	Repo   string `json:"repo" form:"repo"`
}

// Get role by id structure
type GetAuthorityId struct {
	AuthorityId string
}

type Empty struct{}

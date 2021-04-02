package request

import "github.com/madneal/gshark/model"

type SearchResultSearch struct {
	model.SearchResult
	PageInfo
}

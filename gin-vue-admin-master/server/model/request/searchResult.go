package request

import "gin-vue-admin/model"

type SearchResultSearch struct{
    model.SearchResult
    PageInfo
}
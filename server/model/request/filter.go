package request

import "github.com/madneal/gshark/model"

type FilterSearch struct{
    model.Filter
    PageInfo
}
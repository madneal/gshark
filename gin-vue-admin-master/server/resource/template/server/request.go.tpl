package request

import "github.com/madneal/gshark/model"

type {{.StructName}}Search struct{
    model.{{.StructName}}
    PageInfo
}
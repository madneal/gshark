package request

import "github.com/madneal/gshark/model"

type RepoSearch struct{
    model.Repo
    PageInfo
}
package request

import "gin-vue-admin/model"

type TokenSearch struct{
    model.Token
    PageInfo
}
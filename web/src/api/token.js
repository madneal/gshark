import service from '@/utils/request'

// @Tags Token
// @Summary 创建Token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Token true "创建Token"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /token/createToken [post]
export const createToken = (data) => {
     return service({
         url: "/token/createToken",
         method: 'post',
         data
     })
 }


// @Tags Token
// @Summary 删除Token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Token true "删除Token"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /token/deleteToken [delete]
 export const deleteToken = (data) => {
     return service({
         url: "/token/deleteToken",
         method: 'delete',
         data
     })
 }

// @Tags Token
// @Summary 删除Token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Token"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /token/deleteToken [delete]
 export const deleteTokenByIds = (data) => {
     return service({
         url: "/token/deleteTokenByIds",
         method: 'delete',
         data
     })
 }

// @Tags Token
// @Summary 更新Token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Token true "更新Token"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /token/updateToken [put]
 export const updateToken = (data) => {
     return service({
         url: "/token/updateToken",
         method: 'put',
         data
     })
 }


// @Tags Token
// @Summary 用id查询Token
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Token true "用id查询Token"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /token/findToken [get]
 export const findToken = (params) => {
     return service({
         url: "/token/findToken",
         method: 'get',
         params
     })
 }


// @Tags Token
// @Summary 分页获取Token列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取Token列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /token/getTokenList [get]
 export const getTokenList = (params) => {
     return service({
         url: "/token/getTokenList",
         method: 'get',
         params
     })
 }
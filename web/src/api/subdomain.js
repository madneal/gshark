import service from '@/utils/request'

// @Tags Subdomain
// @Summary 创建Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "创建Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /subdomain/createSubdomain [post]
export const createSubdomain = (data) => {
     return service({
         url: "/subdomain/createSubdomain",
         method: 'post',
         data
     })
 }


// @Tags Subdomain
// @Summary 删除Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "删除Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /subdomain/deleteSubdomain [delete]
 export const deleteSubdomain = (data) => {
     return service({
         url: "/subdomain/deleteSubdomain",
         method: 'delete',
         data
     })
 }

// @Tags Subdomain
// @Summary 删除Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /subdomain/deleteSubdomain [delete]
 export const deleteSubdomainByIds = (data) => {
     return service({
         url: "/subdomain/deleteSubdomainByIds",
         method: 'delete',
         data
     })
 }

// @Tags Subdomain
// @Summary 更新Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "更新Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /subdomain/updateSubdomain [put]
 export const updateSubdomain = (data) => {
     return service({
         url: "/subdomain/updateSubdomain",
         method: 'put',
         data
     })
 }


// @Tags Subdomain
// @Summary 用id查询Subdomain
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Subdomain true "用id查询Subdomain"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /subdomain/findSubdomain [get]
 export const findSubdomain = (params) => {
     return service({
         url: "/subdomain/findSubdomain",
         method: 'get',
         params
     })
 }


// @Tags Subdomain
// @Summary 分页获取Subdomain列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取Subdomain列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /subdomain/getSubdomainList [get]
 export const getSubdomainList = (params) => {
     return service({
         url: "/subdomain/getSubdomainList",
         method: 'get',
         params
     })
 }
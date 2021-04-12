import service from '@/utils/request'

// @Tags Filter
// @Summary 创建Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "创建Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /filter/createFilter [post]
export const createFilter = (data) => {
     return service({
         url: "/filter/createFilter",
         method: 'post',
         data
     })
 }


// @Tags Filter
// @Summary 删除Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "删除Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /filter/deleteFilter [delete]
 export const deleteFilter = (data) => {
     return service({
         url: "/filter/deleteFilter",
         method: 'delete',
         data
     })
 }

// @Tags Filter
// @Summary 删除Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /filter/deleteFilter [delete]
 export const deleteFilterByIds = (data) => {
     return service({
         url: "/filter/deleteFilterByIds",
         method: 'delete',
         data
     })
 }

// @Tags Filter
// @Summary 更新Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "更新Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /filter/updateFilter [put]
 export const updateFilter = (data) => {
     return service({
         url: "/filter/updateFilter",
         method: 'put',
         data
     })
 }


// @Tags Filter
// @Summary 用id查询Filter
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Filter true "用id查询Filter"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /filter/findFilter [get]
 export const findFilter = (params) => {
     return service({
         url: "/filter/findFilter",
         method: 'get',
         params
     })
 }


// @Tags Filter
// @Summary 分页获取Filter列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取Filter列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /filter/getFilterList [get]
 export const getFilterList = (params) => {
     return service({
         url: "/filter/getFilterList",
         method: 'get',
         params
     })
 }
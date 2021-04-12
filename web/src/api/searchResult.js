import service from '@/utils/request'

// @Tags SearchResult
// @Summary 创建SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "创建SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /searchResult/createSearchResult [post]
export const createSearchResult = (data) => {
     return service({
         url: "/searchResult/createSearchResult",
         method: 'post',
         data
     })
 }


// @Tags SearchResult
// @Summary 删除SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "删除SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /searchResult/deleteSearchResult [delete]
 export const deleteSearchResult = (data) => {
     return service({
         url: "/searchResult/deleteSearchResult",
         method: 'delete',
         data
     })
 }

// @Tags SearchResult
// @Summary 删除SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /searchResult/deleteSearchResult [delete]
 export const deleteSearchResultByIds = (data) => {
     return service({
         url: "/searchResult/deleteSearchResultByIds",
         method: 'delete',
         data
     })
 }

// @Tags SearchResult
// @Summary 更新SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "更新SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /searchResult/updateSearchResult [put]
 export const updateSearchResult = (data) => {
     return service({
         url: "/searchResult/updateSearchResult",
         method: 'post',
         data
     })
 }


// @Tags SearchResult
// @Summary 用id查询SearchResult
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SearchResult true "用id查询SearchResult"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /searchResult/findSearchResult [get]
 export const findSearchResult = (params) => {
     return service({
         url: "/searchResult/findSearchResult",
         method: 'get',
         params
     })
 }


// @Tags SearchResult
// @Summary 分页获取SearchResult列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取SearchResult列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /searchResult/getSearchResultList [get]
 export const getSearchResultList = (params) => {
     return service({
         url: "/searchResult/getSearchResultList",
         method: 'get',
         params
     })
 }


 export const updateSearchResultStatusByIds = (data) => {
     return service({
         url: '/searchResult/updateSearchResultStatusByIds',
         method: 'post',
         data
     })
 }
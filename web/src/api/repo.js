import service from '@/utils/request'

// @Tags Repo
// @Summary 创建Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "创建Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /repo/createRepo [post]
export const createRepo = (data) => {
     return service({
         url: "/repo/createRepo",
         method: 'post',
         data
     })
 }


// @Tags Repo
// @Summary 删除Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "删除Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /repo/deleteRepo [delete]
 export const deleteRepo = (data) => {
     return service({
         url: "/repo/deleteRepo",
         method: 'delete',
         data
     })
 }

// @Tags Repo
// @Summary 删除Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /repo/deleteRepo [delete]
 export const deleteRepoByIds = (data) => {
     return service({
         url: "/repo/deleteRepoByIds",
         method: 'delete',
         data
     })
 }

// @Tags Repo
// @Summary 更新Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "更新Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /repo/updateRepo [put]
 export const updateRepo = (data) => {
     return service({
         url: "/repo/updateRepo",
         method: 'put',
         data
     })
 }


// @Tags Repo
// @Summary 用id查询Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "用id查询Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /repo/findRepo [get]
 export const findRepo = (params) => {
     return service({
         url: "/repo/findRepo",
         method: 'get',
         params
     })
 }


// @Tags Repo
// @Summary 分页获取Repo列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取Repo列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /repo/getRepoList [get]
 export const getRepoList = (params) => {
     return service({
         url: "/repo/getRepoList",
         method: 'get',
         params
     })
 }
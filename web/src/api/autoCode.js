import service from '@/utils/request'

export const preview = (data) => {
    return service({
        url: "/autoCode/preview",
        method: 'post',
        data,
    })
}

export const createTemp = (data) => {
    return service({
        url: "/autoCode/createTemp",
        method: 'post',
        data,
        responseType: 'blob'
    })
}

// @Tags SysApi
// @Summary 获取当前所有数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /autoCode/getDatabase [get]
export const getDB = () => {
    return service({
        url: "/autoCode/getDB",
        method: 'get',
    })
}



// @Tags SysApi
// @Summary 获取当前数据库所有表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /autoCode/getTables [get]
export const getTable = (params) => {
    return service({
        url: "/autoCode/getTables",
        method: 'get',
        params,
    })
}

// @Tags SysApi
// @Summary 获取当前数据库所有表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /autoCode/getColumn [get]
export const getColumn = (params) => {
    return service({
        url: "/autoCode/getColumn",
        method: 'get',
        params,
    })
}
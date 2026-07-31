import service from '@/utils/request'

export const getScanLogList = (params) => {
    return service({
        url: '/scanLog/getScanLogList',
        method: 'get',
        params
    })
}

export const getScanLogOverview = () => {
    return service({
        url: '/scanLog/getScanLogOverview',
        method: 'get'
    })
}

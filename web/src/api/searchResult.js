import service from '@/utils/request'

export const createSearchResult = (data) => {
     return service({
         url: "/searchResult/createSearchResult",
         method: 'post',
         data
     })
 }

 export const deleteSearchResult = (data) => {
     return service({
         url: "/searchResult/deleteSearchResult",
         method: 'delete',
         data
     })
 }

 export const deleteSearchResultByIds = (data) => {
     return service({
         url: "/searchResult/deleteSearchResultByIds",
         method: 'delete',
         data
     })
 }

 export const updateSearchResult = (data) => {
     return service({
         url: "/searchResult/updateSearchResult",
         method: 'post',
         data
     })
 }

 export const findSearchResult = (params) => {
     return service({
         url: "/searchResult/findSearchResult",
         method: 'get',
         params
     })
 }

 export const getSearchResultList = (params) => {
     return service({
         url: "/searchResult/getSearchResultList",
         method: 'get',
         params
     })
 }

export const exportSearchResult = async (params) => {
    try {
        const response = await service({
            url: '/searchResult/exportSearchResult',
            method: 'get',
            params,
            responseType: 'blob',  // Set response type to blob
            headers: {
                'Accept': 'text/csv'
            }
        });
        const contentDisposition = response.headers?.['content-disposition'] || response.headers?.get('content-disposition');
        const filename = contentDisposition
            ? contentDisposition.split('filename=')[1]?.replace(/"/g, '')
            : 'export.csv';
        const blob = new Blob([response.data], {
            type: 'text/csv;charset=utf-8;'
        });
        const downloadUrl = window.URL.createObjectURL(blob);

        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = filename;
        document.body.appendChild(link);
        link.click();

        document.body.removeChild(link);
        window.URL.revokeObjectURL(downloadUrl);
    } catch (error) {
        console.error('Export failed:', error);
        throw error; // Propagate error to caller
    }
};

 export const updateSearchResultStatusByIds = (data) => {
     return service({
         url: '/searchResult/updateSearchResultStatusByIds',
         method: 'post',
         data
     })
 }

 export const startFilterTask = () => {
    return service({
        url: "/searchResult/startSecFilterTask",
        method: "post"
    })
 }

 export const getTaskStatus = () => {
    return service({
        url: "/searchResult/getTaskStatus",
        method: "get"
    })
 }
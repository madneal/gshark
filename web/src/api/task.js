import service from "@/utils/request";

export const getTaskList = (params) => {
    return service({
        url: "/task/getTaskList",
        method: "get",
        params
    })
}

export const createTask = (data) => {
    return service({
        url: "/task/createTask",
        method: "post",
        data
    })
}

export const switchTaskStatus = (data) => {
    return service({
        url: "/task/switchTaskStatus",
        method: "post",
        data
    })
}
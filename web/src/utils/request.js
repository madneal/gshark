import axios from 'axios'; // 引入axios
import { ElMessage } from 'element-plus';
import { store } from '@/store/index'
import router from '@/router/index'
import { bus } from '@/utils/bus'

const service = axios.create({
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 99999
})
let acitveAxios = 0
let timer
const showLoading = () => {
    acitveAxios++
    if (timer) {
        clearTimeout(timer)
    }
    timer = setTimeout(() => {
        if (acitveAxios > 0) {
            bus.emit("showLoading")
        }
    }, 400);
}

const closeLoading = () => {
        acitveAxios--
        if (acitveAxios <= 0) {
            clearTimeout(timer)
            bus.emit("closeLoading")
        }
    }
    //http request 拦截器
service.interceptors.request.use(
    config => {
        if (!config.donNotShowLoading) {
            showLoading()
        }
        const token = store.getters['user/token']
        const user = store.getters['user/userInfo']
        config.data = JSON.stringify(config.data);
        config.headers = {
            'Content-Type': 'application/json',
            'x-token': token,
            'x-user-id': user.ID
        }
        return config;
    },
    error => {
        closeLoading()
        ElMessage({
            showClose: true,
            message: error,
            type: 'error'
        })
        return error;
    }
);


//http response 拦截器
service.interceptors.response.use(
    response => {
        closeLoading()

        if (response.headers["new-token"]) {
            store.commit('user/setToken', response.headers["new-token"])
        }
        if(response.data.code == 0){
            if(response.data.data.needInit){
                ElMessage({
                    type:"info",
                    message:"您是第一次使用，请初始化"
                })
                    store.commit("user/NeedInit")
                    router.push({name:"init"})
            }
        }
        if (response.data.code == 0 || response.headers.success === "true" ) {
            return response.data
        } else {
            if (response.headers['content-type'] !== 'text/csv') {
                ElMessage({
                    showClose: true,
                    message: response.data.msg || decodeURI(response.headers.msg),
                    type: response.headers.msgtype||'error',
                })
            }

            if (response.data.data && response.data.data.reload) {
                store.commit('user/LoginOut')
            }
            return response.data.msg ? response.data : response
        }
    },
    error => {
        closeLoading()
        ElMessage({
            showClose: true,
            message: error,
            type: 'error'
        })
        return error
    }
)

export default service

import axios from "axios";
import { ElLoading } from "element-plus";

let loadingInstance;
let service = axios.create();

service.interceptors.request.use((config) => {
    loadingInstance = ElLoading.service({ fullscreen: true });
    return config;
});

service.interceptors.response.use((resp) => {
    loadingInstance.close();
    return resp;
}, (error) => {
    loadingInstance.close();
    return error;
});

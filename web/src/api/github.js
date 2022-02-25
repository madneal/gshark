import axios from "axios";
import { Loading } from "element-ui";

let loadingInstance;
let service = axios.create();

service.interceptors.request.use((config) => {
    loadingInstance = Loading.service({ fullscreen: true });
    return config;
});

service.interceptors.response.use((resp) => {
    loadingInstance.close();
    return resp;
}, (error) => {
    loadingInstance.close();
    return error;
});

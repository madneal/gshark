import { configureCompat, createApp } from 'vue'
import App from './App.vue'

import ElementPlus, { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 引入封装的router
import router from '@/router/index'

import '@/permission'
import { store } from '@/store/index'

// 路由守卫
import Bus from '@/utils/bus.js'

import { auth } from '@/directive/auth'

const app = createApp(App)

configureCompat({
    MODE: 3,
    COMPONENT_V_MODEL: false
})

app.config.compatConfig = {
    MODE: 3,
    COMPONENT_V_MODEL: false
}

app.use(ElementPlus)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.config.globalProperties.$loading = ElLoading.service
app.config.globalProperties.$message = ElMessage
app.config.globalProperties.$confirm = ElMessageBox.confirm

app.use(router)
app.use(store)
app.use(Bus)

// 按钮权限指令
auth(app)

app.mount('#app')

export default app

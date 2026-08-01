import { createStore } from 'vuex'
import VuexPersistence from 'vuex-persist'

import { user } from "@/store/module/user"
import { router } from "@/store/module/router"

const vuexLocal = new VuexPersistence({
    storage: window.localStorage,
    modules: ['user']
})
export const store = createStore({
    modules: {
        user,
        router
    },
    plugins: [vuexLocal.plugin]
})

import { createRouter, createWebHashHistory } from 'vue-router'

const baseRouters = [{
    path: '/',
    redirect: '/login'
},
{
    path: "/init",
    name: 'init',
    component: () =>
        import('@/view/init/init.vue')
},
{
    path: '/login',
    name: 'login',
    component: () =>
        import('@/view/login/login.vue')
}
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: baseRouters
})

router.addRoutes = (routes) => {
    routes.forEach(route => router.addRoute(route))
}

export default router

// build.config.js (终极混合模式版)

'use strict'

module.exports = {
    title: 'GShark',
    cdns: [
        /**
         * 终极规则：
         * 1. 如果有 customUrl，则直接使用，拥有最高优先级（用于特殊情况）。
         * 2. 如果有 urlTemplate，则从 package.json 读取版本号并替换 {version} 占位符。
         */
        {
            name: 'vue',
            scope: 'Vue',
            urlTemplate: 'https://cdnjs.cloudflare.com/ajax/libs/vue/{version}/vue.min.js'
        },
        {
            name: 'vue-router',
            scope: 'VueRouter',
            urlTemplate: 'https://cdnjs.cloudflare.com/ajax/libs/vue-router/{version}/vue-router.min.js'
        },
        {
            name: 'vuex',
            scope: 'Vuex',
            urlTemplate: 'https://cdnjs.cloudflare.com/ajax/libs/vuex/{version}/vuex.min.js'
        },
        {
            name: 'axios',
            scope: 'axios',
            urlTemplate: 'https://cdnjs.cloudflare.com/ajax/libs/axios/{version}/axios.min.js'
        },
        {
            name: 'element-ui',
            scope: 'ELEMENT',
            urlTemplate: 'https://cdnjs.cloudflare.com/ajax/libs/element-ui/{version}/index.js'
        },
    ]
};
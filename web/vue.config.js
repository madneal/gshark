'use strict'

const path = require('path')
const buildConf = require('./build.config')
const packageConf = require('./package.json')

function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    publicPath: './',
    outputDir: 'dist',
    assetsDir: 'static',
    lintOnSave: process.env.NODE_ENV === 'development',
    productionSourceMap: false,
    devServer: {
        port: 8080,
        open: true,
        proxy: {
            [process.env.VUE_APP_BASE_API]: {
                target: `http://127.0.0.1:8888/`,
                changeOrigin: true,
                pathRewrite: {
                    ['^' + process.env.VUE_APP_BASE_API]: ''
                }
            }
        },
    },
    configureWebpack: {
        resolve: {
            alias: {
                '@': resolve('src')
            }
        }
    },
    chainWebpack(config) {
        config.module
            .rule('vue')
            .use('vue-loader')
            .loader('vue-loader')
            .tap(options => {
                options.compilerOptions.preserveWhitespace = true
                return options
            })
            .end()

        config
            .when(process.env.NODE_ENV === 'development',
                config => config.devtool('eval-source-map')
            )

        config
            .when(process.env.NODE_ENV !== 'development',
                config => {
                    config.set('externals', buildConf.cdns.reduce((p, a) => {
                        p[a.name] = a.scope
                        return p
                    },{}))

                    config.plugin('html')
                        .tap(args => {
                            if(buildConf.title) {
                                args[0].title = buildConf.title
                            }
                            if(buildConf.cdns.length > 0) {
                                args[0].cdns = buildConf.cdns.map(conf => {
                                    const version = packageConf.dependencies[conf.name].replace(/[\^~]/g, '');

                                    if (conf.customUrl) {
                                        conf.js = conf.customUrl;
                                    } else if (conf.urlTemplate) {
                                        conf.js = conf.urlTemplate.replace('{version}', version);
                                    } else if (conf.path) {
                                        conf.js = `${buildConf.baseCdnUrl}${conf.path}`;
                                    } else {
                                        conf.js = `${buildConf.baseCdnUrl}/${conf.name}/${version}/${conf.name}.min.js`;
                                    }
                                    return conf;
                                })
                            }
                            return args
                        })

                    config
                        .optimization.splitChunks({
                        chunks: 'all',
                        cacheGroups: {
                            libs: {
                                name: 'chunk-libs',
                                test: /[\\/]node_modules[\\/]/,
                                priority: 10,
                                chunks: 'initial'
                            },
                            elementUI: {
                                name: 'chunk-elementUI',
                                priority: 20,
                                test: /[\\/]node_modules[\\/]_?element-ui(.*)/
                            },
                            commons: {
                                name: 'chunk-commons',
                                test: resolve('src/components'),
                                minChunks: 3,
                                priority: 5,
                                reuseExistingChunk: true
                            }
                        }
                    })
                    config.optimization.runtimeChunk('single')
                }
            )
    }
}
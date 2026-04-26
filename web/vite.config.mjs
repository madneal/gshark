import { fileURLToPath, URL } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import buildConf from './build.config.js'
import packageConf from './package.json' assert { type: 'json' }

const cdnConfigs = buildConf.cdns.map((conf) => {
  const version = packageConf.dependencies[conf.name]?.replace(/[\^~]/g, '')

  if (conf.customUrl) {
    return { ...conf, js: conf.customUrl }
  }
  if (conf.urlTemplate && version) {
    return { ...conf, js: conf.urlTemplate.replace('{version}', version) }
  }
  if (conf.path) {
    return { ...conf, js: `${buildConf.baseCdnUrl}${conf.path}` }
  }
  if (version) {
    return { ...conf, js: `${buildConf.baseCdnUrl}/${conf.name}/${version}/${conf.name}.min.js` }
  }
  return conf
})

const htmlTemplatePlugin = () => ({
  name: 'gshark-html-template',
  transformIndexHtml(html) {
    const cdnScripts = process.env.NODE_ENV === 'development'
      ? ''
      : cdnConfigs
        .filter((item) => item.js)
        .map((item) => `    <script type="text/javascript" src="${item.js}"></script>`)
        .join('\n')

    return html
      .replace('%APP_TITLE%', buildConf.title || 'GShark')
      .replace('<!-- CDN_SCRIPTS -->', cdnScripts)
  }
})

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const baseApi = env.VITE_BASE_API || '/api'

  return {
    base: './',
    plugins: [
      vue({
        template: {
          compilerOptions: {
            whitespace: 'preserve',
            compatConfig: {
              MODE: 2
            }
          }
        }
      }),
      htmlTemplatePlugin()
    ],
    resolve: {
      alias: {
        vue: '@vue/compat',
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      port: 8080,
      open: false,
      proxy: {
        [baseApi]: {
          target: 'http://127.0.0.1:8888/',
          changeOrigin: true,
          rewrite: (path) => path.replace(new RegExp(`^${baseApi}`), '')
        }
      }
    },
    build: {
      outDir: 'dist',
      assetsDir: 'static',
      sourcemap: false,
      rollupOptions: {
        external: mode === 'development'
          ? []
          : cdnConfigs.map((conf) => conf.name),
        output: {
          globals: cdnConfigs.reduce((globals, conf) => {
            globals[conf.name] = conf.scope
            return globals
          }, {}),
          manualChunks: {
            'chunk-element-plus': ['element-plus']
          }
        }
      }
    }
  }
})

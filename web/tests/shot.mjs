// 截图工具：node tests/shot.mjs <route> <outfile> [light|dark]
// 例：node tests/shot.mjs /login /tmp/login.png dark
import { chromium } from 'playwright-core'
import { baseUrl, chromePath, installAuth } from './auth.mjs'

const route = process.argv[2] || '/'
const out = process.argv[3] || `/tmp/gshark-shot.png`
const theme = process.argv[4] || 'dark'

let browser
try {
browser = await chromium.launch({
    executablePath: chromePath,
    headless: true,
    args: ['--no-sandbox']
})
const context = await browser.newContext({ viewport: { width: 1440, height: 900 } })
await context.addInitScript(
    ([t]) => localStorage.setItem('gshark-theme', t),
    [theme]
)
// 登录页之外注入伪造会话，模拟已登录
if (route !== '/login') await installAuth(context)
const page = await context.newPage()
const errors = []
page.on('console', m => { if (m.type() === 'error') errors.push(m.text()) })
page.on('pageerror', e => errors.push(String(e)))
// 路由是 hash 模式（createWebHashHistory），统一转成 /#/route 形式导航
const target = route.startsWith('http') ? route : new URL(`#${route}`, baseUrl).href
await page.goto(target, { waitUntil: 'load' })
await page.waitForTimeout(1500)
await page.screenshot({ path: out, fullPage: false })
console.log(`shot: ${out}${errors.length ? '\nconsole errors:\n' + errors.join('\n') : ''}`)
} finally {
  if (browser) await browser.close()
}

import { chromium } from 'playwright-core'
import { createHmac } from 'node:crypto'

const baseUrl = process.env.BUTTON_AUDIT_URL || 'http://localhost:8080/'
const chromePath = process.env.CHROME_PATH || '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'

const userInfo = {
  ID: 1,
  uuid: 'fb7d4a9e-e362-4eea-824d-6fb90c9fc128',
  userName: 'gshark',
  nickName: '超级管理员',
  headerImg: '',
  authorityId: '888',
  authority: {
    authorityId: '888',
    authorityName: '超级管理员',
    defaultRouter: 'state'
  }
}

const base64url = value => Buffer.from(JSON.stringify(value)).toString('base64url')
const token = (() => {
  const now = Math.floor(Date.now() / 1000)
  const payload = {
    UUID: userInfo.uuid,
    ID: userInfo.ID,
    Username: userInfo.userName,
    NickName: userInfo.nickName,
    AuthorityId: userInfo.authorityId,
    BufferTime: 86400,
    nbf: now - 60,
    exp: now + 86400,
    iss: 'qmPlus'
  }
  const unsigned = `${base64url({ alg: 'HS256', typ: 'JWT' })}.${base64url(payload)}`
  const signature = createHmac('sha256', 'qmPlus').update(unsigned).digest('base64url')
  return `${unsigned}.${signature}`
})()

const menus = [
  {
    path: 'admin',
    name: 'superAdmin',
    component: 'view/superAdmin/index.vue',
    meta: { title: '超级管理员', icon: 'user-solid' },
    children: [
      { path: 'authority', name: 'authority', component: 'view/superAdmin/authority/authority.vue', meta: { title: '角色管理', icon: 's-custom' } },
      { path: 'menu', name: 'menu', component: 'view/superAdmin/menu/menu.vue', meta: { title: '菜单管理', icon: 's-order' } },
      { path: 'api', name: 'api', component: 'view/superAdmin/api/api.vue', meta: { title: 'api管理', icon: 's-platform' } },
      { path: 'user', name: 'user', component: 'view/superAdmin/user/user.vue', meta: { title: '用户管理', icon: 'coordinate' } },
      { path: 'dictionary', name: 'dictionary', component: 'view/superAdmin/dictionary/sysDictionary.vue', meta: { title: '字典管理', icon: 'notebook-2' } },
      { path: 'dictionaryDetail/:id', name: 'dictionaryDetail', hidden: true, component: 'view/superAdmin/dictionary/sysDictionaryDetail.vue', meta: { title: '字典详情', icon: 's-order' } },
      { path: 'operation', name: 'operation', component: 'view/superAdmin/operation/sysOperationRecord.vue', meta: { title: '操作历史', icon: 'time' } }
    ]
  },
  {
    path: 'systemTools',
    name: 'systemTools',
    component: 'view/systemTools/index.vue',
    meta: { title: '系统工具', icon: 's-cooperation' },
    children: [
      { path: 'system', name: 'system', component: 'view/systemTools/system/system.vue', meta: { title: '系统配置', icon: 's-operation' } }
    ]
  },
  { path: 'person', name: 'person', hidden: true, component: 'view/person/person.vue', meta: { title: '个人信息', icon: 'message-solid' } },
  { path: 'state', name: 'state', component: 'view/system/state.vue', meta: { title: '服务器状态', icon: 'success' } },
  {
    path: 'setting',
    name: 'setting',
    component: 'view/routerHolder.vue',
    meta: { title: '管理', icon: 'setting' },
    children: [
      { path: 'rule', name: 'rule', component: 'view/rule/rule.vue', meta: { title: '规则管理', icon: 's-order' } },
      { path: 'token', name: 'token', component: 'view/token/token.vue', meta: { title: 'token管理', icon: 'collection' } },
      { path: 'filter', name: 'filter', component: 'view/filter/filter.vue', meta: { title: '过滤规则', icon: 's-opportunity' } }
    ]
  },
  {
    path: 'result',
    name: 'result',
    component: 'view/routerHolder.vue',
    meta: { title: '搜索结果', icon: 'search' },
    children: [
      { path: 'code', name: 'code', component: 'view/searchResult/searchResult.vue', meta: { title: '代码', icon: 'document' } },
      { path: 'subdomain', name: 'subdomain', component: 'view/subdomain/subdomain.vue', meta: { title: '子域名资产报告', icon: 's-data' } }
    ]
  }
]

const routes = [
  '/layout/state',
  '/layout/admin/authority',
  '/layout/admin/menu',
  '/layout/admin/api',
  '/layout/admin/user',
  '/layout/admin/dictionary',
  '/layout/admin/dictionaryDetail/1',
  '/layout/admin/operation',
  '/layout/systemTools/system',
  '/layout/person',
  '/layout/setting/rule',
  '/layout/setting/token',
  '/layout/setting/filter',
  '/layout/result/code',
  '/layout/result/subdomain'
]

const publicRoutes = ['/login', '/init']

const listPayload = list => ({ code: 0, data: { list, total: list.length, page: 1, pageSize: 10 }, msg: 'ok' })
const ok = (data = {}) => ({ code: 0, data, msg: 'ok' })

const sampleRows = {
  api: [{ ID: 1, path: '/demo/path', apiGroup: 'demo', description: 'demo api', method: 'GET' }],
  authority: [{ authorityId: '888', authorityName: '超级管理员', parentId: '0', dataAuthorityId: ['888'], children: [] }],
  menu: [{ ID: 1, path: 'demo', name: 'demo', component: 'view/routerHolder.vue', hidden: false, parentId: '0', meta: { title: 'demo', icon: 'setting' }, parameters: [] }],
  user: [{ ID: 1, uuid: 'uuid-1', userName: 'gshark', nickName: '超级管理员', headerImg: '', authority: { authorityName: '超级管理员' } }],
  dictionary: [{ ID: 1, name: 'demo', type: 'demo', status: true, desc: 'demo' }],
  dictionaryDetail: [{ ID: 1, label: '启用', value: 1, status: true, sort: 1 }],
  operation: [{ ID: 1, ip: '127.0.0.1', method: 'POST', path: '/demo', status: 200, latency: '1ms', agent: 'audit', body: '{}', resp: '{}', user: { userName: 'gshark', nickName: '超级管理员' } }],
  rule: [{ ID: 1, name: 'demo rule', ruleType: 'keyword,regex', content: 'secret', desc: 'demo', status: true }],
  token: [{ ID: 1, type: 'github', token: 'token-demo', status: 1 }],
  filter: [{ ID: 1, file: '.log', isFork: false, type: 'suffix', content: '.log' }],
  searchResult: [{ ID: 1, repoName: 'demo/repo', content: 'secret', keyword: 'secret', path: 'README.md', url: 'https://example.test', hash: 'abc', status: 0 }],
  subdomain: [{ ID: 1, subdomain: 'www.example.test', domain: 'example.test', status: 1 }]
}

const responseFor = url => {
  if (url.includes('/init/checkdb')) return ok({ needInit: false })
  if (url.includes('/base/captcha')) return ok({ picPath: '', captchaId: 'audit' })
  if (url.includes('/base/login')) return ok({ user: userInfo, token })
  if (url.includes('/menu/getMenu')) return ok({ menus })
  if (url.includes('/menu/getMenuList')) return listPayload(sampleRows.menu)
  if (url.includes('/menu/getBaseMenuTree')) return ok({ menus })
  if (url.includes('/menu/getMenuAuthority')) return ok({ menus: [] })
  if (url.includes('/authority/getAuthorityList')) return ok({ list: sampleRows.authority })
  if (url.includes('/api/getApiList')) return listPayload(sampleRows.api)
  if (url.includes('/api/getAllApis')) return ok({ apis: sampleRows.api })
  if (url.includes('/api/getApiById')) return ok({ api: sampleRows.api[0] })
  if (url.includes('/user/getUserList')) return listPayload(sampleRows.user)
  if (url.includes('/sysDictionary/getSysDictionaryList')) return listPayload(sampleRows.dictionary)
  if (url.includes('/sysDictionaryDetail/getSysDictionaryDetailList')) return listPayload(sampleRows.dictionaryDetail)
  if (url.includes('/sysOperationRecord/getSysOperationRecordList')) return listPayload(sampleRows.operation)
  if (url.includes('/system/getSystemConfig')) return ok({
    config: {
      system: { env: 'public', addr: 8888, dbType: 'mysql', useMultipoint: false },
      jwt: { signingKey: 'qmPlus' },
      zap: { level: 'info', format: 'console', prefix: '[GShark]', director: 'log', linkName: 'latest_log', encodeLevel: 'LowercaseColorLevelEncoder', stacktraceKey: 'stacktrace', showLine: false, logInConsole: true },
      email: { enable: false, to: 'audit@example.test', port: 465, from: 'audit@example.test', host: 'smtp.example.test', secret: '' },
      casbin: { modelPath: './resource/rbac_model.conf' },
      captcha: { keyLong: 6, imgWidth: 240, imgHeight: 80 },
      mysql: { username: 'root', password: '', path: 'localhost:3306', dbname: 'gshark', maxIdleConns: 0, maxOpenConns: 0, logMode: false },
      local: { path: 'uploads/file' },
      wechat: { enable: false, url: '' }
    }
  })
  if (url.includes('/system/getServerInfo')) return ok({
    server: {
      os: { goos: 'darwin', numCpu: 8, compiler: 'gc', goVersion: 'go1.23', numGoroutine: 12 },
      disk: { totalMb: 100000, usedMb: 50000, totalGb: 100, usedGb: 50, usedPercent: 50 },
      cpu: { cores: 8, cpus: [10, 20, 30, 40] },
      ram: { totalMb: 32768, usedMb: 16384, usedPercent: 50 }
    }
  })
  if (url.includes('/rule/getRuleList')) return listPayload(sampleRows.rule)
  if (url.includes('/token/getTokenList')) return listPayload(sampleRows.token)
  if (url.includes('/filter/getFilterList')) return listPayload(sampleRows.filter)
  if (url.includes('/searchResult/getSearchResultList')) return listPayload(sampleRows.searchResult)
  if (url.includes('/subdomain/getSubdomainList')) return listPayload(sampleRows.subdomain)
  if (url.includes('/rule/findRule')) return ok({ rule: sampleRows.rule[0] })
  if (url.includes('/token/findToken')) return ok({ retoken: sampleRows.token[0], token: sampleRows.token[0] })
  if (url.includes('/filter/findFilter')) return ok({ refilter: sampleRows.filter[0], filter: sampleRows.filter[0] })
  if (url.includes('/subdomain/findSubdomain')) return ok({ resubdomain: sampleRows.subdomain[0], subdomain: sampleRows.subdomain[0] })
  if (url.includes('/searchResult/findSearchResult')) return ok({ researchResult: sampleRows.searchResult[0], searchResult: sampleRows.searchResult[0] })
  if (url.includes('/sysDictionary/findSysDictionary')) return ok({ resysDictionary: sampleRows.dictionary[0], sysDictionary: sampleRows.dictionary[0] })
  if (url.includes('/sysDictionaryDetail/findSysDictionaryDetail')) return ok({ resysDictionaryDetail: sampleRows.dictionaryDetail[0], sysDictionaryDetail: sampleRows.dictionaryDetail[0] })
  if (url.includes('/find')) return ok({})
  if (url.includes('/get') || url.includes('/create') || url.includes('/update') || url.includes('/delete') || url.includes('/set') || url.includes('/copy') || url.includes('/switch') || url.includes('/start') || url.includes('/Test')) return ok({})
  return ok({})
}

const installMocks = async (context, withAuth = true) => {
  if (withAuth) {
    await context.addInitScript(({ token, userInfo }) => {
      localStorage.setItem('vuex', JSON.stringify({ user: { token, userInfo } }))
    }, { token, userInfo })
  }

  await context.route('**/*', async route => {
    const request = route.request()
    const url = request.url()
    if (!new URL(url).pathname.startsWith('/api/')) {
      await route.continue()
      return
    }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(responseFor(url))
    })
  })
}

const buttonInfo = async page => page.evaluate(() => {
  const seen = new Map()
  let visibleIndex = 0
  return Array.from(document.querySelectorAll('button')).map((button, index) => {
    const rect = button.getBoundingClientRect()
    const span = button.querySelector('span')
    const text = button.innerText.replace(/\s+/g, ' ').trim()
    const key = text || button.className || `button-${index}`
    const visible = rect.width > 0 && rect.height > 0
    const textOccurrence = visible ? (seen.get(key) || 0) : -1
    if (visible) seen.set(key, textOccurrence + 1)
    const currentVisibleIndex = visible ? visibleIndex++ : -1
    return {
      index,
      visibleIndex: currentVisibleIndex,
      text,
      key,
      textOccurrence,
      className: button.className,
      disabled: button.disabled || button.getAttribute('aria-disabled') === 'true',
      visible,
      width: Math.round(rect.width),
      height: Math.round(rect.height),
      textOverflow: span ? span.scrollWidth > span.clientWidth + 1 || span.scrollHeight > span.clientHeight + 1 : false
    }
  })
})

const unsafePattern = /登 出|规则导入/i

const clickVisibleButton = async (page, buttonInfo) => {
  const button = buttonInfo.text
    ? page.locator('button:visible').filter({ hasText: buttonInfo.text }).nth(buttonInfo.textOccurrence)
    : page.locator('button:visible').nth(buttonInfo.visibleIndex)
  await button.scrollIntoViewIfNeeded().catch(() => {})
  await button.click({ timeout: 3000 })
  await page.waitForTimeout(250)
}

const auditRoute = async (browser, routePath) => {
  const context = await browser.newContext({ viewport: { width: 1440, height: 900 } })
  await installMocks(context)
  const page = await context.newPage()
  const runtimeErrors = []
  const failedRequests = []

  page.on('pageerror', error => runtimeErrors.push(error.message))
  page.on('console', msg => {
    if (msg.type() === 'error' && !/compat|Vue Devtools/i.test(msg.text())) runtimeErrors.push(msg.text())
  })
  page.on('requestfailed', request => failedRequests.push(`${request.method()} ${request.url()} ${request.failure()?.errorText || ''}`))

  const loadRoute = async () => {
    await page.goto(`${baseUrl}#${routePath}`, { waitUntil: 'networkidle', timeout: 30000 })
    try {
      await page.locator('.admin-box').waitFor({ timeout: 10000 })
    } catch (error) {
      const bodyText = await page.locator('body').innerText({ timeout: 1000 }).catch(() => '')
      throw new Error(`layout did not mount for ${routePath}; url=${page.url()}; body="${bodyText.replace(/\s+/g, ' ').slice(0, 300)}"; errors="${runtimeErrors.join(' | ')}"; failed="${failedRequests.join(' | ')}"`)
    }
    await page.waitForTimeout(600)
  }

  await loadRoute()

  const initialButtons = await buttonInfo(page)
  const issues = []
  const clicked = []
  const skipped = []

  for (const button of initialButtons) {
    if (!button.visible) continue
    if (button.textOverflow) {
      issues.push(`text overflow in button #${button.index} "${button.text}" (${button.width}x${button.height})`)
    }
    if (button.disabled) continue
    if (unsafePattern.test(button.text)) {
      skipped.push(button.text || button.className)
      continue
    }
    const beforeErrors = runtimeErrors.length
    try {
      await loadRoute()
      await clickVisibleButton(page, button)
      clicked.push(button.text || button.className)
      const close = page.locator('.el-dialog button, .el-message-box button, .el-popover button').filter({ hasText: /取消|取 消|关闭/ }).last()
      if (await close.count()) {
        await close.click({ timeout: 2000 }).catch(() => {})
        await page.waitForTimeout(150)
      }
    } catch (error) {
      issues.push(`click failed for button #${button.index} "${button.text}": ${error.message}`)
    }
    if (runtimeErrors.length > beforeErrors) {
      issues.push(`runtime error after "${button.text}": ${runtimeErrors.slice(beforeErrors).join(' | ')}`)
    }
  }

  const finalButtons = await buttonInfo(page)
  for (const button of finalButtons) {
    if (button.visible && button.textOverflow) {
      issues.push(`post-click text overflow in button #${button.index} "${button.text}" (${button.width}x${button.height})`)
    }
  }

  if (runtimeErrors.length) issues.push(`runtime errors: ${runtimeErrors.join(' | ')}`)
  if (failedRequests.length) issues.push(`failed requests: ${failedRequests.join(' | ')}`)

  await context.close()
  return { routePath, buttons: initialButtons.map(button => button.text || button.className), clicked, skipped, issues }
}

const preparePublicRoute = async (page, routePath) => {
  if (routePath === '/login') {
    await page.locator('input[placeholder="请输入用户名"]').fill('gshark')
    await page.locator('input[placeholder="请输入密码"]').fill('password')
    await page.locator('input[placeholder="请输入验证码"]').fill('123456')
  }
}

const auditPublicRoute = async (browser, routePath) => {
  const context = await browser.newContext({ viewport: { width: 1440, height: 900 } })
  await installMocks(context, false)
  const page = await context.newPage()
  const runtimeErrors = []
  const failedRequests = []

  page.on('pageerror', error => runtimeErrors.push(error.message))
  page.on('console', msg => {
    if (msg.type() === 'error' && !/compat|Vue Devtools/i.test(msg.text())) runtimeErrors.push(msg.text())
  })
  page.on('requestfailed', request => failedRequests.push(`${request.method()} ${request.url()} ${request.failure()?.errorText || ''}`))

  const loadRoute = async () => {
    await page.goto(`${baseUrl}#${routePath}`, { waitUntil: 'networkidle', timeout: 30000 })
    await page.locator('body').waitFor({ timeout: 10000 })
    await page.waitForTimeout(500)
    await preparePublicRoute(page, routePath)
  }

  await loadRoute()
  const initialButtons = (await buttonInfo(page)).filter(button => button.visible)
  const issues = []
  const clicked = []
  const skipped = []

  for (const button of initialButtons) {
    if (button.textOverflow) {
      issues.push(`text overflow in button #${button.index} "${button.text}" (${button.width}x${button.height})`)
    }
    if (button.disabled) continue
    if (unsafePattern.test(button.text)) {
      skipped.push(button.text || button.className)
      continue
    }
    const beforeErrors = runtimeErrors.length
    try {
      await loadRoute()
      await clickVisibleButton(page, button)
      clicked.push(button.text || button.className)
      await page.waitForTimeout(500)
    } catch (error) {
      issues.push(`click failed for button #${button.index} "${button.text}": ${error.message}`)
    }
    if (runtimeErrors.length > beforeErrors) {
      issues.push(`runtime error after "${button.text}": ${runtimeErrors.slice(beforeErrors).join(' | ')}`)
    }
  }

  if (runtimeErrors.length) issues.push(`runtime errors: ${runtimeErrors.join(' | ')}`)
  if (failedRequests.length) issues.push(`failed requests: ${failedRequests.join(' | ')}`)

  await context.close()
  return { routePath, buttons: initialButtons.map(button => button.text || button.className), clicked, skipped, issues }
}

const browser = await chromium.launch({
  executablePath: chromePath,
  headless: true,
  args: ['--no-sandbox']
})

try {
  const results = []
  for (const routePath of publicRoutes) {
    results.push(await auditPublicRoute(browser, routePath))
  }
  for (const routePath of routes) {
    results.push(await auditRoute(browser, routePath))
  }

  const failures = results.filter(result => result.issues.length)
  for (const result of results) {
    console.log(`${result.routePath}: ${result.buttons.length} buttons, clicked ${result.clicked.length}, skipped ${result.skipped.length}`)
    if (result.issues.length) {
      for (const issue of result.issues) console.log(`  - ${issue}`)
    }
  }

  const totalButtons = results.reduce((sum, result) => sum + result.buttons.length, 0)
  const totalClicked = results.reduce((sum, result) => sum + result.clicked.length, 0)
  const totalSkipped = results.reduce((sum, result) => sum + result.skipped.length, 0)
  console.log(`button audit summary: ${publicRoutes.length + routes.length} routes, ${totalButtons} buttons, ${totalClicked} clicks, ${totalSkipped} skipped upload/logout buttons`)

  if (failures.length) {
    throw new Error(`${failures.length} routes have button audit issues`)
  }
} finally {
  await browser.close()
}

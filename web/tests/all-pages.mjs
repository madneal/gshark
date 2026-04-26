import { chromium } from 'playwright-core'
import { authHeaders, baseUrl, chromePath, installAuth } from './auth.mjs'

const fatalPattern = /TypeError|ReferenceError|Cannot read|Cannot set|is not a function|Unhandled/i
const ignoredConsolePatterns = [
  /compat/i,
  /Download the Vue Devtools extension/i
]

const browser = await chromium.launch({
  executablePath: chromePath,
  headless: true,
  args: ['--no-sandbox']
})

const isContainerMenu = menu => {
  return menu.component === 'view/routerHolder.vue' || menu.component === 'view/superAdmin/index.vue'
}

const collectRoutes = (menus, parents = [], titleParents = []) => {
  const routes = []
  for (const menu of menus || []) {
    const parts = [...parents, menu.path].filter(Boolean)
    const title = menu.meta?.title || menu.name || parts.join('/')
    const visibleChildren = (menu.children || []).filter(child => !child.hidden)
    const children = collectRoutes(menu.children, parts, [...titleParents, title])
    if (!menu.hidden && visibleChildren.length === 0 && !isContainerMenu(menu) && !parts.some(part => part.includes(':'))) {
      routes.push({
        title,
        path: `/layout/${parts.join('/')}`.replace(/\/+/g, '/'),
        component: menu.component,
        ancestors: titleParents,
        isLeaf: visibleChildren.length === 0
      })
    }
    routes.push(...children)
  }
  return routes
}

const normalizeError = message => String(message || '').replace(/\s+/g, ' ').trim()

const pageSnapshot = async page => {
  return page.evaluate(() => {
    const rect = selector => {
      const el = document.querySelector(selector)
      if (!el) return null
      const box = el.getBoundingClientRect()
      return { x: box.x, y: box.y, width: box.width, height: box.height }
    }
    return {
      layout: rect('.layout-cont'),
      aside: rect('.main-left'),
      main: rect('.main-right'),
      adminBox: rect('.admin-box'),
      menuItems: document.querySelectorAll('.el-menu-item').length,
      submenuTitles: document.querySelectorAll('.el-sub-menu__title').length,
      tables: document.querySelectorAll('.admin-box .el-table').length,
      empty: document.querySelectorAll('.admin-box .el-empty').length,
      cards: document.querySelectorAll('.admin-box .el-card').length,
      forms: document.querySelectorAll('.admin-box .el-form').length,
      adminText: document.querySelector('.admin-box')?.innerText.replace(/\s+/g, ' ').trim() || '',
      bodyText: document.body.innerText.replace(/\s+/g, ' ').trim()
    }
  })
}

const assertUsablePage = (route, snapshot, url) => {
  const issues = []
  if (!snapshot.layout || !snapshot.main || !snapshot.adminBox) {
    issues.push('missing layout/main/admin panel')
  }
  if (url.includes('#/login')) {
    issues.push('redirected to login')
  }
  if (route.path === '/' && !url.includes('#/layout/')) {
    issues.push(`did not resolve to a layout route: ${url}`)
  }
  if (snapshot.aside && snapshot.aside.width < 180) {
    issues.push(`sidebar collapsed unexpectedly: ${snapshot.aside.width}px`)
  }
  if (snapshot.main && snapshot.aside && snapshot.main.x < snapshot.aside.width - 1) {
    issues.push(`main overlaps sidebar: main.x=${snapshot.main.x}, aside.width=${snapshot.aside.width}`)
  }
  if (snapshot.adminBox && (snapshot.adminBox.width < 500 || snapshot.adminBox.height < 100)) {
    issues.push(`admin panel has bad size: ${snapshot.adminBox.width}x${snapshot.adminBox.height}`)
  }
  if (snapshot.menuItems < 8 || snapshot.submenuTitles < 3) {
    issues.push(`sidebar menu looks incomplete: items=${snapshot.menuItems}, submenus=${snapshot.submenuTitles}`)
  }
  if (snapshot.adminText.length < 8 && snapshot.tables + snapshot.empty + snapshot.cards + snapshot.forms === 0) {
    issues.push('main panel has no visible data, empty state, card, table, or form')
  }
  return issues
}

const openAncestor = async (page, title) => {
  const submenu = page.locator('.el-sub-menu').filter({
    has: page.locator('.el-sub-menu__title').filter({ hasText: title })
  }).first()
  if (await submenu.count() === 0) return
  const className = await submenu.getAttribute('class')
  if (!className?.includes('is-opened')) {
    await submenu.locator('.el-sub-menu__title').first().click()
    await page.waitForTimeout(150)
  }
}

const clickLeafMenu = async (page, route) => {
  for (const ancestor of route.ancestors) {
    await openAncestor(page, ancestor)
  }
  const item = page.locator('.el-menu-item').filter({ hasText: route.title }).first()
  await item.click({ timeout: 5000 })
  await page.waitForFunction(expectedPath => {
    return location.hash === `#${expectedPath}`
  }, route.path, { timeout: 5000 }).catch(() => {})
  await page.waitForLoadState('networkidle', { timeout: 5000 }).catch(() => {})
  await page.waitForFunction(() => {
    const adminBox = document.querySelector('.admin-box')
    if (!adminBox) return false
    const hasStructuredContent = adminBox.querySelector('.el-table, .el-empty, .el-card, .el-form')
    const text = adminBox.innerText.replace(/\s+/g, ' ').trim()
    return Boolean(hasStructuredContent || text.length >= 8)
  }, null, { timeout: 5000 }).catch(() => {})
}

try {
  const apiContext = await browser.newContext()
  const response = await apiContext.request.post(`${baseUrl.replace(/\/$/, '')}/api/menu/getMenu`, {
    headers: authHeaders,
    data: {}
  })
  const menuJson = await response.json()
  await apiContext.close()

  if (menuJson.code !== 0) {
    throw new Error(`getMenu failed: ${JSON.stringify(menuJson)}`)
  }

  const routes = collectRoutes(menuJson.data?.menus)
  if (routes.length === 0) {
    throw new Error('no frontend routes found in backend menu')
  }

  const failures = []
  const entries = [
    { title: 'authenticated default entry', path: '/', hash: '' },
    ...routes.map(route => ({ ...route, hash: `#${route.path}` }))
  ]

  for (const route of entries) {
    const context = await browser.newContext({
      viewport: { width: 1440, height: 900 }
    })
    await installAuth(context)
    const page = await context.newPage()
    const errors = []

    page.on('pageerror', error => {
      errors.push(normalizeError(error.message))
    })
    page.on('console', msg => {
      if (msg.type() === 'error' && !ignoredConsolePatterns.some(pattern => pattern.test(msg.text()))) {
        errors.push(normalizeError(msg.text()))
      }
    })

    try {
      await page.goto(`${baseUrl}${route.hash}`, { waitUntil: 'networkidle', timeout: 30000 })
      await page.waitForTimeout(800)

      const url = page.url()
      await page.locator('.admin-box').waitFor({ timeout: 5000 })
      const snapshot = await pageSnapshot(page)
      const fatalErrors = errors.filter(error => fatalPattern.test(error))
      const pageIssues = assertUsablePage(route, snapshot, url)

      if (fatalErrors.length > 0) {
        failures.push(`${route.path} (${route.title}) runtime errors: ${fatalErrors.join(' | ')}`)
      } else if (pageIssues.length > 0) {
        failures.push(`${route.path} (${route.title}) page/style issues: ${pageIssues.join(' | ')}`)
      } else {
        console.log(`ok ${route.path} - ${route.title}`)
      }
    } catch (error) {
      failures.push(`${route.path} (${route.title}) navigation failed: ${normalizeError(error.message)}`)
    } finally {
      await context.close()
    }
  }

  const menuContext = await browser.newContext({
    viewport: { width: 1440, height: 900 }
  })
  await installAuth(menuContext)
  const menuPage = await menuContext.newPage()
  const menuErrors = []
  menuPage.on('pageerror', error => {
    menuErrors.push(normalizeError(error.message))
  })
  menuPage.on('console', msg => {
    if (msg.type() === 'error' && !ignoredConsolePatterns.some(pattern => pattern.test(msg.text()))) {
      menuErrors.push(normalizeError(msg.text()))
    }
  })

  try {
    await menuPage.goto(baseUrl, { waitUntil: 'networkidle', timeout: 30000 })
    await menuPage.waitForTimeout(800)
    for (const route of routes.filter(route => route.isLeaf)) {
      await clickLeafMenu(menuPage, route)
      const url = menuPage.url()
      const snapshot = await pageSnapshot(menuPage)
      const fatalErrors = menuErrors.filter(error => fatalPattern.test(error))
      const pageIssues = assertUsablePage(route, snapshot, url)

      if (!url.includes(`#${route.path}`)) {
        pageIssues.push(`menu click did not navigate to ${route.path}: ${url}`)
      }
      if (fatalErrors.length > 0) {
        failures.push(`menu click ${route.path} (${route.title}) runtime errors: ${fatalErrors.join(' | ')}`)
      } else if (pageIssues.length > 0) {
        failures.push(`menu click ${route.path} (${route.title}) issues: ${pageIssues.join(' | ')}`)
      } else {
        console.log(`menu ok ${route.path} - ${route.title}`)
      }
      menuErrors.length = 0
    }
  } catch (error) {
    failures.push(`menu click regression failed: ${normalizeError(error.message)}`)
  } finally {
    await menuContext.close()
  }

  if (failures.length > 0) {
    throw new Error(`page crawl failed:\n${failures.join('\n')}`)
  }

  console.log(`all page regression passed (${routes.length} routes + default entry)`)
} finally {
  await browser.close()
}

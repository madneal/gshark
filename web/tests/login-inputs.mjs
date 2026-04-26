import { chromium } from 'playwright-core'

const chromePath = process.env.CHROME_PATH || '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'
const baseUrl = process.env.LOGIN_TEST_URL || 'http://localhost:8081/'

const browser = await chromium.launch({
  executablePath: chromePath,
  headless: true,
  args: ['--no-sandbox']
})

try {
  const page = await browser.newPage()
  const errors = []

  page.on('pageerror', error => errors.push(error.message))
  page.on('console', msg => {
    if (msg.type() === 'error') {
      errors.push(msg.text())
    }
  })

  await page.goto(baseUrl, { waitUntil: 'networkidle' })

  const cases = [
    ['请输入用户名', 'aliceuser'],
    ['请输入密码', 'secretpass'],
    ['请输入验证码', 'abcd']
  ]

  const focused = []
  for (const [placeholder, text] of cases) {
    const input = page.locator(`input[placeholder="${placeholder}"]`)
    await input.click()
    await page.keyboard.type(text, { delay: 20 })
    await page.waitForTimeout(100)
    focused.push(await input.inputValue())
  }

  await page.locator('body').click({ position: { x: 5, y: 5 } })

  const blurred = []
  for (const [placeholder] of cases) {
    blurred.push(await page.locator(`input[placeholder="${placeholder}"]`).inputValue())
  }

  const expected = cases.map(([, text]) => text)
  if (JSON.stringify(focused) !== JSON.stringify(expected)) {
    throw new Error(`login inputs did not keep typed values while focused: ${JSON.stringify(focused)}`)
  }

  if (JSON.stringify(blurred) !== JSON.stringify(expected)) {
    throw new Error(`login inputs did not keep typed values after blur: ${JSON.stringify(blurred)}`)
  }

  const fatalErrors = errors.filter(error => /TypeError|Cannot read|Unhandled|ReferenceError/.test(error))
  if (fatalErrors.length > 0) {
    throw new Error(`browser errors: ${fatalErrors.join(' | ')}`)
  }

  console.log('login input regression passed')
} finally {
  await browser.close()
}

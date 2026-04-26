import { createHmac } from 'node:crypto'

export const baseUrl = process.env.FE_TEST_URL || 'http://localhost:8081/'
export const chromePath = process.env.CHROME_PATH || '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'

export const userInfo = {
  ID: 1,
  uuid: 'fb7d4a9e-e362-4eea-824d-6fb90c9fc128',
  userName: 'gshark',
  nickName: '超级管理员',
  headerImg: '',
  authorityId: '888',
  authority: {
    authorityId: '888',
    authorityName: '普通用户',
    defaultRouter: 'state'
  }
}

const base64url = value => Buffer.from(JSON.stringify(value)).toString('base64url')

export const token = (() => {
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

export const installAuth = async context => {
  await context.addInitScript(({ token, userInfo }) => {
    localStorage.setItem('vuex', JSON.stringify({
      user: {
        token,
        userInfo
      }
    }))
  }, { token, userInfo })
}

export const authHeaders = {
  'content-type': 'application/json',
  'x-token': token,
  'x-user-id': String(userInfo.ID)
}

export function cliAuthorization(query) {
    if (typeof query.callback !== 'string' || typeof query.state !== 'string' || !/^[a-f0-9]{64}$/.test(query.state)) {
        throw new Error('CLI 授权链接无效，请重新运行 auth login')
    }
    const callback = new URL(query.callback)
    if (callback.protocol !== 'http:' || callback.hostname !== '127.0.0.1' || !callback.port ||
        callback.username || callback.password || callback.pathname !== '/callback' || callback.search || callback.hash) {
        throw new Error('CLI 授权回调无效，请重新运行 auth login')
    }
    return { callback: callback.href, state: query.state }
}

export function cliSessionActive(token) {
    try {
        const payload = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
        const { exp } = JSON.parse(atob(payload))
        return typeof exp === 'number' && exp > Date.now() / 1000
    } catch {
        return false
    }
}

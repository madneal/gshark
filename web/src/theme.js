// 主题初始化与切换：Element Plus 的暗色由 html.dark 驱动（dark/css-vars.css 已在 main.js 引入），
// 用户选择持久化在 localStorage。默认暗色 —— 管理控制台的既有习惯。
const KEY = 'gshark-theme'

export function currentTheme() {
    return localStorage.getItem(KEY) === 'light' ? 'light' : 'dark'
}

export function applyTheme(theme) {
    document.documentElement.classList.toggle('dark', theme === 'dark')
    localStorage.setItem(KEY, theme)
}

applyTheme(currentTheme())

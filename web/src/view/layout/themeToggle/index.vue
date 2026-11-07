<template>
    <el-tooltip :content="isDark ? '切换到亮色' : '切换到暗色'" placement="bottom">
        <button class="theme-toggle" type="button" @click="toggle">
            <el-icon><Sunny v-if="isDark" /><Moon v-else /></el-icon>
        </button>
    </el-tooltip>
</template>

<script>
import { currentTheme, applyTheme } from '@/theme'

export default {
    name: 'ThemeToggle',
    data() {
        return { isDark: currentTheme() === 'dark' }
    },
    methods: {
        toggle() {
            this.isDark = !this.isDark
            applyTheme(this.isDark ? 'dark' : 'light')
        }
    }
}
</script>

<style scoped>
.theme-toggle {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    border: none;
    border-radius: var(--gs-radius);
    background: transparent;
    color: var(--gs-text-secondary);
    font-size: 17px;
    cursor: pointer;
    transition: background-color 0.15s, color 0.15s;
}
.theme-toggle:hover {
    background: var(--gs-bg-subtle);
    color: var(--gs-text);
}
</style>

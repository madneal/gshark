<template>
  <main class="cli-auth-page">
    <section class="cli-auth-card" aria-labelledby="cli-auth-title">
      <img src="@/assets/nav_logo.png" alt="GShark" class="cli-auth-logo" />
      <p class="cli-auth-eyebrow">GShark CLI</p>
      <h1 id="cli-auth-title">授权命令行访问</h1>
      <p v-if="error" role="alert">{{ error }}</p>
      <template v-else-if="!canceled">
        <p>以 <strong>{{ username }}</strong> 的身份授权本机 CLI，使用当前账号的权限管理 GShark。</p>
        <p class="cli-auth-hint">仅在你刚刚运行了 <code>gshark-cli auth login</code> 时继续。</p>
        <form :action="authorization.callback" method="post" @submit="submitting = true">
          <input type="hidden" name="state" :value="authorization.state" />
          <input type="hidden" name="token" :value="token" />
          <input type="hidden" name="username" :value="username" />
          <div class="cli-auth-actions">
            <el-button @click="canceled = true" :disabled="submitting">取消</el-button>
            <el-button type="primary" native-type="submit" :loading="submitting">授权 CLI</el-button>
          </div>
        </form>
      </template>
      <p v-else role="status">已取消授权。可以关闭此页面，并在终端按 Ctrl+C 结束等待。</p>
    </section>
  </main>
</template>

<script>
import { cliAuthorization } from '@/utils/cliAuth'

export default {
    name: 'CliAuth',
    data() {
        return { authorization: null, error: '', canceled: false, submitting: false }
    },
    computed: {
        token() { return this.$store.getters['user/token'] },
        username() { return this.$store.getters['user/userInfo'].userName }
    },
    watch: {
        '$route.query': {
            immediate: true,
            handler(query) {
                this.error = ''
                this.authorization = null
                this.canceled = false
                this.submitting = false
                try {
                    this.authorization = cliAuthorization(query)
                } catch (error) {
                    this.error = error.message
                }
            }
        }
    }
}
</script>

<style scoped>
.cli-auth-page {
    min-height: 100%;
    display: grid;
    place-items: center;
    padding: 24px;
    box-sizing: border-box;
    background: var(--gs-bg-page, var(--gs-dark-bg));
    color: var(--gs-text, var(--gs-dark-text));
}
.cli-auth-card {
    width: 100%;
    max-width: 440px;
    box-sizing: border-box;
    padding: 32px;
    background: var(--gs-bg-card, var(--gs-dark-panel));
    border: 1px solid var(--gs-border, var(--gs-dark-border));
    border-radius: var(--gs-radius-lg, 12px);
}
.cli-auth-logo { width: 44px; height: 44px; object-fit: contain; }
.cli-auth-eyebrow { color: var(--gs-text-secondary, var(--gs-dark-muted)); font-size: 13px; letter-spacing: .08em; }
h1 { margin: 16px 0; font-size: 24px; }
p { line-height: 1.7; overflow-wrap: anywhere; }
.cli-auth-hint { color: var(--gs-text-secondary, var(--gs-dark-muted)); font-size: 13px; }
code { font-family: var(--gs-font-mono, monospace); }
.cli-auth-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 24px; }
</style>

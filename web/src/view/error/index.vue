<template>
  <div class="big">
    <div class="inner">
      <img src="../../assets/notFound.png" alt="页面不存在" />
      <p>页面被神秘力量吸走了</p>
      <p class="hint">常见问题为菜单未分配默认页，可自行分配菜单或修改默认路由</p>
      <div class="actions">
        <el-button type="primary" @click="goHome">返回首页</el-button>
        <el-button @click="goLogin">去登录</el-button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "NotFound",
  methods: {
    goHome() {
      this.$router.push({ path: "/" }).catch(() => {
        this.$router.push({ name: "login" });
      });
    },
    goLogin() {
      this.$router.push({ name: "login" });
    },
  },
};
</script>

<style lang="scss" scoped>
.big {
  min-height: calc(100vh - 220px);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
}
.inner {
  text-align: center;
  color: var(--gs-dark-text, #e5e7eb);
}
.inner img {
  max-width: 220px;
  margin-bottom: 12px;
}
.inner p {
  font-size: 20px;
  margin: 0 0 8px;
}
.hint {
  font-size: 14px !important;
  color: var(--gs-dark-muted, #94a3b8);
}
.actions {
  margin-top: 16px;
}
</style>

<template>
  <div class="error-page">
    <div class="error-card">
      <img src="../../assets/notFound.png" alt="页面不存在" />
      <h1>页面不存在</h1>
      <p>页面被神秘力量吸走了。常见原因是菜单未分配默认页，可在权限中分配菜单或修改默认路由。</p>
      <div class="error-actions">
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
.error-page {
  min-height: calc(100vh - 220px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  background: transparent;
}

.error-card {
  max-width: 520px;
  text-align: center;
  padding: 32px 28px;
  background: rgba(24, 34, 53, 0.8);
  border: 1px solid rgba(148, 163, 184, 0.12);
  border-radius: 12px;
  box-shadow: 0 18px 48px rgba(0, 0, 0, 0.28);

  img {
    max-width: 220px;
    width: 100%;
    height: auto;
    margin-bottom: 16px;
    opacity: 0.92;
  }

  h1 {
    margin: 0 0 10px;
    font-size: 22px;
    font-weight: 700;
    color: #f1f5f9;
  }

  p {
    margin: 0 0 22px;
    font-size: 14px;
    line-height: 1.6;
    color: #94a3b8;
  }
}

.error-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}
</style>

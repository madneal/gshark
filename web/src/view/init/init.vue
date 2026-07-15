<template>
  <div class="init-page">
    <div class="init-shell">
      <div class="init-brand">
        <img class="logo" src="@/assets/nav_logo.png" alt="GShark" />
        <h1>欢迎使用 GShark</h1>
        <p>首次使用需要初始化数据库并写入基础数据</p>
      </div>
      <div class="init-card">
        <el-form ref="form" :model="form" label-width="100px" label-position="top">
          <el-form-item label="数据库类型">
            <el-select disabled v-model="form.sqlType" placeholder="请选择" style="width: 100%">
              <el-option key="mysql" label="MySQL（目前仅支持）" value="mysql" />
            </el-select>
          </el-form-item>
          <el-form-item label="Host">
            <el-input v-model="form.host" placeholder="例如 127.0.0.1" />
          </el-form-item>
          <el-form-item label="Port">
            <el-input v-model="form.port" placeholder="例如 3306" />
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="form.userName" placeholder="数据库用户名" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input
              v-model="form.password"
              type="password"
              show-password
              placeholder="无密码可留空"
            />
          </el-form-item>
          <el-form-item label="数据库名">
            <el-input v-model="form.dbName" placeholder="例如 gshark" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" class="init-submit" @click="onSubmit">
              立即初始化
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { initDB } from "@/api/initdb";
export default {
  name: "Init",
  data() {
    return {
      form: {
        sqlType: "mysql",
        host: "127.0.0.1",
        port: "3306",
        userName: "root",
        password: "",
        dbName: "gshark",
      },
    };
  },
  methods: {
    async onSubmit() {
      const loading = this.$loading({
        lock: true,
        text: "正在初始化数据库，请稍候",
        background: "rgba(0, 0, 0, 0.7)",
      });
      try {
        const res = await initDB(this.form);
        if (res.code == 0) {
          this.$message({
            type: "success",
            message: res.msg,
          });
          this.$router.push({ name: "login" });
        }
      } finally {
        loading.close();
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.init-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 16px;
  background: linear-gradient(135deg, #0b1120 0%, #0f172a 50%, #111827 100%);
}

.init-shell {
  width: 100%;
  max-width: 480px;
}

.init-brand {
  text-align: center;
  margin-bottom: 28px;

  .logo {
    width: 52px;
    height: 52px;
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.08);
    padding: 6px;
  }

  h1 {
    margin: 16px 0 8px;
    font-size: 26px;
    font-weight: 700;
    color: #f1f5f9;
  }

  p {
    margin: 0;
    color: #94a3b8;
    font-size: 14px;
  }
}

.init-card {
  padding: 28px 24px;
  background: rgba(24, 34, 53, 0.8);
  border: 1px solid rgba(148, 163, 184, 0.12);
  border-radius: 12px;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(12px);
}

.init-submit {
  width: 100%;
  height: 40px;
  font-weight: 600;
}
</style>

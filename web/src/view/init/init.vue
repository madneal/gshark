<template>
  <div class="init">
    <p class="in-one">欢迎使用 GShark</p>
    <p class="in-two">您需要初始化数据库并填充初始数据</p>
    <div class="form-card">
      <el-form ref="form" :model="form" label-width="100px">
        <el-form-item label="数据库类型">
          <el-select disabled v-model="form.sqlType" placeholder="请选择">
            <el-option key="mysql" label="mysql（目前只支持 mysql）" value="mysql" />
          </el-select>
        </el-form-item>
        <el-form-item label="host">
          <el-input v-model="form.host" placeholder="请输入数据库地址" />
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model="form.port" placeholder="请输入数据库端口" />
        </el-form-item>
        <el-form-item label="userName">
          <el-input v-model="form.userName" placeholder="请输入数据库用户名" />
        </el-form-item>
        <el-form-item label="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            placeholder="请输入数据库密码（没有则为空）"
          />
        </el-form-item>
        <el-form-item label="dbName">
          <el-input v-model="form.dbName" placeholder="请输入数据库名称" />
        </el-form-item>
        <el-form-item>
          <div style="text-align: right">
            <el-button type="primary" @click="onSubmit">立即初始化</el-button>
          </div>
        </el-form-item>
      </el-form>
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
          this.$message({ type: "success", message: res.msg });
          this.$router.push({ name: "login" });
        }
      } finally {
        loading.close();
      }
    },
  },
};
</script>

<style lang="scss">
.init {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--gs-dark-bg, #0f172a);
  color: var(--gs-dark-text, #e5e7eb);
}
.in-one {
  font-size: 26px;
  margin-bottom: 8px;
}
.in-two {
  font-size: 16px;
  color: var(--gs-dark-muted, #94a3b8);
}
.form-card {
  margin-top: 32px;
  width: min(60vw, 560px);
  padding: 20px;
  border-radius: 8px;
  background: var(--gs-dark-panel, #182235);
  border: 1px solid var(--gs-dark-border-soft, rgba(148, 163, 184, 0.18));
}
</style>

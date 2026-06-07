<template>
  <div id="userLayout" class="user-layout-wrapper">
    <div class="container">
      <div class="top">
        <div class="desc">
          <img class="logo_login" src="@/assets/nav_logo.png" alt="" />
        </div>
        <div class="header">
          <a href="/">
            <span class="title">GShark</span>
          </a>
        </div>
      </div>
      <div class="main">
        <el-form
          :model="loginForm"
          ref="loginForm"
          @keyup.enter="submitForm"
        >
          <el-form-item prop="username">
            <el-input placeholder="请输入用户名" v-model="loginForm.username">
              <template #suffix>
                <el-icon><User /></el-icon>
              </template>
            ></el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              :type="lock === 'lock' ? 'password' : 'text'"
              placeholder="请输入密码"
              v-model="loginForm.password"
            >
              <template #suffix>
                <el-icon class="password-toggle" @click="changeLock">
                  <Unlock v-if="lock === 'lock'" />
                  <Lock v-else />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item style="position: relative">
            <el-input
              v-model="loginForm.captcha"
              name="logVerify"
              placeholder="请输入验证码"
              style="width: 60%"
            />
            <button
              class="vPic"
              type="button"
              role="button"
              title="刷新验证码"
              @click.prevent.stop="loginVefify"
            >
              <img
                v-if="picPath"
                :src="picPath"
                width="100%"
                height="100%"
                alt="请输入验证码"
              />
              <span v-else class="captcha-placeholder">
                {{ captchaLoading ? "加载中" : "刷新" }}
              </span>
            </button>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitForm" style="width: 100%"
              >登 录</el-button
            >
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
/* eslint-disable vue/multi-word-component-names */
import { mapActions } from "vuex";
import { captcha } from "@/api/user";
export default {
  name: "Login",
  data() {
    return {
      curYear: 0,
      lock: "lock",
      loginForm: {
        username: "",
        password: "",
        captcha: "",
        captchaId: "",
      },
      logVerify: "",
      picPath: "",
      captchaLoading: false,
    };
  },
  created() {
    this.loginVefify();
    this.curYear = new Date().getFullYear();
  },
  methods: {
    ...mapActions("user", ["LoginIn"]),
    async login() {
      return await this.LoginIn(this.loginForm);
    },
    async submitForm() {
      this.$refs.loginForm.validate(async (v) => {
        if (v) {
          const flag = await this.login();
          if (!flag) {
            this.loginVefify();
          }
        } else {
          this.$message({
            type: "error",
            message: "请正确填写登录信息",
            showClose: true,
          });
          this.loginVefify();
          return false;
        }
      });
    },
    changeLock() {
      this.lock === "lock" ? (this.lock = "unlock") : (this.lock = "lock");
    },
    async loginVefify() {
      if (this.captchaLoading) {
        return;
      }
      this.captchaLoading = true;
      try {
        const result = await captcha({});
        const captchaData = result?.data;
        if (!captchaData?.picPath || !captchaData?.captchaId) {
          if (!this.picPath) {
            this.loginForm.captchaId = "";
          }
          return;
        }
        this.picPath = captchaData.picPath;
        this.loginForm.captchaId = captchaData.captchaId;
        this.loginForm.captcha = "";
      } finally {
        this.captchaLoading = false;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/style/login.scss";

.password-toggle {
  cursor: pointer;
}
</style>

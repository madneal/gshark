<template>
  <div class="person-page">
    <el-row :gutter="16">
      <el-col :span="6" :xs="24">
        <div class="user-card">
          <div
            class="user-headpic-update"
            :style="{
              'background-image': `url(${
                userInfo.headerImg && userInfo.headerImg.slice(0, 4) !== 'http'
                  ? path + userInfo.headerImg
                  : userInfo.headerImg
              })`,
              'background-repeat': 'no-repeat',
              'background-size': 'cover',
            }"
          >
            <span class="update" @click="openChooseImg">
              <el-icon><Edit /></el-icon>
              重新上传
            </span>
          </div>
          <p class="nickname">{{ userInfo.nickName }}</p>
          <p class="username">
            <el-icon><User /></el-icon>
            {{ userInfo.userName || userInfo.nickName }}
          </p>
        </div>
      </el-col>
      <el-col :span="18" :xs="24">
        <div class="user-panel">
          <el-tabs v-model="activeName">
            <el-tab-pane label="账号绑定" name="second">
              <div class="account-row">
                <div>
                  <p class="title">修改密码</p>
                  <p class="desc">定期更新密码有助于账号安全</p>
                </div>
                <el-button type="primary" link @click="showPassword = true">
                  修改密码
                </el-button>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </el-col>
    </el-row>

    <ChooseImg ref="chooseImg" @enter-img="enterImg" />

    <el-dialog
      v-model="showPassword"
      @close="clearPassword"
      title="修改密码"
      width="400px"
    >
      <el-form
        :model="pwdModify"
        :rules="rules"
        label-width="80px"
        ref="modifyPwdForm"
      >
        <el-form-item label="原密码" prop="password">
          <el-input show-password v-model="pwdModify.password" />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input show-password v-model="pwdModify.newPassword" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input show-password v-model="pwdModify.confirmPassword" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showPassword = false">取 消</el-button>
          <el-button @click="savePassword" type="primary">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import ChooseImg from "@/components/chooseImg/index.vue";
import { setUserInfo, changePassword } from "@/api/user";
import { mapGetters, mapMutations } from "vuex";

const path = import.meta.env.VITE_BASE_API;

export default {
  name: "Person",
  data() {
    return {
      path: path,
      activeName: "second",
      showPassword: false,
      pwdModify: {},
      rules: {
        password: [
          { required: true, message: "请输入密码", trigger: "blur" },
          { min: 6, message: "最少6个字符", trigger: "blur" },
        ],
        newPassword: [
          { required: true, message: "请输入新密码", trigger: "blur" },
          { min: 6, message: "最少6个字符", trigger: "blur" },
        ],
        confirmPassword: [
          { required: true, message: "请输入确认密码", trigger: "blur" },
          { min: 6, message: "最少6个字符", trigger: "blur" },
          {
            validator: (rule, value, callback) => {
              if (value !== this.pwdModify.newPassword) {
                callback(new Error("两次密码不一致"));
              } else {
                callback();
              }
            },
            trigger: "blur",
          },
        ],
      },
    };
  },
  components: {
    ChooseImg,
  },
  computed: {
    ...mapGetters("user", ["userInfo", "token"]),
  },
  methods: {
    ...mapMutations("user", ["ResetUserInfo"]),
    savePassword() {
      this.$refs.modifyPwdForm.validate((valid) => {
        if (valid) {
          changePassword({
            username: this.userInfo.userName,
            password: this.pwdModify.password,
            newPassword: this.pwdModify.newPassword,
          }).then((res) => {
            if (res.code == 0) {
              this.$message.success("修改密码成功！");
            }
            this.showPassword = false;
          });
        }
      });
    },
    clearPassword() {
      this.pwdModify = {
        password: "",
        newPassword: "",
        confirmPassword: "",
      };
      this.$refs.modifyPwdForm?.clearValidate();
    },
    openChooseImg() {
      this.$refs.chooseImg.open();
    },
    async enterImg(url) {
      const res = await setUserInfo({ headerImg: url, ID: this.userInfo.ID });
      if (res.code == 0) {
        this.ResetUserInfo({ headerImg: url });
        this.$message({
          type: "success",
          message: "设置成功",
        });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.person-page {
  color: var(--gs-dark-text, #e5e7eb);
}

.user-card {
  min-height: 280px;
  padding: 28px 20px;
  text-align: center;
  background: rgba(15, 23, 42, 0.55);
  border: 1px solid rgba(148, 163, 184, 0.12);
  border-radius: 10px;
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.2);

  .nickname {
    margin: 18px 0 8px;
    font-size: 22px;
    font-weight: 600;
    color: #f1f5f9;
  }

  .username {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    margin: 0;
    color: #94a3b8;
    font-size: 14px;
  }
}

.user-panel {
  min-height: 280px;
  padding: 8px 18px 18px;
  background: rgba(15, 23, 42, 0.55);
  border: 1px solid rgba(148, 163, 184, 0.12);
  border-radius: 10px;
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.2);
}

.account-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 4px 6px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);

  .title {
    margin: 0 0 6px;
    font-size: 16px;
    color: #e5e7eb;
  }

  .desc {
    margin: 0;
    font-size: 13px;
    color: #94a3b8;
  }
}

.user-headpic-update {
  width: 120px;
  height: 120px;
  line-height: 120px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  border-radius: 16px;
  background-color: #1e293b;
  overflow: hidden;

  &:hover {
    color: #fff;
    background-blend-mode: multiply;

    .update {
      color: #fff;
      background: rgba(15, 23, 42, 0.55);
    }
  }

  .update {
    height: 120px;
    width: 120px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 4px;
    text-align: center;
    color: transparent;
    cursor: pointer;
  }
}
</style>

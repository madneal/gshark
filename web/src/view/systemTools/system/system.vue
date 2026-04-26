<template>
  <div class="system-config-page">
    <el-form :model="config" ref="form" class="system-config-form" label-position="top">
      <section class="config-section">
        <div class="section-header">
          <h2>系统配置</h2>
        </div>
        <div class="field-grid two">
          <el-form-item label="环境值">
            <el-input v-model="config.system.env"></el-input>
          </el-form-item>
          <el-form-item label="端口值">
            <el-input v-model.number="config.system.addr"></el-input>
          </el-form-item>
          <el-form-item class="inline-control" label="数据初始化">
            <el-checkbox v-model="config.system.needInitData">开启</el-checkbox>
          </el-form-item>
          <el-form-item class="inline-control" label="多点登录拦截">
            <el-checkbox v-model="config.system.useMultipoint">开启</el-checkbox>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>JWT签名</h2>
        </div>
        <div class="field-grid one">
          <el-form-item label="jwt签名">
            <el-input v-model="config.jwt.signingKey"></el-input>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>Zap日志配置</h2>
        </div>
        <div class="field-grid two">
          <el-form-item label="级别">
            <el-input v-model.number="config.zap.level"></el-input>
          </el-form-item>
          <el-form-item label="输出">
            <el-input v-model="config.zap.format"></el-input>
          </el-form-item>
          <el-form-item label="日志前缀">
            <el-input v-model="config.zap.prefix"></el-input>
          </el-form-item>
          <el-form-item label="日志文件夹">
            <el-input v-model="config.zap.director"></el-input>
          </el-form-item>
          <el-form-item label="软链接名称">
            <el-input v-model="config.zap.linkName"></el-input>
          </el-form-item>
          <el-form-item label="编码级">
            <el-input v-model="config.zap.encodeLevel"></el-input>
          </el-form-item>
          <el-form-item label="栈名">
            <el-input v-model="config.zap.stacktraceKey"></el-input>
          </el-form-item>
          <el-form-item class="inline-control" label="显示行">
            <el-checkbox v-model="config.zap.showLine">开启</el-checkbox>
          </el-form-item>
          <el-form-item class="inline-control wide" label="输出控制台">
            <el-checkbox v-model="config.zap.logInConsole">开启</el-checkbox>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>企业微信机器人</h2>
          <el-button @click="sendBot" size="small" type="primary">测试</el-button>
        </div>
        <div class="field-grid one">
          <el-form-item class="inline-control" label="是否启用">
            <el-switch
              v-model="config.wechat.enable"
              active-color="#13ce66"
              inactive-color="#ff4949"
            ></el-switch>
          </el-form-item>
          <el-form-item label="企业微信">
            <el-input v-model="config.wechat.url" placeholder="请填写企业微信机器人webhook地址"></el-input>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>邮箱配置</h2>
          <el-button @click="email" size="small" type="primary">测试邮件</el-button>
        </div>
        <div class="field-grid two">
          <el-form-item class="inline-control wide" label="是否启用">
            <el-switch
              v-model="config.email.enable"
              active-color="#13ce66"
              inactive-color="#ff4949"
            ></el-switch>
          </el-form-item>
          <el-form-item label="接收者邮箱">
            <el-input v-model="config.email.to" placeholder="可多个，以逗号分隔"></el-input>
          </el-form-item>
          <el-form-item label="端口">
            <el-input v-model.number="config.email.port" placeholder="请输入端口号"></el-input>
          </el-form-item>
          <el-form-item label="发送者邮箱">
            <el-input v-model="config.email.from"></el-input>
          </el-form-item>
          <el-form-item label="host">
            <el-input v-model="config.email.host"></el-input>
          </el-form-item>
          <el-form-item label="secret">
            <el-input v-model="config.email.secret"></el-input>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>Casbin配置</h2>
        </div>
        <div class="field-grid one">
          <el-form-item label="模型地址">
            <el-input v-model="config.casbin.modelPath"></el-input>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>验证码配置</h2>
        </div>
        <div class="field-grid three">
          <el-form-item label="keyLong">
            <el-input v-model.number="config.captcha.keyLong"></el-input>
          </el-form-item>
          <el-form-item label="imgWidth">
            <el-input v-model.number="config.captcha.imgWidth"></el-input>
          </el-form-item>
          <el-form-item label="imgHeight">
            <el-input v-model.number="config.captcha.imgHeight"></el-input>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>MySQL数据库配置</h2>
        </div>
        <div class="field-grid two">
          <el-form-item label="username">
            <el-input v-model="config.mysql.username"></el-input>
          </el-form-item>
          <el-form-item label="password">
            <el-input v-model="config.mysql.password" type="password" show-password></el-input>
          </el-form-item>
          <el-form-item label="path">
            <el-input v-model="config.mysql.path"></el-input>
          </el-form-item>
          <el-form-item label="dbname">
            <el-input v-model="config.mysql.dbname"></el-input>
          </el-form-item>
          <el-form-item label="maxIdleConns">
            <el-input v-model.number="config.mysql.maxIdleConns"></el-input>
          </el-form-item>
          <el-form-item label="maxOpenConns">
            <el-input v-model.number="config.mysql.maxOpenConns"></el-input>
          </el-form-item>
          <el-form-item class="inline-control wide" label="logMode">
            <el-checkbox v-model="config.mysql.logMode">启用日志模式</el-checkbox>
          </el-form-item>
        </div>
      </section>

      <section class="config-section">
        <div class="section-header">
          <h2>本地上传配置</h2>
        </div>
        <div class="field-grid one">
          <el-form-item label="本地文件路径">
            <el-input v-model="config.local.path"></el-input>
          </el-form-item>
        </div>
      </section>

      <div class="form-actions">
        <el-button @click="update" type="primary">立即更新</el-button>
        <el-button @click="reload">重启服务（开发中）</el-button>
      </div>
    </el-form>
  </div>
</template>

<script>
import { getSystemConfig, setSystemConfig } from "@/api/system";
import { emailTest, botTest } from "@/api/email";
export default {
  name: "Config",
  data() {
    return {
      config: {
        system: {},
        jwt: {},
        casbin: {},
        mysql: {},
        captcha: {},
        zap: {},
        local: {},
        email: {},
        wechat: {}
      }
    };
  },
  async created() {
    await this.initForm();
  },
  methods: {
    async initForm() {
      const res = await getSystemConfig();
      if (res.code === 0) {
        this.config = res.data.config;
      }
    },
    reload() {},
    async update() {
      const res = await setSystemConfig({ config: this.config });
      if (res.code === 0) {
        this.$message({
          type: "success",
          message: "配置文件设置成功"
        });
        await this.initForm();
      }
    },
    async email() {
      const res = await emailTest();
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "邮件发送成功"
        });
        await this.initForm();
      } else {
        this.$message({
          type: "error",
          message: "邮件发送失败"
        });
      }
    },
    async sendBot() {
      const res = await botTest();
      if (res.code === 0) {
        this.$message.success("消息发送成功");
        await this.initForm();
      } else {
        this.$message.error("消息发送失败");
      }
    }
  }
};
</script>
<style scoped lang="scss">
.system-config-page {
  color: var(--gs-dark-text);
}

.system-config-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  max-width: 1040px;
  margin: 0 auto;
}

.config-section {
  padding: 0 0 18px;
  border-bottom: 1px solid var(--gs-dark-border-soft);

  &:last-of-type {
    border-bottom: none;
  }
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 38px;
  margin-bottom: 14px;

  h2 {
    position: relative;
    margin: 0;
    padding-left: 12px;
    color: #ffffff;
    font-size: 16px;
    font-weight: 600;
    line-height: 1.35;

    &::before {
      content: "";
      position: absolute;
      top: 3px;
      bottom: 3px;
      left: 0;
      width: 3px;
      border-radius: 3px;
      background: var(--gs-dark-primary);
    }
  }
}

.field-grid {
  display: grid;
  gap: 14px 18px;

  &.one {
    grid-template-columns: minmax(0, 1fr);
  }

  &.two {
    grid-template-columns: repeat(2, minmax(260px, 1fr));
  }

  &.three {
    grid-template-columns: repeat(3, minmax(180px, 1fr));
  }
}

.wide {
  grid-column: 1 / -1;
}

.inline-control {
  :deep(.el-form-item__content) {
    min-height: 34px;
    align-items: center;
  }
}

.form-actions {
  position: sticky;
  bottom: 0;
  z-index: 2;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin: 4px -18px -18px;
  padding: 16px 18px 0;
  background: linear-gradient(180deg, rgba(24, 34, 53, 0), var(--gs-dark-panel) 34%);
}

:deep(.el-form-item) {
  margin-bottom: 0;
}

:deep(.el-form-item__label) {
  height: auto;
  margin-bottom: 6px;
  padding: 0;
  color: var(--gs-dark-muted);
  font-size: 13px;
  font-weight: 600;
  line-height: 1.35;
}

:deep(.el-form-item__content) {
  min-width: 0;
}

:deep(.el-input),
:deep(.el-select),
:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-checkbox) {
  height: 32px;
}

@media screen and (max-width: 1100px) {
  .system-config-form {
    max-width: none;
  }

  .field-grid.two,
  .field-grid.three {
    grid-template-columns: repeat(2, minmax(220px, 1fr));
  }
}

@media screen and (max-width: 700px) {
  .system-config-form {
    gap: 16px;
  }

  .field-grid,
  .field-grid.two,
  .field-grid.three {
    grid-template-columns: minmax(0, 1fr);
  }

  .section-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .form-actions {
    flex-wrap: wrap;

    :deep(.el-button) {
      flex: 1 1 160px;
    }
  }
}
</style>

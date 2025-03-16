<template>
  <div class="system">
    <el-form :model="config" label-width="100px" ref="form" class="system">
      <!--  System start  -->
      <h2>系统配置</h2>
      <el-row>
        <el-col :lg="12">
          <el-form-item label="环境值">
            <el-input v-model="config.system.env"></el-input>
          </el-form-item>
        </el-col>
        <el-col :lg="12">
          <el-form-item label="端口值">
            <el-input v-model.number="config.system.addr"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :lg="12">
          <el-form-item label="数据初始化">
            <el-checkbox v-model="config.system.needInitData">开启</el-checkbox>
          </el-form-item>
        </el-col>
        <el-col :lg="12">
          <el-form-item label="多点登录拦截">
            <el-checkbox v-model="config.system.useMultipoint">开启</el-checkbox>
          </el-form-item>
        </el-col>
      </el-row>
      <!--  System end  -->

      <!--  JWT start  -->
      <h2>jwt签名</h2>
      <el-form-item label="jwt签名">
        <el-input v-model="config.jwt.signingKey"></el-input>
      </el-form-item>
      <!--  JWT end  -->

      <!--  Zap start  -->
      <h2>Zap日志配置</h2>
      <el-row>
        <el-col :span="12">
          <el-form-item label="级别">
            <el-input v-model.number="config.zap.level"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="输出">
            <el-input v-model="config.zap.format"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="日志前缀">
            <el-input v-model="config.zap.prefix"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="日志文件夹">
            <el-input v-model="config.zap.director"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="软链接名称">
            <el-input v-model="config.zap.linkName"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="编码级">
            <el-input v-model="config.zap.encodeLevel"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="栈名">
            <el-input v-model="config.zap.stacktraceKey"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="显示行">
            <el-checkbox v-model="config.zap.showLine"></el-checkbox>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="24">
          <el-form-item label="输出控制台">
            <el-checkbox v-model="config.zap.logInConsole"></el-checkbox>
          </el-form-item>
        </el-col>
      </el-row>
      <!--  Zap end  -->

      <h2>企业微信机器人</h2>
      <el-form-item label="是否启用">
        <el-switch
            v-model="config.wechat.enable"
            active-color="#13ce66"
            inactive-color="#ff4949">
        </el-switch>
      </el-form-item>
      <el-form-item label="企业微信">
        <el-input v-model="config.wechat.url" placeholder="请填写企业微信机器人webhook地址"></el-input>
      </el-form-item>
      <el-form-item label="测试机器人">
        <el-button @click="sendBot">测试</el-button>
      </el-form-item>

      <!--  Email start  -->
      <h2>邮箱配置</h2>
      <el-row>
        <el-col :span="24">
          <el-form-item label="是否启用">
            <el-switch
                v-model="config.email.enable"
                active-color="#13ce66"
                inactive-color="#ff4949">
            </el-switch>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="接收者邮箱">
            <el-input
                v-model="config.email.to"
                placeholder="可多个，以逗号分隔">
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="端口">
            <el-input
                v-model.number="config.email.port"
                placeholder="请输入端口号">
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="发送者邮箱">
            <el-input v-model="config.email.from"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="host">
            <el-input v-model="config.email.host"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="secret">
            <el-input v-model="config.email.secret"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="测试邮件">
            <el-button type="primary" @click="email">测试邮件</el-button>
          </el-form-item>
        </el-col>
      </el-row>
      <!--  Email end  -->

      <!--  Casbin start  -->
      <h2>casbin配置</h2>
      <el-form-item label="模型地址">
        <el-input v-model="config.casbin.modelPath"></el-input>
      </el-form-item>
      <!--  Casbin end  -->

      <!--  Captcha start  -->
      <h2>验证码配置</h2>
      <el-row>
        <el-col :span="12">
          <el-form-item label="keyLong">
            <el-input v-model.number="config.captcha.keyLong"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="imgWidth">
            <el-input v-model.number="config.captcha.imgWidth"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="imgHeight">
            <el-input v-model.number="config.captcha.imgHeight"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <!-- Optional: Empty column or additional field -->
        </el-col>
      </el-row>
      <!--  Captcha end  -->

      <!--  dbType start  -->
        <h2>mysql 数据库配置</h2>
      <el-row>
        <el-col :span="12">
          <el-form-item label="username">
            <el-input v-model="config.mysql.username"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="password">
            <el-input
                v-model="config.mysql.password"
                type="password"
                show-password>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="path">
            <el-input v-model="config.mysql.path"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="dbname">
            <el-input v-model="config.mysql.dbname"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="12">
          <el-form-item label="maxIdleConns">
            <el-input v-model.number="config.mysql.maxIdleConns"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="maxOpenConns">
            <el-input v-model.number="config.mysql.maxOpenConns"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="24">
          <el-form-item label="logMode">
            <el-checkbox v-model="config.mysql.logMode">启用日志模式</el-checkbox>
          </el-form-item>
        </el-col>
      </el-row>

      <!--  dbType end  -->

      <!--  ossType start  -->
        <h2>本地上传配置</h2>
        <el-form-item label="本地文件路径">
          <el-input v-model="config.local.path"></el-input>
        </el-form-item>
      <!--  ossType end  -->

      <el-form-item>
        <el-button @click="update" type="primary">立即更新</el-button>
        <el-button @click="reload" type="primary">重启服务（开发中）</el-button>
      </el-form-item>
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
<style lang="scss">
.system {
  h2 {
    padding: 10px;
    margin: 10px 0;
    font-size: 16px;
    box-shadow: -4px 1px 3px 0px #e7e8e8;
  }
}
.el-input {
  width: 50%;
}
</style>

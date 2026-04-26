<template>
  <div class="server-state-page">
    <el-empty v-if="stateError && !hasStateData" :description="stateError"></el-empty>
    <el-row :gutter="15" class="system_state">
      <el-col :span="12">
        <el-card v-if="state.os" class="card_item">
          <template #header>Runtime</template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">os:</el-col>
              <el-col :span="12" v-text="state.os.goos"></el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">cpu nums:</el-col>
              <el-col :span="12" v-text="state.os.numCpu"></el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">compiler:</el-col>
              <el-col :span="12" v-text="state.os.compiler"></el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">go version:</el-col>
              <el-col :span="12" v-text="state.os.goVersion"></el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">goroutine nums:</el-col>
              <el-col :span="12" v-text="state.os.numGoroutine"></el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card v-if="state.disk" class="card_item">
          <template #header>Disk</template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-row :gutter="10">
                  <el-col :span="12">total (MB)</el-col>
                  <el-col :span="12" v-text="state.disk.totalMb"></el-col>
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (MB)</el-col>
                  <el-col :span="12" v-text="state.disk.usedMb"></el-col>
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">total (GB)</el-col>
                  <el-col :span="12" v-text="state.disk.totalGb"></el-col>
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (GB)</el-col>
                  <el-col :span="12" v-text="state.disk.usedGb"></el-col>
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-progress
                  type="dashboard"
                  :percentage="state.disk.usedPercent"
                  :color="colors"
                ></el-progress>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="15" class="system_state">
      <el-col :span="12">
        <el-card
          class="card_item"
          v-if="state.cpu"
          :body-style="{ height: '180px', 'overflow-y': 'scroll' }"
        >
          <template #header>CPU</template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">physical number of cores:</el-col>
              <el-col :span="12" v-text="state.cpu.cores"> </el-col>
            </el-row>
            <template v-for="(item, index) in state.cpu.cpus" :key="index">
              <el-row :gutter="10">
                <el-col :span="12">core {{ index }}:</el-col>
                <el-col :span="12"
                  ><el-progress
                    type="line"
                    :percentage="+item.toFixed(0)"
                    :color="colors"
                  ></el-progress
                ></el-col>
              </el-row>
            </template>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="card_item" v-if="state.ram">
          <template #header>Ram</template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-row :gutter="10">
                  <el-col :span="12">total (MB)</el-col>
                  <el-col :span="12" v-text="state.ram.totalMb"></el-col>
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (MB)</el-col>
                  <el-col :span="12" v-text="state.ram.usedMb"></el-col>
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">total (GB)</el-col>
                  <el-col :span="12" v-text="state.ram.totalMb / 1024"></el-col>
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (GB)</el-col>
                  <el-col
                    :span="12"
                    v-text="(state.ram.usedMb / 1024).toFixed(2)"
                  ></el-col>
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-progress
                  type="dashboard"
                  :percentage="state.ram.usedPercent"
                  :color="colors"
                ></el-progress>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getSystemState } from "@/api/system.js";
export default {
  name: "State",
  data() {
    return {
      timer:null,
      state: {},
      stateError: "",
      colors: [
        { color: "#5cb87a", percentage: 20 },
        { color: "#e6a23c", percentage: 40 },
        { color: "#f56c6c", percentage: 80 },
      ],
    };
  },
  created() { 
    this.reload();
    this.timer = setInterval(() => {
      this.reload();
    }, 1000*10);
  },
  beforeUnmount(){
    clearInterval(this.timer)
    this.timer = null
  },
  computed: {
    hasStateData() {
      return Boolean(this.state.os || this.state.disk || this.state.cpu || this.state.ram)
    },
  },
  methods: {
    async reload() {
      const res = await getSystemState();
      this.state = res.data?.server || {};
      this.stateError = this.hasStateData ? "" : (res.msg || "服务器状态暂不可用");
    },
  },
};
</script>

<style lang="scss">
.server-state-page {
  .system_state {
    padding: 0;
    margin: 0 -8px 16px;

    &:last-child {
      margin-bottom: 0;
    }

    > .el-col {
      padding-bottom: 16px;
    }
  }

  .card_item {
    min-height: 280px;
    height: 100%;
    background: var(--gs-dark-panel);
    border: 1px solid var(--gs-dark-border-soft);
    border-radius: 8px;
    color: var(--gs-dark-text);
    overflow: hidden;

    .el-card__header {
      padding: 14px 18px;
      color: #ffffff;
      font-size: 15px;
      font-weight: 600;
      line-height: 1.3;
      background: rgba(15, 23, 42, .34);
      border-bottom: 1px solid var(--gs-dark-border-soft);
    }

    .el-card__body {
      padding: 16px 18px;
    }

    .el-row {
      padding: 5px 0;
      align-items: center;
      row-gap: 8px;
    }

    .el-col {
      min-width: 0;
      line-height: 1.45;
      color: var(--gs-dark-text);
      overflow-wrap: anywhere;
    }

    .el-col:nth-child(odd) {
      color: var(--gs-dark-muted);
    }

    .el-progress {
      width: 100%;
    }

    .el-progress-bar__outer {
      background-color: var(--gs-dark-bg-soft);
    }

    .el-progress__text {
      min-width: 40px;
      color: var(--gs-dark-text);
    }

    .el-progress--dashboard {
      display: flex;
      justify-content: center;
      padding: 4px 0;
    }
  }

  .el-empty {
    padding: 36px 0;
  }
}

@media screen and (max-width: 900px) {
  .server-state-page {
    .system_state {
      margin-bottom: 0;

      > .el-col {
        flex: 0 0 100%;
        max-width: 100%;
      }
    }

    .card_item {
      min-height: 0;
    }
  }
}

@media screen and (max-width: 520px) {
  .server-state-page {
    .card_item {
      .el-card__header {
        padding: 12px 14px;
      }

      .el-card__body {
        padding: 12px 14px;
      }

      .el-row .el-col {
        flex: 0 0 100%;
        max-width: 100%;
      }
    }
  }
}
</style>

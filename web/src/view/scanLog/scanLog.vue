<template>
  <div class="scan-log-page">
    <section class="scan-summary">
      <div class="summary-copy">
        <span class="summary-label">当前扫描周期</span>
        <strong>{{ shortCycleId(overview.cycleId) }}</strong>
        <span class="last-activity">最近活动 {{ formatDate(overview.lastActivityAt) || "--" }}</span>
      </div>
      <div class="summary-progress">
        <div class="progress-meta">
          <span>{{ overview.completed || 0 }}/{{ overview.total || providerOrder.length }} 已完成</span>
          <span>{{ overview.progress || 0 }}%</span>
        </div>
        <el-progress
          :percentage="overview.progress || 0"
          :stroke-width="8"
          :show-text="false"
          :status="overview.abnormal > 0 ? 'exception' : undefined"
        />
      </div>
      <div class="summary-metrics">
        <div>
          <span>运行中</span>
          <strong>{{ overview.running || 0 }}</strong>
        </div>
        <div>
          <span>异常</span>
          <strong :class="{ danger: overview.abnormal > 0 }">{{ overview.abnormal || 0 }}</strong>
        </div>
      </div>
    </section>

    <section class="provider-grid" aria-label="扫描任务状态">
      <article v-for="provider in providerStates" :key="provider.provider" class="provider-card">
        <header>
          <span class="provider-name">{{ providerLabel(provider.provider) }}</span>
          <el-tag :type="statusTagType(provider)" effect="dark" size="small">
            {{ statusLabel(provider) }}
          </el-tag>
        </header>
        <p>{{ provider.message || "等待扫描数据" }}</p>
        <footer>
          <span>{{ providerTime(provider) }}</span>
          <span v-if="provider.durationMs">{{ formatDuration(provider.durationMs) }}</span>
        </footer>
      </article>
    </section>

    <div class="search-term scan-toolbar">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="平台">
          <el-select v-model="searchInfo.provider" clearable placeholder="全部平台" style="width: 150px">
            <el-option
              v-for="provider in providerOrder"
              :key="provider"
              :label="providerLabel(provider)"
              :value="provider"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="全部状态" style="width: 150px">
            <el-option v-for="status in statusOptions" :key="status.value" :label="status.label" :value="status.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="applyFilters">查询</el-button>
          <el-tooltip content="重置筛选" placement="top">
            <el-button :icon="RefreshLeft" circle aria-label="重置筛选" @click="resetFilters" />
          </el-tooltip>
        </el-form-item>
      </el-form>
      <div class="toolbar-actions">
        <span>自动刷新</span>
        <el-switch v-model="autoRefresh" @change="syncTimer" />
        <el-tooltip content="立即刷新" placement="top">
          <el-button :icon="Refresh" circle :loading="loading" aria-label="立即刷新" @click="refreshAll" />
        </el-tooltip>
      </div>
    </div>

    <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%">
      <el-table-column label="开始时间" min-width="170">
        <template #default="scope">{{ formatDate(scope.row.startedAt || scope.row.CreatedAt) }}</template>
      </el-table-column>
      <el-table-column label="平台" min-width="110">
        <template #default="scope">{{ providerLabel(scope.row.provider) }}</template>
      </el-table-column>
      <el-table-column label="状态" min-width="110">
        <template #default="scope">
          <el-tag :type="statusTagType(scope.row)" effect="dark" size="small">
            {{ statusLabel(scope.row) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="耗时" min-width="100">
        <template #default="scope">{{ formatDuration(scope.row.durationMs) }}</template>
      </el-table-column>
      <el-table-column label="周期" min-width="130">
        <template #default="scope">
          <span class="cycle-id" :title="scope.row.cycleId">{{ shortCycleId(scope.row.cycleId) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="信息" prop="message" min-width="320" show-overflow-tooltip />
      <el-table-column label="最近心跳" min-width="170">
        <template #default="scope">{{ formatDate(scope.row.heartbeatAt) || "--" }}</template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<script>
import { markRaw } from "vue";
import { Refresh, RefreshLeft, Search } from "@element-plus/icons-vue";
import { getScanLogList, getScanLogOverview } from "@/api/scanLog";
import { formatDate } from "@/utils/date";

const EMPTY_OVERVIEW = {
  cycleId: "",
  logs: [],
  total: 0,
  completed: 0,
  running: 0,
  abnormal: 0,
  progress: 0,
  lastActivityAt: null,
};

export default {
  name: "ScanLog",
  data() {
    return {
      Refresh: markRaw(Refresh),
      RefreshLeft: markRaw(RefreshLeft),
      Search: markRaw(Search),
      providerOrder: ["gitlab", "sourcegraph", "github", "gobuster", "postman"],
      statusOptions: [
        { value: "pending", label: "排队中" },
        { value: "running", label: "运行中" },
        { value: "success", label: "成功" },
        { value: "skipped", label: "已跳过" },
        { value: "failed", label: "失败" },
        { value: "timeout", label: "超时" },
        { value: "interrupted", label: "已中断" },
      ],
      overview: { ...EMPTY_OVERVIEW },
      searchInfo: { provider: "", status: "" },
      tableData: [],
      page: 1,
      pageSize: 30,
      total: 0,
      loading: false,
      autoRefresh: true,
      timer: null,
    };
  },
  computed: {
    providerStates() {
      const logs = new Map((this.overview.logs || []).map((item) => [item.provider, item]));
      return this.providerOrder.map((provider) => logs.get(provider) || {
        provider,
        status: "pending",
        message: "等待扫描数据",
      });
    },
  },
  methods: {
    formatDate,
    providerLabel(provider) {
      const labels = {
        gitlab: "GitLab",
        sourcegraph: "Sourcegraph",
        searchcode: "Searchcode",
        github: "GitHub",
        gobuster: "Domain",
        postman: "Postman",
      };
      return labels[provider] || provider || "--";
    },
    statusLabel(item) {
      if (item.stale) return "心跳异常";
      const option = this.statusOptions.find((status) => status.value === item.status);
      return option ? option.label : item.status || "未知";
    },
    statusTagType(item) {
      if (item.stale || ["failed", "timeout", "interrupted"].includes(item.status)) return "danger";
      if (item.status === "success") return "success";
      if (item.status === "running") return "primary";
      if (item.status === "skipped") return "info";
      return "warning";
    },
    shortCycleId(cycleId) {
      return cycleId ? cycleId.slice(0, 8) : "--";
    },
    formatDuration(durationMs) {
      if (!durationMs) return "--";
      if (durationMs < 1000) return `${durationMs} ms`;
      if (durationMs < 60000) return `${(durationMs / 1000).toFixed(1)} s`;
      const minutes = Math.floor(durationMs / 60000);
      const seconds = Math.floor((durationMs % 60000) / 1000);
      return `${minutes}m ${seconds}s`;
    },
    providerTime(provider) {
      if (provider.status === "running") return `心跳 ${formatDate(provider.heartbeatAt) || "--"}`;
      return formatDate(provider.finishedAt || provider.CreatedAt) || "--";
    },
    async refreshAll() {
      if (this.loading) return;
      this.loading = true;
      try {
        const [overviewRes, listRes] = await Promise.all([
          getScanLogOverview(),
          getScanLogList({
            page: this.page,
            pageSize: this.pageSize,
            ...this.searchInfo,
          }),
        ]);
        if (overviewRes.code === 0) this.overview = overviewRes.data || { ...EMPTY_OVERVIEW };
        if (listRes.code === 0) {
          this.tableData = listRes.data.list || [];
          this.total = listRes.data.total || 0;
          this.page = listRes.data.page || this.page;
          this.pageSize = listRes.data.pageSize || this.pageSize;
        }
      } finally {
        this.loading = false;
      }
    },
    applyFilters() {
      this.page = 1;
      this.refreshAll();
    },
    resetFilters() {
      this.searchInfo = { provider: "", status: "" };
      this.page = 1;
      this.refreshAll();
    },
    handleCurrentChange(page) {
      this.page = page;
      this.refreshAll();
    },
    handleSizeChange(pageSize) {
      this.pageSize = pageSize;
      this.page = 1;
      this.refreshAll();
    },
    syncTimer() {
      if (this.timer) clearInterval(this.timer);
      this.timer = this.autoRefresh ? setInterval(this.refreshAll, 5000) : null;
    },
  },
  created() {
    this.refreshAll();
    this.syncTimer();
  },
  beforeUnmount() {
    if (this.timer) clearInterval(this.timer);
  },
};
</script>

<style scoped>
.scan-log-page {
  min-width: 0;
}

.scan-summary {
  display: grid;
  grid-template-columns: minmax(220px, 0.8fr) minmax(280px, 1.4fr) auto;
  gap: 28px;
  align-items: center;
  padding: 18px 2px 20px;
  border-bottom: 1px solid var(--el-border-color);
}

.summary-copy,
.summary-progress {
  min-width: 0;
}

.summary-copy {
  display: grid;
  gap: 4px;
}

.summary-label,
.last-activity,
.progress-meta,
.summary-metrics span,
.provider-card footer {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.summary-copy strong {
  font-size: 22px;
  letter-spacing: 0;
}

.progress-meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.summary-metrics {
  display: flex;
  gap: 28px;
}

.summary-metrics div {
  display: grid;
  gap: 2px;
  min-width: 58px;
}

.summary-metrics strong {
  font-size: 24px;
  font-variant-numeric: tabular-nums;
}

.summary-metrics strong.danger {
  color: var(--el-color-danger);
}

.provider-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(150px, 1fr));
  gap: 12px;
  margin: 18px 0;
}

.provider-card {
  min-width: 0;
  padding: 14px;
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  background: var(--el-bg-color-overlay);
}

.provider-card header,
.provider-card footer,
.scan-toolbar,
.toolbar-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.provider-name {
  overflow: hidden;
  font-weight: 600;
  text-overflow: ellipsis;
}

.provider-card p {
  height: 40px;
  margin: 12px 0;
  overflow: hidden;
  color: var(--el-text-color-regular);
  font-size: 13px;
  line-height: 20px;
}

.provider-card footer span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.scan-toolbar {
  flex-wrap: wrap;
}

.scan-toolbar :deep(.el-form-item) {
  margin-bottom: 0;
}

.toolbar-actions {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.cycle-id {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

@media (max-width: 1100px) {
  .provider-grid {
    grid-template-columns: repeat(3, minmax(160px, 1fr));
  }
}

@media (max-width: 760px) {
  .scan-summary {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .provider-grid {
    grid-template-columns: 1fr;
  }

  .summary-metrics {
    justify-content: flex-start;
  }

  .scan-toolbar,
  .toolbar-actions {
    align-items: flex-start;
  }
}
</style>

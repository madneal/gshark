<template>
  <div class="search-result-page">
    <div class="search-term page-toolbar">
      <el-form :inline="true" :model="searchInfo" class="page-toolbar__filters">
        <el-form-item label="搜索条件">
          <el-input
            placeholder="仓库名称 | 匹配内容"
            clearable
            v-model="searchInfo.query"
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            placeholder="关键词"
            clearable
            v-model="searchInfo.keyword"
            style="width: 140px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchInfo.status"
            clearable
            placeholder="全部"
            style="width: 120px"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="仅二次关键词">
          <el-switch
            v-model="searchInfo.onlySecKeyword"
            @change="secKeywordChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="exportResult">导出</el-button>
        </el-form-item>
      </el-form>

      <div class="page-toolbar__actions">
        <el-button
          type="success"
          :disabled="!hasSelection"
          @click="confirmBulk(false)"
        >
          批量确认
        </el-button>
        <el-button
          type="danger"
          :disabled="!hasSelection"
          @click="confirmBulk(true)"
        >
          批量忽略
        </el-button>
        <el-button @click="startAITask">启动 AI 分析</el-button>
        <el-button
          :disabled="taskBtnDisable"
          @click="confirmFilterTask"
        >
          {{ taskButtonTxt }}
        </el-button>
      </div>
    </div>

    <el-table
      :data="tableData"
      @selection-change="handleSelectionChange"
      border
      ref="multipleTable"
      stripe
      style="width: 100%"
      tooltip-effect="dark"
      empty-text="暂无匹配结果"
    >
      <el-table-column type="selection" width="48" />
      <el-table-column label="ID" prop="ID" width="64" />
      <el-table-column label="文件" min-width="180">
        <template #default="scope">
          <a
            class="table-link"
            :href="scope.row.url"
            target="_blank"
            rel="noopener noreferrer"
          >
            {{ scope.row.repo + "/" + scope.row.path }}
          </a>
        </template>
      </el-table-column>
      <el-table-column label="匹配内容" min-width="320">
        <template #default="scope">
          <pre
            v-if="scope.row.text_matches"
            class="match-snippet"
            v-html="fragmentsFilter(scope.row.text_matches)"
          />
          <pre v-else-if="scope.row.matches" class="match-snippet">{{
            scope.row.matches
          }}</pre>
        </template>
      </el-table-column>
      <el-table-column label="关键词" prop="keyword" width="110">
        <template #default="scope">
          <span v-if="scope.row.keyword" class="keyword-chip">{{
            scope.row.keyword
          }}</span>
        </template>
      </el-table-column>
      <el-table-column label="二级关键词" prop="sec_keyword" width="120">
        <template #default="scope">
          <span v-if="scope.row.sec_keyword" class="keyword-chip">{{
            scope.row.sec_keyword
          }}</span>
        </template>
      </el-table-column>
      <el-table-column label="日期" width="110">
        <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag
            class="status-tag"
            size="small"
            :type="statusTagType(scope.row.status)"
            effect="dark"
          >
            {{ statusFilter(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="170" fixed="right">
        <template #default="scope">
          <el-button
            size="small"
            type="primary"
            :disabled="scope.row.status === 1"
            @click="updateSearchResult(scope.row, 1)"
          >
            确认
          </el-button>
          <el-button
            size="small"
            type="danger"
            :disabled="scope.row.status === 2"
            @click="updateSearchResult(scope.row, 2)"
          >
            忽略
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    />
  </div>
</template>

<script>
import {
  findSearchResult,
  getSearchResultList,
  updateSearchResult,
  updateSearchResultStatusByIds,
  startFilterTask,
  getTaskStatus,
  exportSearchResult,
  startAITask,
} from "@/api/searchResult";
import { formatDate } from "@/utils/date";
import infoList from "@/mixins/infoList";

export default {
  name: "SearchResult",
  mixins: [infoList],
  data() {
    return {
      listApi: getSearchResultList,
      dialogFormVisible: false,
      type: "",
      taskButtonTxt: "启动二次过滤",
      taskBtnDisable: false,
      statusOptions: [
        { label: "未处理", value: 0 },
        { label: "已确认", value: 1 },
        { label: "已忽略", value: 2 },
      ],
      multipleSelection: [],
      formData: {
        repo: "",
        matches: "",
        keyword: "",
        path: "",
        url: "",
        status: 0,
      },
    };
  },
  computed: {
    hasSelection() {
      return this.multipleSelection.length > 0;
    },
  },
  methods: {
    formatDate,
    escapeHtml(value) {
      return String(value)
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#39;");
    },
    statusFilter(val) {
      const statusOptions = {
        0: "未处理",
        1: "已确认",
        2: "已忽略",
      };
      return statusOptions[val] ?? "未知";
    },
    statusTagType(val) {
      if (val === 1) return "success";
      if (val === 2) return "info";
      return "warning";
    },
    fragmentsFilter(val) {
      if (!val) {
        return "";
      }
      if (typeof val === "string") {
        return this.escapeHtml(val);
      }
      const parts = [];
      for (let i = 0; i < val.length; i++) {
        const fragment = val[i].fragment || "";
        if (!val[i].matches || !val[i].matches.length) {
          parts.push(this.escapeHtml(fragment));
        } else {
          const ranges = [];
          val[i].matches.forEach((ele) => {
            if (Array.isArray(ele.indices) && ele.indices.length >= 2) {
              ranges.push([ele.indices[0], ele.indices[1]]);
            }
          });
          ranges.sort((a, b) => a[0] - b[0]);
          let cursor = 0;
          let html = "";
          ranges.forEach(([start, end]) => {
            const safeStart = Math.max(0, start);
            const safeEnd = Math.min(fragment.length, end);
            if (safeEnd <= safeStart || safeStart < cursor) {
              return;
            }
            html += this.escapeHtml(fragment.slice(cursor, safeStart));
            html += `<mark>${this.escapeHtml(
              fragment.slice(safeStart, safeEnd)
            )}</mark>`;
            cursor = safeEnd;
          });
          html += this.escapeHtml(fragment.slice(cursor));
          parts.push(html);
        }
        if (i !== val.length - 1) {
          parts.push("\n<span class=\"match-sep\">────</span>\n");
        }
      }
      return parts.join("");
    },
    onSubmit() {
      this.page = 1;
      this.pageSize = 100;
      this.getTableData();
    },
    async exportResult() {
      await exportSearchResult(this.formData);
    },
    secKeywordChange() {
      this.getTableData();
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    async startAITask() {
      try {
        await startAITask();
        this.$message({ type: "success", message: "已启动 AI 分析任务" });
      } catch (e) {
        this.$message({ type: "error", message: "启动 AI 分析失败" });
      }
    },
    confirmFilterTask() {
      this.$confirm(
        "过滤任务将自动忽略未匹配二次过滤关键词的结果，确定启动吗？",
        "启动二次过滤",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }
      )
        .then(() => this.startFilterTask())
        .catch(() => {});
    },
    async startFilterTask() {
      await startFilterTask();
      const resp = await getTaskStatus();
      if (resp.msg === "running") {
        this.taskButtonTxt = "任务运行中";
        this.taskBtnDisable = true;
      }
      this.$message({ type: "success", message: "二次过滤任务已启动" });
    },
    confirmBulk(isIgnore) {
      if (!this.hasSelection) {
        this.$message({ type: "warning", message: "请选择要操作的数据" });
        return;
      }
      const action = isIgnore ? "忽略" : "确认";
      this.$confirm(
        `确定批量${action}选中的 ${this.multipleSelection.length} 条结果吗？`,
        `批量${action}`,
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: isIgnore ? "warning" : "info",
        }
      )
        .then(() => this.onChange(isIgnore))
        .catch(() => {});
    },
    async onChange(isIgnore) {
      const ids = this.multipleSelection.map((item) => item.ID);
      const res = await updateSearchResultStatusByIds({
        ids,
        status: isIgnore ? 2 : 1,
      });
      if (res.code === 0) {
        this.$message({ type: "success", message: "操作成功" });
        if (this.tableData.length == ids.length) {
          this.page--;
        }
        await this.getTableData();
      }
    },
    async updateSearchResult(row, status) {
      const res = await findSearchResult({ ID: row.ID });
      this.type = "update";
      if (res.code === 0) {
        res.data.searchResult.status = status;
        this.formData = res.data.searchResult;
        const data = {
          repo: this.formData.repo,
          status: status,
        };
        const updateRes = await updateSearchResult(data);
        if (updateRes.code === 0) {
          this.$message({
            type: "success",
            message: status === 1 ? "已确认" : "已忽略",
          });
          await this.getTableData();
        }
      }
    },
  },
  async created() {
    await this.getTableData();
    const resp = await getTaskStatus();
    if (resp.msg === "running") {
      this.taskButtonTxt = "任务运行中";
      this.taskBtnDisable = true;
    } else {
      this.taskButtonTxt = "启动二次过滤";
      this.taskBtnDisable = false;
    }
  },
};
</script>

<style scoped>
.search-result-page .page-toolbar {
  margin-bottom: 14px;
}
</style>

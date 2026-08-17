<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline filter-row">
        <el-form-item label="搜索条件">
          <el-input
            placeholder="仓库名称|匹配内容"
            clearable
            v-model="searchInfo.query"
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            placeholder="搜索条件"
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
        <el-form-item class="filter-actions">
          <div class="filter-actions__btns">
            <el-button type="primary" @click="onSubmit">查询</el-button>
            <el-button @click="exportResult">导出</el-button>
          </div>
        </el-form-item>
      </el-form>

      <div class="action-row">
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
      <el-table-column type="selection" width="55" />
      <el-table-column label="ID" prop="ID" width="50" />
      <el-table-column label="文件" width="180">
        <template #default="scope">
          <a :href="scope.row.url" target="_blank" rel="noopener noreferrer">
            {{ scope.row.repo + "/" + scope.row.path }}
          </a>
        </template>
      </el-table-column>
      <el-table-column label="匹配内容" prop="matches" min-width="320">
        <template #default="scope">
          <pre
            v-if="scope.row.text_matches"
            class="search-result-matches"
            v-html="renderTextMatches(scope.row.text_matches, scope.row.keyword)"
          ></pre>
          <pre
            v-else-if="scope.row.matches"
            class="search-result-matches"
            v-html="renderMatches(scope.row.matches, scope.row.keyword)"
          ></pre>
        </template>
      </el-table-column>
      <el-table-column label="关键词" prop="keyword" width="100" show-overflow-tooltip />
      <el-table-column label="日期" width="100">
        <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag size="small" :type="statusTagType(scope.row.status)" effect="dark">
            {{ statusFilter(scope.row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="170">
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
  exportSearchResult,
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
      // null shows placeholder "全部"; normalizeSearchInfo maps it to -1 for the API
      // (backend skips status filter when Status < 0). Do not bind -1 to the select —
      // it is not in statusOptions and would render as the literal "-1".
      searchInfo: {
        status: null,
      },
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
    // Backend treats Status < 0 as "no status filter". Null/empty from clearable
    // would otherwise omit the param and bind to Go's zero value (0 = 未处理).
    normalizeSearchInfo() {
      const info = { ...this.searchInfo };
      if (info.status === null || info.status === undefined || info.status === "") {
        info.status = -1;
      }
      return info;
    },
    async getTableData(page = this.page, pageSize = this.pageSize) {
      const table = await this.listApi({
        page,
        pageSize,
        ...this.normalizeSearchInfo(),
      });
      if (table.code === 0) {
        this.tableData = table.data.list;
        this.total = table.data.total;
        this.page = table.data.page;
        this.pageSize = table.data.pageSize;
      }
    },
    statusFilter(val) {
      return { 0: "未处理", 1: "已确认", 2: "已忽略" }[val] ?? "未知";
    },
    statusTagType(val) {
      if (val === 1) return "success";
      if (val === 2) return "info";
      return "warning";
    },
    escapeHtml(value) {
      const entities = {
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        '"': "&quot;",
        "'": "&#39;",
      };
      return String(value ?? "").replace(/[&<>"']/g, (character) => entities[character]);
    },
    escapeRegExp(value) {
      return String(value ?? "").replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    },
    highlightKeyword(value, keyword) {
      const text = String(value ?? "");
      if (!keyword) return this.escapeHtml(text);
      const pattern = new RegExp(this.escapeRegExp(keyword), "gi");
      let cursor = 0;
      let rendered = "";
      text.replace(pattern, (match, offset) => {
        rendered += this.escapeHtml(text.slice(cursor, offset));
        rendered += `<mark class="search-result-highlight">${this.escapeHtml(match)}</mark>`;
        cursor = offset + match.length;
        return match;
      });
      return rendered + this.escapeHtml(text.slice(cursor));
    },
    renderTextMatches(value, keyword) {
      if (!value) return "";
      if (typeof value === "string") return this.highlightKeyword(value, keyword);
      if (!Array.isArray(value)) return this.highlightKeyword(value, keyword);

      return value
        .map((item) => {
          const fragment = String(item?.fragment ?? "");
          const ranges = (item?.matches || [])
            .filter((match) => Array.isArray(match?.indices) && match.indices.length >= 2)
            .map((match) => [Number(match.indices[0]), Number(match.indices[1])])
            .filter(([start, end]) => Number.isFinite(start) && Number.isFinite(end))
            .sort((a, b) => a[0] - b[0]);

          if (!ranges.length) return this.escapeHtml(fragment);

          let cursor = 0;
          let rendered = "";
          ranges.forEach(([start, end]) => {
            const safeStart = Math.max(0, start);
            const safeEnd = Math.min(fragment.length, end);
            if (safeEnd <= safeStart || safeStart < cursor) return;
            rendered += this.escapeHtml(fragment.slice(cursor, safeStart));
            rendered += `<mark class="search-result-highlight">${this.escapeHtml(
              fragment.slice(safeStart, safeEnd)
            )}</mark>`;
            cursor = safeEnd;
          });
          rendered += this.escapeHtml(fragment.slice(cursor));
          return rendered;
        })
        .join("\n=====================================\n");
    },
    renderMatches(value, keyword) {
      return this.highlightKeyword(value, keyword);
    },
    onSubmit() {
      this.page = 1;
      this.pageSize = 100;
      this.getTableData();
    },
    async exportResult() {
      try {
        await exportSearchResult(this.normalizeSearchInfo());
      } catch (e) {
        this.$message({ type: "error", message: "导出失败" });
      }
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
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
        const updateRes = await updateSearchResult({
          repo: this.formData.repo,
          status: status,
        });
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
  },
};
</script>

<style scoped>
.filter-row {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  column-gap: 12px;
  row-gap: 8px;
  margin-bottom: 0;
}

.filter-row :deep(.el-form-item) {
  display: inline-flex;
  align-items: center;
  margin-right: 0;
  margin-bottom: 0;
}

.filter-row :deep(.el-form-item__label),
.filter-row :deep(.el-form-item__content) {
  line-height: 32px;
  height: auto;
}

/* no-label actions sit on the same baseline as inputs */
.filter-actions :deep(.el-form-item__content) {
  display: flex;
  align-items: center;
  min-height: 32px;
}

.filter-actions__btns {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.filter-actions__btns :deep(.el-button + .el-button) {
  margin-left: 0;
}

.action-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  padding-top: 10px;
  border-top: 1px solid rgba(148, 163, 184, 0.12);
  margin-top: 10px;
}

.action-row :deep(.el-button + .el-button) {
  margin-left: 0;
}
</style>

<style>
.el-table pre {
  white-space: pre-line;
}

.search-result-matches {
  margin: 0;
  overflow-wrap: anywhere;
}

.search-result-highlight {
  display: inline;
  padding: 0 0.18em;
  border-radius: 0.2em;
  color: #1f2937;
  background: #f6c945;
  box-shadow: 0 0 0 1px rgba(255, 220, 90, 0.35);
  font-weight: 700;
}
</style>

<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-button type="primary" @click="openDialog">新增过滤规则</el-button>
          <el-button
            type="danger"
            :disabled="multipleSelection.length === 0"
            @click="confirmDelete"
          >
            批量删除
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      :data="tableData"
      @selection-change="handleSelectionChange"
      border
      ref="multipleTable"
      stripe
      style="width: 100%"
      tooltip-effect="dark"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="日期" width="180">
        <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
      </el-table-column>
      <el-table-column label="过滤类型" prop="filter_type" width="120" />
      <el-table-column label="过滤种类" prop="filter_class" width="120" />
      <el-table-column label="内容" prop="content" width="120" />
      <el-table-column label="操作">
        <template #default="scope">
          <el-button size="small" type="primary" @click="updateFilter(scope.row)">
            变更
          </el-button>
          <el-button size="small" type="danger" @click="deleteRow(scope.row)">
            删除
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

    <el-dialog
      :before-close="closeDialog"
      v-model="dialogFormVisible"
      title="新增过滤规则（目前仅适用于Github）"
    >
      <el-form :model="formData" label-position="right" label-width="120px">
        <el-form-item label="过滤类型：">
          <el-radio-group v-model="formData.filter_type">
            <el-radio label="whitelist">白名单</el-radio>
            <el-radio label="blacklist">黑名单</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="过滤种类：">
          <el-radio-group v-model="formData.filter_class">
            <el-radio label="extension">文件后缀</el-radio>
            <el-radio label="keyword">关键词</el-radio>
            <el-radio label="sec_keyword">二次关键词</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="内容：">
          <el-input
            v-model="formData.content"
            clearable
            placeholder="仅适用于Github，排除关键词，以,分隔"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button @click="enterDialog" type="primary">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  createFilter,
  deleteFilter,
  deleteFilterByIds,
  updateFilter,
  findFilter,
  getFilterList,
} from "@/api/filter";
import { formatDate } from "@/utils/date";
import infoList from "@/mixins/infoList";

export default {
  name: "FilterView",
  mixins: [infoList],
  data() {
    return {
      listApi: getFilterList,
      dialogFormVisible: false,
      type: "",
      multipleSelection: [],
      formData: {
        filter_type: "blacklist",
        filter_class: "extension",
        content: "",
      },
    };
  },
  methods: {
    formatDate,
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    deleteRow(row) {
      this.$confirm("确定要删除吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => this.deleteFilter(row))
        .catch(() => {});
    },
    confirmDelete() {
      if (!this.multipleSelection.length) {
        this.$message({ type: "warning", message: "请选择要删除的数据" });
        return;
      }
      this.$confirm(
        `确定删除选中的 ${this.multipleSelection.length} 条过滤规则吗？`,
        "批量删除",
        { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" }
      )
        .then(() => this.onDelete())
        .catch(() => {});
    },
    async onDelete() {
      const ids = this.multipleSelection.map((item) => item.ID);
      const res = await deleteFilterByIds({ ids });
      if (res.code == 0) {
        this.$message({ type: "success", message: "删除成功" });
        if (this.tableData.length == ids.length) this.page--;
        this.getTableData();
      }
    },
    async updateFilter(row) {
      const res = await findFilter({ ID: row.ID });
      this.type = "update";
      if (res.code === 0) {
        this.formData = res.data.filter;
        this.dialogFormVisible = true;
      }
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
        filter_type: "blacklist",
        filter_class: "extension",
        content: "",
      };
    },
    async deleteFilter(row) {
      const res = await deleteFilter({ ID: row.ID });
      if (res.code === 0) {
        this.$message({ type: "success", message: "删除成功" });
        if (this.tableData.length == 1) this.page--;
        await this.getTableData();
      }
    },
    async enterDialog() {
      let res;
      switch (this.type) {
        case "update":
          res = await updateFilter(this.formData);
          break;
        default:
          res = await createFilter(this.formData);
          break;
      }
      if (res.code === 0) {
        this.$message({ type: "success", message: "创建/更改成功" });
        this.closeDialog();
        await this.getTableData();
      }
    },
    openDialog() {
      this.type = "create";
      this.dialogFormVisible = true;
    },
  },
  async created() {
    await this.getTableData();
  },
};
</script>

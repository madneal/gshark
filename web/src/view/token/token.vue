<template>
  <div>
    <div class="search-term page-toolbar">
      <div class="page-toolbar__filters">
        <el-button type="primary" @click="onSubmit">刷新</el-button>
      </div>
      <div class="page-toolbar__actions">
        <el-button type="primary" @click="openDialog">新增 Token</el-button>
        <el-button
          type="danger"
          :disabled="multipleSelection.length === 0"
          @click="confirmDelete"
        >
          批量删除
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
      empty-text="暂无 Token"
    >
      <el-table-column type="selection" width="48" />
      <el-table-column label="日期" width="170">
        <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
      </el-table-column>
      <el-table-column label="类型" prop="type" width="120">
        <template #default="scope">
          <span class="keyword-chip">{{ scope.row.type }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Token" prop="content" min-width="200" show-overflow-tooltip />
      <el-table-column label="操作" width="170" fixed="right">
        <template #default="scope">
          <el-button size="small" type="primary" @click="updateToken(scope.row)">
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

    <el-dialog :before-close="closeDialog" v-model="dialogFormVisible" title="添加 Token">
      <el-form :model="formData" label-position="right" label-width="80px">
        <el-form-item label="类型:">
          <el-radio-group v-model="formData.type">
            <el-radio label="github">github</el-radio>
            <el-radio label="gitlab">gitlab</el-radio>
            <el-radio label="postman">postman</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="Token:">
          <el-input v-model="formData.content" clearable placeholder="请输入" />
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
  createToken,
  deleteToken,
  deleteTokenByIds,
  updateToken,
  findToken,
  getTokenList,
} from "@/api/token";
import { formatDate } from "@/utils/date";
import infoList from "@/mixins/infoList";

export default {
  name: "Token",
  mixins: [infoList],
  data() {
    return {
      listApi: getTokenList,
      dialogFormVisible: false,
      type: "",
      multipleSelection: [],
      formData: {
        type: "github",
        content: "",
      },
      typeOptions: [
        { label: "github", value: "github" },
        { label: "gitlab", value: "gitlab" },
        { label: "postman", value: "postman" },
      ],
    };
  },
  methods: {
    formatDate,
    onSubmit() {
      this.page = 1;
      this.pageSize = 10;
      this.getTableData();
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    deleteRow(row) {
      this.$confirm("确定要删除吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.deleteToken(row);
        })
        .catch(() => {});
    },
    confirmDelete() {
      if (this.multipleSelection.length === 0) {
        this.$message({ type: "warning", message: "请选择要删除的数据" });
        return;
      }
      this.$confirm(
        `确定删除选中的 ${this.multipleSelection.length} 个 Token 吗？`,
        "批量删除",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }
      )
        .then(() => this.onDelete())
        .catch(() => {});
    },
    async onDelete() {
      const ids = this.multipleSelection.map((item) => item.ID);
      const res = await deleteTokenByIds({ ids });
      if (res.code == 0) {
        this.$message({ type: "success", message: "删除成功" });
        if (this.tableData.length == ids.length) {
          this.page--;
        }
        this.getTableData();
      }
    },
    async updateToken(row) {
      const res = await findToken({ ID: row.ID });
      this.type = "update";
      if (res.code == 0) {
        this.formData = res.data.retoken;
        this.dialogFormVisible = true;
      }
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
        type: "github",
        content: "",
      };
    },
    async deleteToken(row) {
      const res = await deleteToken({ ID: row.ID });
      if (res.code == 0) {
        this.$message({ type: "success", message: "删除成功" });
        if (this.tableData.length == 1) {
          this.page--;
        }
        this.getTableData();
      }
    },
    async enterDialog() {
      let res;
      switch (this.type) {
        case "create":
          res = await createToken(this.formData);
          break;
        case "update":
          res = await updateToken(this.formData);
          break;
        default:
          res = await createToken(this.formData);
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

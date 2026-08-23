<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="规则类型">
          <el-input placeholder="搜索条件" clearable v-model="searchInfo.type" />
        </el-form-item>
        <el-form-item label="规则内容">
          <el-input placeholder="搜索条件" clearable v-model="searchInfo.content" />
        </el-form-item>
        <el-form-item label="规则名称">
          <el-input placeholder="搜索条件" clearable v-model="searchInfo.name" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button type="primary" @click="openDialog">新增规则</el-button>
          <el-button
            type="danger"
            :disabled="multipleSelection.length === 0"
            @click="confirmDelete"
          >
            批量删除
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-upload
            action="/api/rule/uploadRules"
            ref="ruleData"
            :with-credentials="true"
            :headers="{ 'x-token': $store.getters['user/token'] }"
            :show-file-list="false"
            :on-success="uploadSuccess"
          >
            <template #trigger>
              <el-button>规则导入</el-button>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-link
            href="https://github.com/madneal/gshark/blob/master/template.csv"
            type="primary"
            target="_blank"
          >
            规则导入模板
          </el-link>
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
      <el-table-column label="规则类型" prop="ruleType" width="120" />
      <el-table-column label="规则内容" prop="content" width="120" />
      <el-table-column label="本地匹配正则" prop="matchPattern" width="180" />
      <el-table-column label="规则名称" prop="name" width="120" />
      <el-table-column label="规则描述" prop="desc" width="120" />
      <el-table-column label="状态" width="120">
        <template #default="scope">
          <el-switch
            v-model="scope.row.status"
            :active-value="true"
            :inactive-value="false"
            @change="switchStatus(scope.row.ID, scope.row.status)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="scope">
          <el-button size="small" type="primary" @click="updateRule(scope.row)">
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

    <el-dialog :before-close="closeDialog" v-model="dialogFormVisible" title="新增规则">
      <el-form :model="formData" label-position="right" label-width="100px">
        <el-form-item label="规则类型:" required>
          <el-checkbox-group v-model="formData.ruleType">
            <el-checkbox v-for="ruleType in types" :label="ruleType" :key="ruleType">
              {{ typeLabels[ruleType] || ruleType }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="规则内容:" required>
          <el-input v-model="formData.content" clearable placeholder="请输入关键词内容" />
        </el-form-item>
        <el-form-item label="本地匹配正则:">
          <el-input
            v-model="formData.matchPattern"
            type="textarea"
            :rows="3"
            clearable
            placeholder="可选，仅对返回的代码片段进行本地校验"
          />
        </el-form-item>
        <el-form-item label="规则名称:">
          <el-input v-model="formData.name" clearable placeholder="请输入" />
        </el-form-item>
        <el-form-item label="规则描述:">
          <el-input v-model="formData.desc" clearable placeholder="请输入" />
        </el-form-item>
        <el-form-item label="状态:">
          <el-switch v-model="formData.status" />
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
  createRule,
  deleteRule,
  deleteRuleByIds,
  updateRule,
  findRule,
  getRuleList,
  switchRule,
} from "@/api/rule";
import { formatDate } from "@/utils/date";
import infoList from "@/mixins/infoList";

export default {
  name: "Rule",
  mixins: [infoList],
  data() {
    return {
      listApi: getRuleList,
      dialogFormVisible: false,
      type: "",
      multipleSelection: [],
      formData: {
        ruleType: [],
        content: "",
        matchPattern: "",
        name: "",
        desc: "",
        status: true,
      },
      types: ["github", "github_issue", "gist", "gitlab", "sourcegraph", "domain", "postman"],
      typeLabels: {
        github: "github",
        github_issue: "github issue/pr",
        gist: "gist",
        gitlab: "gitlab",
        sourcegraph: "sourcegraph",
        domain: "domain",
        postman: "postman",
      },
    };
  },
  methods: {
    formatDate,
    async switchStatus(id, status) {
      status = status ? 1 : 0;
      const res = await switchRule({ id, status });
      if (res) await this.getTableData();
    },
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
        .then(() => this.deleteRule(row))
        .catch(() => {});
    },
    confirmDelete() {
      if (!this.multipleSelection.length) {
        this.$message({ type: "warning", message: "请选择要删除的数据" });
        return;
      }
      this.$confirm(
        `确定删除选中的 ${this.multipleSelection.length} 条规则吗？`,
        "批量删除",
        { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" }
      )
        .then(() => this.onDelete())
        .catch(() => {});
    },
    uploadSuccess() {
      this.$message({ type: "success", message: "规则导入成功" });
      this.getTableData();
    },
    async onDelete() {
      const ids = this.multipleSelection.map((item) => item.ID);
      const res = await deleteRuleByIds({ ids });
      if (res.code === 0) {
        this.$message({ type: "success", message: "删除成功" });
        if (this.tableData.length == ids.length) this.page--;
        await this.getTableData();
      }
    },
    async updateRule(row) {
      const res = await findRule({ ID: row.ID });
      this.type = "update";
      if (res.code === 0) {
        this.formData = res.data.rule;
        this.formData.ruleType = this.formData.ruleType.split(",");
        this.dialogFormVisible = true;
      }
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
        ruleType: [],
        content: "",
        matchPattern: "",
        name: "",
        desc: "",
        status: 0,
      };
    },
    async deleteRule(row) {
      const res = await deleteRule({ ID: row.ID });
      if (res.code == 0) {
        this.$message({ type: "success", message: "删除成功" });
        if (this.tableData.length == 1) this.page--;
        await this.getTableData();
      }
    },
    async enterDialog() {
      this.formData.ruleType = this.formData.ruleType.join(",");
      let res;
      switch (this.type) {
        case "update":
          res = await updateRule(this.formData);
          break;
        default:
          res = await createRule(this.formData);
          break;
      }
      if (res.code == 0) {
        this.$message({ type: "success", message: "创建/更改成功" });
        this.closeDialog();
        await this.getTableData();
      }
    },
    openDialog() {
      this.type = "create";
      this.dialogFormVisible = true;
      this.formData.ruleType = [];
    },
  },
  async created() {
    await this.getTableData();
  },
};
</script>

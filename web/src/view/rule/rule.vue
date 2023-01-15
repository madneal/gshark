<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="规则类型">
          <el-input placeholder="搜索条件" v-model="searchInfo.type"></el-input>
        </el-form-item>
        <el-form-item label="规则内容">
          <el-input
            placeholder="搜索条件"
            v-model="searchInfo.content"
          ></el-input>
        </el-form-item>
        <el-form-item label="规则名称">
          <el-input placeholder="搜索条件" v-model="searchInfo.name"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="openDialog" type="primary">新增规则</el-button>
        </el-form-item>
        <el-form-item>
          <el-popover placement="top" v-model="deleteVisible" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin: 0">
              <el-button @click="deleteVisible = false" size="mini" type="text"
                >取消</el-button
              >
              <el-button @click="onDelete" size="mini" type="primary"
                >确定</el-button
              >
            </div>
            <el-button
              icon="el-icon-delete"
              size="mini"
              slot="reference"
              type="danger"
              >批量删除</el-button
            >
          </el-popover>
        </el-form-item>
        <el-form-item>
          <el-upload action="/api/rule/uploadRules" ref="ruleData" :with-credentials="true"
          :headers="{ 'x-token': $store.getters['user/token'] }" :show-file-list="false"
          :on-success="uploadSuccess">
            <template #trigger>
              <el-button type="primary">规则导入</el-button>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-link href="https://github.com/madneal/gshark/blob/master/template.csv" type="primary">规则导入模板</el-link>
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
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column label="日期" width="180">
        <template slot-scope="scope">{{
          scope.row.CreatedAt | formatDate
        }}</template>
      </el-table-column>

      <el-table-column
        label="规则类型"
        prop="ruleType"
        width="120"
      ></el-table-column>

      <el-table-column
        label="规则内容"
        prop="content"
        width="120"
      ></el-table-column>

      <el-table-column
        label="规则名称"
        prop="name"
        width="120"
      ></el-table-column>

      <el-table-column
        label="规则描述"
        prop="desc"
        width="120"
      ></el-table-column>

      <el-table-column label="状态" width="120">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.status"
            :active-value="true"
            :inactive-value="false"
            @change="switchStatus(scope.row.ID, scope.row.status)"
          />
        </template>
      </el-table-column>

      <el-table-column label="按钮组">
        <template slot-scope="scope">
          <el-button
            class="table-button"
            @click="updateRule(scope.row)"
            size="small"
            type="primary"
            icon="el-icon-edit"
            >变更</el-button
          >
          <el-button
            type="danger"
            icon="el-icon-delete"
            size="mini"
            @click="deleteRow(scope.row)"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{ float: 'right', padding: '20px' }"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>

    <el-dialog
      :before-close="closeDialog"
      :visible.sync="dialogFormVisible"
      title="新增规则"
    >
      <el-form :model="formData" label-position="right" label-width="100px">
        <el-form-item label="规则类型:" required>
          <el-checkbox-group v-model="formData.ruleType">
            <el-checkbox v-for="ruleType in types" :label="ruleType" :key="ruleType">{{ ruleType }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="规则内容:" required>
          <el-input
            v-model="formData.content"
            clearable
            placeholder="请输入关键词内容"
          ></el-input>
        </el-form-item>

        <el-form-item label="规则名称:">
          <el-input
            v-model="formData.name"
            clearable
            placeholder="请输入"
          ></el-input>
        </el-form-item>

        <el-form-item label="规则描述:">
          <el-input
            v-model="formData.desc"
            clearable
            placeholder="请输入"
          ></el-input>
        </el-form-item>

        <el-form-item label="状态:">
          <el-switch v-model="formData.status"></el-switch>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button @click="enterDialog" type="primary">确 定</el-button>
      </div>
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
  switchRule, batchCreateRules,
} from "@/api/rule";
import { formatTimeToStr } from "@/utils/date";
import infoList from "@/mixins/infoList";
import { store } from '@/store/index';

export default {
  name: "Rule",
  mixins: [infoList],
  data() {
    return {
      listApi: getRuleList,
      dialogFormVisible: false,
      dialogBatchRules: false,
      type: "",
      deleteVisible: false,
      multipleSelection: [],
      formData: {
        ruleType: [],
        content: "",
        name: "",
        desc: "",
        status: true,
      },
      batchRulesForm: {
        type: [],
        contents: ''
      },
      statusOptions: [
        {
          label: "disabled",
          value: 0,
        },
        {
          label: "enabled",
          value: 1,
        },
      ],
      types: ['github', 'gitlab', 'searchcode', 'domain'],
      typeOptions: [
        {
          label: "github",
          value: "github",
        },
        {
          label: "gitlab",
          value: "gitlab",
        },
        {
          label: "searchcode",
          value: "searchcode",
        },
        {
          label: "postman",
          value: "postman"
        }
      ],
    };
  },
  filters: {
    formatDate: function (time) {
      if (time != null && time != "") {
        var date = new Date(time);
        return formatTimeToStr(date, "yyyy-MM-dd hh:mm:ss");
      } else {
        return "";
      }
    },
    formatBoolean: function (bool) {
      if (bool != null) {
        return bool ? "是" : "否";
      } else {
        return "";
      }
    },
    formatStatus: function (status) {
      const statusMap = {
        0: "disabled",
        1: "enabled",
      };
      return statusMap[status];
    },
  },
  methods: {
    async switchStatus(id, status) {
      const data = {
        id,
        status,
      };
      const res = await switchRule(data);
      if (res) {
        await this.getTableData();
      }
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
      }).then(() => {
        this.deleteRule(row);
      });
    },
    uploadSuccess() {
      this.$message({
        type: "success",
        message: "规则导入成功"
      });
      this.getTableData();
    },
    async onDelete() {
      const ids = [];
      if (this.multipleSelection.length == 0) {
        this.$message({
          type: "warning",
          message: "请选择要删除的数据",
        });
        return;
      }
      this.multipleSelection &&
        this.multipleSelection.map((item) => {
          ids.push(item.ID);
        });
      const res = await deleteRuleByIds({ ids });
      if (res.code === 0) {
        this.$message({
          type: "success",
          message: "删除成功",
        });
        if (this.tableData.length == ids.length) {
          this.page--;
        }
        this.deleteVisible = false;
        await this.getTableData();
      }
    },
    async updateRule(row) {
      const res = await findRule({ ID: row.ID });
      this.type = "update";
      if (res.code === 0) {
        this.formData = res.data.rule;
        this.dialogFormVisible = true;
      }
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
        ruleType: [],
        content: "",
        name: "",
        desc: "",
        status: 0,
      };
    },
    async deleteRule(row) {
      const res = await deleteRule({ ID: row.ID });
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "删除成功",
        });
        if (this.tableData.length == 1) {
          this.page--;
        }
        await this.getTableData();
      }
    },
    async enterDialog() {
      let res;
      this.formData.ruleType = this.formData.ruleType.join(',');
      switch (this.type) {
        case "create":
          res = await createRule(this.formData);
          break;
        case "update":
          res = await updateRule(this.formData);
          break;
        default:
          res = await createRule(this.formData);
          break;
      }
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "创建/更改成功",
        });
        this.closeDialog();
        await this.getTableData();
      }
    },
    openDialog() {
      this.type = "create";
      this.dialogFormVisible = true;
      this.formData.ruleType = [];
    },
    async batchCreateRules() {
      this.batchRulesForm.type = this.batchRulesForm.type.join(',');
      const res = await batchCreateRules(this.batchRulesForm);
      if (res.code === 0) {
        this.dialogBatchRules = false;
        this.$message({
          type: 'success',
          message: '批量创建规则成功'
        });
        await this.getTableData();
      }
    }
  },
  async created() {
    await this.getTableData();
  },
};
</script>

<style>
</style>

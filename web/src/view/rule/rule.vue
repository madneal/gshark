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
        <el-form-item label="状态">
          <el-input
            placeholder="搜索条件"
            v-model="searchInfo.status"
          ></el-input>
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
        prop="type"
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

      <el-table-column label="状态" prop="status" width="120">
        <template slot-scope="scope">
          {{ scope.row.status | formatStatus }}
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
      title="弹窗操作"
    >
      <el-form :model="formData" label-position="right" label-width="80px">
        <el-form-item label="规则类型:">
          <el-select v-model="formData.type">
            <el-option
              v-for="item in typeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="规则内容:">
          <el-input
            v-model="formData.content"
            clearable
            placeholder="请输入"
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
          <el-select placeholder="请选择" v-model="formData.status">
            <el-option
              :key="item.value"
              :label="item.label"
              :value="item.value"
              v-for="item in statusOptions"
            ></el-option>
          </el-select>
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
} from "@/api/rule"; //  此处请自行替换地址
import { formatTimeToStr } from "@/utils/date";
import infoList from "@/mixins/infoList";
export default {
  name: "Rule",
  mixins: [infoList],
  data() {
    return {
      listApi: getRuleList,
      dialogFormVisible: false,
      type: "",
      deleteVisible: false,
      multipleSelection: [],
      formData: {
        type: "",
        content: "",
        name: "",
        desc: "",
        status: 1,
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
          value: "searchcode"
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
    //条件搜索前端看此方法
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
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "删除成功",
        });
        if (this.tableData.length == ids.length) {
          this.page--;
        }
        this.deleteVisible = false;
        this.getTableData();
      }
    },
    async updateRule(row) {
      const res = await findRule({ ID: row.ID });
      this.type = "update";
      if (res.code == 0) {
        this.formData = res.data.rerule;
        this.dialogFormVisible = true;
      }
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
        type: "",
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
        this.getTableData();
      }
    },
    async enterDialog() {
      let res;
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
        this.getTableData();
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

<style>
</style>

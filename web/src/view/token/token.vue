<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="openDialog" type="primary">新增token</el-button>
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

      <el-table-column label="类型" prop="type" width="120"></el-table-column>

      <el-table-column
        label="token"
        prop="content"
        width="120"
      ></el-table-column>

      <el-table-column label="描述" prop="desc" width="120"></el-table-column>

      <el-table-column label="limit" prop="limit" width="120"></el-table-column>

      <el-table-column
        label="remaining"
        prop="remaining"
        width="120"
      ></el-table-column>

      <el-table-column
        label="重置时间"
        prop="resetTime"
        width="120"
      ></el-table-column>

      <el-table-column label="按钮组">
        <template slot-scope="scope">
          <el-button
            class="table-button"
            @click="updateToken(scope.row)"
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
      title="添加token"
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

        <el-form-item label="token:">
          <el-input
            v-model="formData.content"
            clearable
            placeholder="请输入"
          ></el-input>
        </el-form-item>

        <el-form-item label="描述:">
          <el-input
            v-model="formData.desc"
            clearable
            placeholder="请输入"
          ></el-input>
        </el-form-item>

        <el-form-item label="limit:"
          ><el-input
            v-model.number="formData.limit"
            clearable
            placeholder="请输入"
          ></el-input>
        </el-form-item>

        <el-form-item label="remaining:"
          ><el-input
            v-model.number="formData.remaining"
            clearable
            placeholder="请输入"
          ></el-input>
        </el-form-item>

        <el-form-item label="重置时间:">
          <el-date-picker
            type="date"
            placeholder="选择日期"
            v-model="formData.resetTime"
            clearable
          ></el-date-picker>
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
  createToken,
  deleteToken,
  deleteTokenByIds,
  updateToken,
  findToken,
  getTokenList,
} from "@/api/token"; //  此处请自行替换地址
import { formatTimeToStr } from "@/utils/date";
import infoList from "@/mixins/infoList";
export default {
  name: "Token",
  mixins: [infoList],
  data() {
    return {
      listApi: getTokenList,
      dialogFormVisible: false,
      type: "",
      deleteVisible: false,
      multipleSelection: [],
      formData: {
        type: "",
        content: "",
        desc: "",
        limit: 0,
        remaining: 0,
        resetTime: new Date(),
      },
      typeOptions: [
        {
          label: "github",
          value: "github",
        },
        {
          label: "gitlab",
          value: "gitlab",
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
        this.deleteToken(row);
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
      const res = await deleteTokenByIds({ ids });
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
        type: "",
        content: "",
        desc: "",
        limit: 0,
        remaining: 0,
        resetTime: new Date(),
      };
    },
    async deleteToken(row) {
      const res = await deleteToken({ ID: row.ID });
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
          res = await createToken(this.formData);
          break;
        case "update":
          res = await updateToken(this.formData);
          break;
        default:
          res = await createToken(this.formData);
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

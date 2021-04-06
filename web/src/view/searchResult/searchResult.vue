<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="仓库名称">
          <el-input placeholder="搜索条件" v-model="searchInfo.repo"></el-input>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            placeholder="搜索条件"
            v-model="searchInfo.keyword"
          ></el-input>
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
              >批量忽略</el-button
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
        label="仓库名称"
        prop="repo"
        width="120"
      ></el-table-column>

      <el-table-column
        label="匹配内容"
        prop="matches"
        width="120"
      ></el-table-column>

      <el-table-column
        label="关键词"
        prop="keyword"
        width="120"
      ></el-table-column>

      <el-table-column label="路径" prop="path" width="120"></el-table-column>

      <el-table-column label="URL" width="120">
        <template slot-scope="scope">
        <el-link :href="scope.row.url" type="primary">{{ scope.row.path }}</el-link>
        </template>
      </el-table-column>

      <el-table-column label="状态" prop="status" width="120">
        
        <template slot-scope="scope">
            {{ scope.row.status | statusFilter }}
        </template>
      </el-table-column>

      <el-table-column label="Review">
        <template slot-scope="scope">
          <el-button
            class="table-button"
            @click="updateSearchResult(scope.row)"
            size="small"
            type="primary"
            icon="el-icon-edit"
            >确认</el-button
          >
          <el-button
            type="danger"
            icon="el-icon-delete"
            size="mini"
            @click="deleteRow(scope.row)"
            >忽略</el-button
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

  </div>
</template>

<script>
import {
  deleteSearchResult,
  deleteSearchResultByIds,
  findSearchResult,
  getSearchResultList,
} from "@/api/searchResult"; //  此处请自行替换地址
import { formatTimeToStr } from "@/utils/date";
import infoList from "@/mixins/infoList";
export default {
  name: "SearchResult",
  mixins: [infoList],
  data() {
    return {
      listApi: getSearchResultList,
      dialogFormVisible: false,
      type: "",
      deleteVisible: false,
      statusOptions: [
        {
          未处理: 0,
          已确认: 1,
          已忽略: 2,
        },
      ],
      multipleSelection: [],
      formData: {
        repo: "",
        matches: "",
        keyword: "",
        path: "",
        url: "",
        textmatchMd5: "",
        status: 0,
      },
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
    statusFilter: function(val) {
      const statusOptions = {
        0: '未处理',
        1: '已处理',
        2: '已忽略'
      };
      return statusOptions[val];
    }
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
        this.deleteSearchResult(row);
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
      const res = await deleteSearchResultByIds({ ids });
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
    async updateSearchResult(row) {
      const res = await findSearchResult({ ID: row.ID });
      this.type = "update";
      if (res.code == 0) {
        this.formData = res.data.researchResult;
        this.dialogFormVisible = true;
      }
    },
    async deleteSearchResult(row) {
      const res = await deleteSearchResult({ ID: row.ID });
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
    }
  },
  async created() {
    await this.getTableData();
  },
};
</script>

<style>
</style>

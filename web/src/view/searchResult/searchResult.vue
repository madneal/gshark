<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="搜索条件">
          <el-input
            placeholder="仓库名称|匹配内容"
            v-model="searchInfo.query"
          ></el-input>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            placeholder="搜索条件"
            v-model="searchInfo.keyword"
          ></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>

        <el-form-item>
          <el-popover placement="top" v-model="deleteVisible" width="160">
            <p>确定要忽略吗？</p>
            <div style="text-align: right; margin: 0">
              <el-button @click="deleteVisible = false" size="mini" type="text"
                >取消</el-button
              >
              <el-button @click="onChange(true)" size="mini" type="primary"
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

        <el-form-item>
          <el-popover placement="top" v-model="confirmVisible" width="160">
            <p>确定要确认吗？</p>
            <div style="text-align: right; margin: 0">
              <el-button @click="confirmVisible = false" size="mini" type="text"
                >取消</el-button
              >
              <el-button @click="onChange(false)" size="mini" type="primary"
                >确定</el-button
              >
            </div>
            <el-button
              size="mini"
              slot="reference"
              type="primary"
              >批量确认</el-button
            >
          </el-popover>
        </el-form-item>

        <el-form-item>
          <el-popover placement="top" v-model="taskVisible" width="160">
            <p>确定启动过滤任务吗？过滤任务将自动忽略未匹配二次过滤关键词的结果</p>
            <div style="text-align: right; margin: 0">
              <el-button @click="taskVisible = false" size="mini" type="text"
              >取消</el-button
              >
              <el-button @click="startFilterTask" size="mini" type="primary"
              >确定</el-button
              >
            </div>
            <el-button
                size="mini"
                slot="reference"
                type="primary"
                :disabled="taskBtnDisable"
            >{{ taskButtonTxt }}</el-button
            >
          </el-popover>
        </el-form-item>

        <el-form-item label="仅展示二次关键词结果">
          <el-switch v-model="searchInfo.onlySecKeyword" @change="secKeywordChange"></el-switch>
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

      <el-table-column label="ID" prop="ID" width="50"></el-table-column>

      <el-table-column label="文件" width="180">
        <template slot-scope="scope">
          <a :href="scope.row.url" target="_blank" style="color: #409EFF;">
            {{ scope.row.repo + "/" + scope.row.path}}
          </a>
        </template>
      </el-table-column>

      <el-table-column label="匹配内容" prop="matches" width="550">
        <template slot-scope="scope">
          <pre>{{ scope.row.text_matches | fragmentsFilter }}</pre>
        </template>
      </el-table-column>

      <el-table-column
        label="关键词"
        prop="keyword"
        width="80"
      ></el-table-column>

      <el-table-column label="二级关键词" prop="sec_keyword" width="100"></el-table-column>

      <el-table-column label="日期" width="100">
        <template slot-scope="scope">{{
          scope.row.CreatedAt | formatDate
        }}</template>
      </el-table-column>

      <el-table-column label="状态" prop="status" width="80">
        <template slot-scope="scope">
          {{ scope.row.status | statusFilter }}
        </template>
      </el-table-column>

      <el-table-column label="Review">
        <template slot-scope="scope">
          <el-button
            class="table-button"
            @click="updateSearchResult(scope.row, 1)"
            size="small"
            type="primary"
            icon="el-icon-edit"
            :disabled="scope.row.status === 1 ? true : false"
            >确认</el-button
          >
          <el-button
            type="danger"
            icon="el-icon-delete"
            size="mini"
            @click="updateSearchResult(scope.row, 2)"
            :disabled="scope.row.status === 2 ? true : false"
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
  findSearchResult,
  getSearchResultList,
  updateSearchResult,
  updateSearchResultStatusByIds,
  startFilterTask,
  getTaskStatus
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
      confirmVisible: false,
      taskVisible: false,
      taskButtonTxt: "启动二次过滤",
      taskBtnDisable: false,
      onlySecKeyword: false,
      statusOptions: [
        {
          label: "未处理",
          value: 0,
        },
        {
          label: "已确认",
          value: 1,
        },
        {
          label: "已忽略",
          value: 2,
        },
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
    statusFilter: function (val) {
      const statusOptions = {
        0: "未处理",
        1: "已处理",
        2: "已忽略",
      };
      return statusOptions[val];
    },
    fragmentsFilter: function (val) {
      if (!val) {
        return "";
      }
      let result = "";
      for (let i = 0; i < val.length; i++) {
        const matches = val[i].matches;
        let index = [];
        matches.forEach(ele => {
          index = index.concat(ele.indices);
        });
        let fragment = val[i].fragment;
        fragment = fragment.slice(0, index[0]) + "【" + fragment.slice(index[0]);
        fragment = fragment.slice(0, index[index.length-1]+1) + "】" + fragment.slice(index[index.length-1]+1);
        result = result + fragment;
        if (i !== val.length - 1) {
          result = result + "\n=====================================\n";
        }
      }
      return result;
    },
  },
  methods: {
    onSubmit() {
      this.page = 1;
      this.pageSize = 100;
      this.getTableData();
    },
    secKeywordChange() {
      this.getTableData();
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    async startFilterTask() {
      await startFilterTask();
      const resp = await getTaskStatus();
      if (resp.msg === "running") {
        this.taskButtonTxt = "任务运行中";
        this.taskBtnDisable = true;
      }
      this.taskVisible = false;
    },
    async onChange(isIgnore) {
      const ids = [];
      if (this.multipleSelection.length === 0) {
        this.$message({
          type: "warning",
          message: "请选择要操作的数据",
        });
        return;
      }
      this.multipleSelection &&
        this.multipleSelection.map((item) => {
          ids.push(item.ID);
        });
      let res;
      if (isIgnore) {
        res = await updateSearchResultStatusByIds({ ids, status: 2 });
      } else {
        res = await updateSearchResultStatusByIds({ ids , status: 1 });
      }
      if (res.code === 0) {
        this.$message({
          type: "success",
          message: "操作成功",
        });
        if (this.tableData.length == ids.length) {
          this.page--;
        }
        if (isIgnore) {
          this.deleteVisible = false;
        } else {
          this.confirmVisible = false;
        }
        await this.getTableData();
      }
    },
    async updateSearchResult(row, status) {
      const res = await findSearchResult({ ID: row.ID });
      this.type = "update";
      res.data.searchResult.status = status;
      if (res.code === 0) {
        this.formData = res.data.searchResult;
        this.dialogFormVisible = true;
        const data = {
          repo: this.formData.repo,
          status: status
        };
        const updateRes = await updateSearchResult(data);
        if (updateRes.code === 0) {
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

<style>
.el-table pre {
  white-space: pre-line;
}
</style>

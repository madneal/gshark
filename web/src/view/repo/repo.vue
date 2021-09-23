<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">              
        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="openDialog" type="primary">新增代码仓库</el-button>
        </el-form-item>
        <el-form-item>
          <el-popover placement="top" v-model="deleteVisible" width="160">
            <p>确定要删除吗？</p>
              <div style="text-align: right; margin: 0">
                <el-button @click="deleteVisible = false" size="mini" type="text">取消</el-button>
                <el-button @click="onDelete" size="mini" type="primary">确定</el-button>
              </div>
            <el-button icon="el-icon-delete" size="mini" slot="reference" type="danger">批量删除</el-button>
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
         <template slot-scope="scope">{{scope.row.CreatedAt|formatDate}}</template>
    </el-table-column>
    
    <el-table-column label="项目ID" prop="projectId" width="120"></el-table-column> 
    
    <el-table-column label="仓库类型" prop="type" width="120"></el-table-column> 
    
    <el-table-column label="描述" prop="desc" width="120"></el-table-column> 
    
    <el-table-column label="URL" prop="url" width="120"></el-table-column> 
    
    <el-table-column label="路径" prop="path" width="120"></el-table-column> 
    
    <el-table-column label="上次活动时间" prop="lastActivityAt" width="120"></el-table-column> 
    
    <el-table-column label="状态" prop="status" width="120"></el-table-column> 
    
      <el-table-column label="按钮组">
        <template slot-scope="scope">
          <el-button class="table-button" @click="updateRepo(scope.row)" size="small" type="primary" icon="el-icon-edit">变更</el-button>
          <el-button type="danger" icon="el-icon-delete" size="mini" @click="deleteRow(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{float:'right',padding:'20px'}"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>

    <el-dialog :before-close="closeDialog" :visible.sync="dialogFormVisible" title="弹窗操作">
      <el-form :model="formData" label-position="right" label-width="80px">
         <el-form-item label="项目ID:"><el-input v-model.number="formData.projectId" clearable placeholder="请输入"></el-input>
      </el-form-item>
       
         <el-form-item label="仓库类型:">
            <el-input v-model="formData.type" clearable placeholder="请输入" ></el-input>
      </el-form-item>
       
         <el-form-item label="描述:">
            <el-input v-model="formData.desc" clearable placeholder="请输入" ></el-input>
      </el-form-item>
       
         <el-form-item label="URL:">
            <el-input v-model="formData.url" clearable placeholder="请输入" ></el-input>
      </el-form-item>
       
         <el-form-item label="路径:">
            <el-input v-model="formData.path" clearable placeholder="请输入" ></el-input>
      </el-form-item>
       
         <el-form-item label="上次活动时间:">
              <el-date-picker type="date" placeholder="选择日期" v-model="formData.lastActivityAt" clearable></el-date-picker>
       </el-form-item>
       
         <el-form-item label="状态:"><el-input v-model.number="formData.status" clearable placeholder="请输入"></el-input>
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
    createRepo,
    deleteRepo,
    deleteRepoByIds,
    updateRepo,
    findRepo,
    getRepoList
} from "@/api/repo";  //  此处请自行替换地址
import { formatTimeToStr } from "@/utils/date";
import infoList from "@/mixins/infoList";
export default {
  name: "Repo",
  mixins: [infoList],
  data() {
    return {
      listApi: getRepoList,
      dialogFormVisible: false,
      type: "",
      deleteVisible: false,
      multipleSelection: [],formData: {
            projectId:0,
            type:"",
            desc:"",
            url:"",
            path:"",
            lastActivityAt:new Date(),
            status:0,
            
      }
    };
  },
  filters: {
    formatDate: function(time) {
      if (time != null && time != "") {
        var date = new Date(time);
        return formatTimeToStr(date, "yyyy-MM-dd hh:mm:ss");
      } else {
        return "";
      }
    },
    formatBoolean: function(bool) {
      if (bool != null) {
        return bool ? "是" :"否";
      } else {
        return "";
      }
    }
  },
  methods: {
      //条件搜索前端看此方法
      onSubmit() {
        this.page = 1
        this.pageSize = 10           
        this.getTableData()
      },
      handleSelectionChange(val) {
        this.multipleSelection = val
      },
      deleteRow(row){
        this.$confirm('确定要删除吗?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
           this.deleteRepo(row);
        });
      },
      async onDelete() {
        const ids = []
        if(this.multipleSelection.length == 0){
          this.$message({
            type: 'warning',
            message: '请选择要删除的数据'
          })
          return
        }
        this.multipleSelection &&
          this.multipleSelection.map(item => {
            ids.push(item.ID)
          })
        const res = await deleteRepoByIds({ ids })
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功'
          })
          if (this.tableData.length == ids.length) {
              this.page--;
          }
          this.deleteVisible = false
          this.getTableData()
        }
      },
    async updateRepo(row) {
      const res = await findRepo({ ID: row.ID });
      this.type = "update";
      if (res.code == 0) {
        this.formData = res.data.rerepo;
        this.dialogFormVisible = true;
      }
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
          projectId:0,
          type:"",
          desc:"",
          url:"",
          path:"",
          lastActivityAt:new Date(),
          status:0,
          
      };
    },
    async deleteRepo(row) {
      const res = await deleteRepo({ ID: row.ID });
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "删除成功"
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
          res = await createRepo(this.formData);
          break;
        case "update":
          res = await updateRepo(this.formData);
          break;
        default:
          res = await createRepo(this.formData);
          break;
      }
      if (res.code == 0) {
        this.$message({
          type:"success",
          message:"创建/更改成功"
        })
        this.closeDialog();
        this.getTableData();
      }
    },
    openDialog() {
      this.type = "create";
      this.dialogFormVisible = true;
    }
  },
  async created() {
    await this.getTableData();
  
}
};
</script>

<style>
</style>

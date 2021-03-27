<template>
<div>
    <el-form :model="formData" label-position="right" label-width="80px">
             <el-form-item label="仓库名称:">
                <el-input v-model="formData.repo" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="匹配内容:">
                <el-input v-model="formData.matches" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="关键词:">
                <el-input v-model="formData.keyword" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="路径:">
                <el-input v-model="formData.path" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="URL:">
                <el-input v-model="formData.url" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="哈希值:">
                <el-input v-model="formData.textmatchMd5" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="状态:"><el-input v-model.number="formData.status" clearable placeholder="请输入"></el-input>
          </el-form-item>
           <el-form-item>
           <el-button @click="save" type="primary">保存</el-button>
           <el-button @click="back" type="primary">返回</el-button>
           </el-form-item>
    </el-form>
</div>
</template>

<script>
import {
    createSearchResult,
    updateSearchResult,
    findSearchResult
} from "@/api/searchResult";  //  此处请自行替换地址
import infoList from "@/mixins/infoList";
export default {
  name: "SearchResult",
  mixins: [infoList],
  data() {
    return {
      type: "",formData: {
            repo:"",
            matches:"",
            keyword:"",
            path:"",
            url:"",
            textmatchMd5:"",
            status:0,
            
      }
    };
  },
  methods: {
    async save() {
      let res;
      switch (this.type) {
        case "create":
          res = await createSearchResult(this.formData);
          break;
        case "update":
          res = await updateSearchResult(this.formData);
          break;
        default:
          res = await createSearchResult(this.formData);
          break;
      }
      if (res.code == 0) {
        this.$message({
          type:"success",
          message:"创建/更改成功"
        })
      }
    },
    back(){
        this.$router.go(-1)
    }
  },
  async created() {
   // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if(this.$route.query.id){
    const res = await findSearchResult({ ID: this.$route.query.id })
    if(res.code == 0){
       this.formData = res.data.researchResult
       this.type = "update"
     }
    }else{
     this.type = "create"
   }
  
}
};
</script>

<style>
</style>
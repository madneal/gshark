<template>
<div>
    <el-form :model="formData" label-position="right" label-width="80px">
             <el-form-item label="文件后缀:">
                <el-input v-model="formData.extension" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="是否fork:">
                <el-switch active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" v-model="formData.isFork" clearable ></el-switch>
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
    createFilter,
    updateFilter,
    findFilter
} from "@/api/filter";  //  此处请自行替换地址
import infoList from "@/mixins/infoList";
export default {
  name: "Filter",
  mixins: [infoList],
  data() {
    return {
      type: "",formData: {
            extension:"",
            isFork:false,
            
      }
    };
  },
  methods: {
    async save() {
      let res;
      switch (this.type) {
        case "create":
          res = await createFilter(this.formData);
          break;
        case "update":
          res = await updateFilter(this.formData);
          break;
        default:
          res = await createFilter(this.formData);
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
    const res = await findFilter({ ID: this.$route.query.id })
    if(res.code == 0){
       this.formData = res.data.refilter
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
<template>
<div>
    <el-form :model="formData" label-position="right" label-width="80px">
             <el-form-item label="规则类型:">
                <el-input v-model="formData.type" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="规则内容:">
                <el-input v-model="formData.content" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="规则名称:">
                <el-input v-model="formData.name" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="规则描述:">
                <el-input v-model="formData.desc" clearable placeholder="请输入" ></el-input>
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
    createRule,
    updateRule,
    findRule
} from "@/api/rule";  //  此处请自行替换地址
import infoList from "@/mixins/infoList";
export default {
  name: "Rule",
  mixins: [infoList],
  data() {
    return {
      type: "",formData: {
            type:"",
            content:"",
            name:"",
            desc:"",
            status:0,
            
      }
    };
  },
  methods: {
    async save() {
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
    const res = await findRule({ ID: this.$route.query.id })
    if(res.code == 0){
       this.formData = res.data.rerule
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
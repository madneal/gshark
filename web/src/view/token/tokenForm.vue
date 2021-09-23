<template>
<div>
    <el-form :model="formData" label-position="right" label-width="80px">
             <el-form-item label="类型:">
                <el-input v-model="formData.type" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="token:">
                <el-input v-model="formData.content" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="描述:">
                <el-input v-model="formData.desc" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="limit:"><el-input v-model.number="formData.limit" clearable placeholder="请输入"></el-input>
          </el-form-item>
           
             <el-form-item label="remaining:"><el-input v-model.number="formData.remaining" clearable placeholder="请输入"></el-input>
          </el-form-item>
           
             <el-form-item label="重置时间:">
                  <el-date-picker type="date" placeholder="选择日期" v-model="formData.resetTime" clearable></el-date-picker>
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
    createToken,
    updateToken,
    findToken
} from "@/api/token";  //  此处请自行替换地址
import infoList from "@/mixins/infoList";
export default {
  name: "Token",
  mixins: [infoList],
  data() {
    return {
      type: "",formData: {
            type:"",
            content:"",
            desc:"",
            limit:0,
            remaining:0,
            resetTime:new Date(),
            
      }
    };
  },
  methods: {
    async save() {
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
    const res = await findToken({ ID: this.$route.query.id })
    if(res.code == 0){
       this.formData = res.data.retoken
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
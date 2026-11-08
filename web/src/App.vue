<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
import { checkDB } from "@/api/initdb"
export default {
  name: 'app',
  async created(){
    await checkDB()
  }  
}
</script>

<style lang="scss">
// 引入初始化样式
@import '@/style/main.scss';
@import '@/style/base.scss';
@import '@/style/mobile.scss';
#app {
  background: var(--gs-bg-page);
  height: 100vh;
  overflow: hidden;
}
</style>

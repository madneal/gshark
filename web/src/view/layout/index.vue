<template>
  <el-container class="layout-cont">
    <el-container :class="[isSider?'openside':'hideside',isMobile ? 'mobile': '']">
      <el-row :class="[isShadowBg?'shadowBg':'']" @click="changeShadow()"></el-row>
      <el-aside class="main-cont main-left">
        <div class="tilte">
          <img alt class="logoimg" src="@/assets/nav_logo.png" />
          <h2 class="tit-text" v-if="isSider">GShark</h2>
        </div>
        <LayoutAside class="aside" />
      </el-aside>
      <el-main class="main-cont main-right">
        <transition :duration="{ enter: 800, leave: 100 }" mode="out-in" name="el-fade-in-linear">
          <div
            :style="{width: `calc(100% - ${isMobile?'0px':isCollapse?'54px':'220px'})`}"
            class="topfix"
          >
          <el-header class="header-cont">
            <div class="header-left">
              <div @click="totalCollapse" class="menu-total">
                <el-icon v-if="isCollapse"><Expand /></el-icon>
                <el-icon v-else><Fold /></el-icon>
              </div>
              <el-breadcrumb class="breadcrumb" separator-class="el-icon-arrow-right">
                <el-breadcrumb-item
                  :key="item.path"
                  v-for="item in matched.slice(1,matched.length)"
                >{{item.meta.title}}</el-breadcrumb-item>
              </el-breadcrumb>
            </div>
            <div class="right-box">
                <Search />
                <Screenfull class="screenfull" :style="{cursor:'pointer'}"></Screenfull>
                <el-dropdown>
                  <span class="header-avatar">
                   <CustomPic/>
                    <span style="margin-left: 5px">{{userInfo.nickName}}</span>
                    <el-icon><ArrowDown /></el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu class="dropdown-group">
                      <el-dropdown-item @click="toPerson">
                        <el-icon><User /></el-icon>
                        个人信息
                      </el-dropdown-item>
                      <el-dropdown-item @click="LoginOut">
                        <el-icon><SwitchButton /></el-icon>
                        登 出
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
            </div>
          </el-header>
            <HistoryComponent />
          </div>
        </transition>
        <div v-loading="loadingFlag" element-loading-text="正在加载中" class="admin-box">
          <router-view v-slot="{ Component }">
            <component :is="Component" v-if="reloadFlag" />
          </router-view>
        </div>
      </el-main>
    </el-container>
   
  </el-container>
</template>

<script>
import LayoutAside from '@/view/layout/aside/index.vue'
import HistoryComponent from '@/view/layout/aside/historyComponent/history.vue'
import Screenfull from '@/view/layout/screenfull/index.vue'
import Search from '@/view/layout/search/search.vue'
import { mapGetters, mapActions } from 'vuex'
import CustomPic from '@/components/customPic/index.vue'
export default {
  name: 'Layout',
  data() {
    return {
      show: false,
      isCollapse: false,
      isSider: true,
      isMobile: false,
      isShadowBg: false,
      loadingFlag:false,
      reloadFlag:true,
      value: ''
    }
  },
  components: {
    LayoutAside,
    HistoryComponent,
    Screenfull,
    Search,
    CustomPic
  },
  methods: {
    ...mapActions('user', ['LoginOut']),
    reload(){
      this.reloadFlag = false
      this.$nextTick(()=>{
        this.reloadFlag = true
      })
    },
    totalCollapse() {
      this.isCollapse = !this.isCollapse
      this.isSider = !this.isCollapse
      this.isShadowBg = !this.isCollapse
      this.$bus.emit('collapse', this.isCollapse)
    },
    toPerson() {
      this.$router.push({ name: 'person' })
    },
    changeShadow() {
      this.isShadowBg = !this.isShadowBg
      this.isSider = !!this.isCollapse
      this.totalCollapse()
    },
  },
  computed: {
    ...mapGetters('user', ['userInfo']),
    title() {
      return this.$route.meta.title || '当前页面'
    },
    matched() {
      return this.$route.matched
    }
  },
  mounted() {
    let screenWidth = document.body.clientWidth
    if (screenWidth < 1000) {
      this.isMobile = true
      this.isSider = false
      this.isCollapse = true
    } else if (screenWidth >= 1000 && screenWidth < 1200) {
      this.isMobile = false
      this.isSider = false
      this.isCollapse = true
    } else {
      this.isMobile = false
      this.isSider = true
      this.isCollapse = false
    }
    this.$bus.emit('collapse', this.isCollapse)
    this.$bus.emit('mobile', this.isMobile)
    this.$bus.on("reload",this.reload)
    this.$bus.on("showLoading",()=>{
      this.loadingFlag = true
    })
    this.$bus.on("closeLoading",()=>{
      this.loadingFlag = false
    })
    window.onresize = () => {
      return (() => {
        let screenWidth = document.body.clientWidth
        if (screenWidth < 1000) {
          this.isMobile = true
          this.isSider = false
          this.isCollapse = true
        } else if (screenWidth >= 1000 && screenWidth < 1200) {
          this.isMobile = false
          this.isSider = false
          this.isCollapse = true
        } else {
          this.isMobile = false
          this.isSider = true
          this.isCollapse = false
        }
        this.$bus.emit('collapse', this.isCollapse)
        this.$bus.emit('mobile', this.isMobile)
      })()
    }
  }
}
</script>
<style lang="scss">
@import '@/style/mobile.scss';
</style>

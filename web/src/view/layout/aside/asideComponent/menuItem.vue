<template>
  <el-menu-item :index="routerInfo.name" :route="{parameters:routerInfo.parameters}">
    <el-icon>
      <component :is="iconName" />
    </el-icon>
    <template #title>{{routerInfo.meta.title}}</template>
  </el-menu-item>
</template>

<script>
import { menuIcon } from '@/utils/menuIcon'

export default {
  name: 'MenuItem',
  props: {
    routerInfo: {
      default: function() {
        return null
      },
      type: Object
    }
  },
  computed: {
    iconName() {
      return menuIcon(this.routerInfo.meta?.icon)
    }
  }
}
</script>
<style lang="scss">
</style>

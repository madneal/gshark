<template>
  <el-sub-menu :popper-append-to-body="false" ref="subMenu" :index="routerInfo.name">
    <template #title>
      <el-icon>
        <component :is="iconName" />
      </el-icon>
      <span>{{routerInfo.meta.title}}</span>
    </template>
    <slot></slot>
  </el-sub-menu>
</template>

<script>
import { menuIcon } from '@/utils/menuIcon'

export default {
  name: 'AsyncSubmenu',
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

<template>
  <!-- Align with routerHolder / superAdmin for Vue Router 4 nested outlets -->
  <div>
    <router-view v-slot="{ Component }">
      <component :is="Component" />
    </router-view>
  </div>
</template>

<script>
export default {
  name: "System",
};
</script>
<style lang="scss"></style>

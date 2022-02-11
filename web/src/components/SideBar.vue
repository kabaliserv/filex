<script lang="ts" setup></script>

<template>
  <div class="sidebar">
    <el-menu
        default-active="1"
        class="sidebar-menu"
    >
      <el-menu-item index="1">
        <template #title>Fichers</template>
      </el-menu-item>
      <el-menu-item index="2">
        <template #title>Partagers</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<style scoped>
.sidebar-menu {
  height: 100%;
  background-color: #eeeeee;
}
</style>

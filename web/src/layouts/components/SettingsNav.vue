<script lang="ts" setup>

type Item = {
  index: string,
  name: string,
  path: string
}
type Props = {
  title: string
  items?: Array<Item>
  activeMenu?: string
}
const props = withDefaults(defineProps<Props>(), {
  title: "Settings",
  items: () => new Array<Item>()
})
</script>

<template>
  <div class="settings-nav">
    <div class="nav-header">
      <span>{{props.title}}</span>
    </div>
    <el-menu
        :default-active="props.activeMenu"
        :router="true"
    >
      <el-menu-item v-for="(item, index) in props.items" :key="index" :index="item.index" :route="{ path: item.path}">
        <span>{{ item.name }}</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<style lang="scss">
.settings-nav {
  --border-color: #dedee0;

  border-radius: 5px;
  border: 1px solid var(--border-color);
  overflow: hidden;
  //max-width: 200px;
  //padding-bottom: 5px;
  .nav-header {
    padding: 10px  10px;
    padding-left: 22px;
    font-weight: bold;
    border-bottom: 1px solid var(--border-color);
  }

  .el-menu {
    border-right: none;

    .el-menu-item {
      line-height: 2rem;
      height: 40px;
      padding: 0 5px;
      border-bottom: 1px solid var(--border-color);
      border-left: 2px solid transparent;

      &:hover {
        background-color: #f6f6f6;
      }

      &:last-child {
        border-bottom: none;
      }

      &.is-active {
        color: initial;
        border-left-color: #58bfff;
      }
    }
  }
}
</style>
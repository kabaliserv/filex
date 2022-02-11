
<script lang="ts" setup>
import {HeaderBar, SideBar} from "@/components";
import {Fragment} from "vue";
</script>

<template>
  <div v-if="true" id="app">
    <HeaderBar />
    <main>
      <div class="main-wrapper">
<!--        <SideBar class="sidebar"/>-->
        <router-view class="pages-view"></router-view>
      </div>
    </main>
  </div>
  <router-view v-else></router-view>
</template>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /*text-align: center ;*/
  color: #2c3e50;
}
main {
  display: flex;
  flex-direction: row;
  height: 100%;
  flex: 1;
}

.main-wrapper {
  display: flex;
  flex-direction: row;
  height: 100%;
}

.sidebar {
  width: 200px;
  height: 100%;
  /*background-color: deepskyblue;*/
}

.pages-view {
  flex: 1 1 auto ;
  height: 100%;
}
</style>
